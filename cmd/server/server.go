/* Copyright (c) 2014-2015 Jason Ish
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED ``AS IS'' AND ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package server

import (
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/jasonish/evebox/config"
	"github.com/jasonish/evebox/core"
	"github.com/jasonish/evebox/elasticsearch"
	"github.com/jasonish/evebox/log"
	"github.com/jasonish/evebox/server"
	"github.com/jessevdk/go-flags"
	"net/http"
)

const DEFAULT_ELASTICSEARCH_URL string = "http://localhost:9200"

var opts struct {
	// We don't provide a default for this one so we can easily
	// detect if its been set or not.
	ElasticSearchUri   string `long:"elasticsearch" short:"e" description:"Elastic Search URI (default: http://localhost:9200)"`
	ElasticSearchIndex string `long:"index" short:"i" description:"Elastic Search Index (default: logstash-*)"`
	Port               string `long:"port" short:"p" default:"5636" description:"Port to bind to"`
	Host               string `long:"host" default:"0.0.0.0" description:"Host to bind to"`
	DevServerUri       string `long:"dev" description:"Frontend development server URI"`
	Version            bool   `long:"version" description:"Show version"`
	Config             string `long:"config" short:"c" description:"Configuration filename"`
	NoCheckCertificate bool   `long:"no-check-certificate" short:"k" description:"Disable certificate check for Elastic Search"`
}

var conf *config.Config

func init() {
	conf = config.NewConfig()
}

func VersionMain() {
	fmt.Printf("EveBox Version %s (rev %s) [%s]\n",
		core.BuildVersion, core.BuildRev, core.BuildDate)
}

func getElasticSearchUrl() string {
	if opts.ElasticSearchUri != "" {
		return opts.ElasticSearchUri
	}
	if os.Getenv("ELASTICSEARCH_URL") != "" {
		return os.Getenv("ELASTICSEARCH_URL")
	}
	return DEFAULT_ELASTICSEARCH_URL
}

func Main(args []string) {

	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		// flags.Parse should have already presented an error message.
		os.Exit(1)
	}

	if opts.Version {
		VersionMain()
		return
	}

	log.SetLevel(log.DEBUG)

	// If no configuration was provided, see if evebox.yaml exists
	// in the current directory.
	if opts.Config == "" {
		_, err = os.Stat("./evebox.yaml")
		if err == nil {
			opts.Config = "./evebox.yaml"
		}
	}
	if opts.Config != "" {
		log.Printf("Loading configuration file %s.\n", opts.Config)
		conf, err = config.LoadConfig(opts.Config)
		if err != nil {
			log.Fatal(err)
		}
	}

	if opts.ElasticSearchIndex != "" {
		conf.ElasticSearchIndex = opts.ElasticSearchIndex
	} else if os.Getenv("ELASTICSEARCH_INDEX") != "" {
		conf.ElasticSearchIndex = os.Getenv("ELASTICSEARCH_INDEX")
	} else {
		conf.ElasticSearchIndex = "logstash-*"
	}
	log.Info("Using ElasticSearch Index %s.", conf.ElasticSearchIndex)

	appContext := server.AppContext{
		Config: conf,
	}
	elasticSearch := elasticsearch.New(getElasticSearchUrl())
	elasticSearch.SetEventIndex(conf.ElasticSearchIndex)
	pingResponse, err := elasticSearch.Ping()
	if err != nil {
		log.Error("Failed to ping Elastic Search: %v", err)
	} else {
		log.Info("Connected to Elastic Search (version: %s)",
			pingResponse.Version.Number)
	}
	appContext.ElasticSearch = elasticSearch
	appContext.ArchiveService = elasticsearch.NewArchiveService(elasticSearch)
	appContext.EventService = elasticsearch.NewEventService(elasticSearch)

	router := mux.NewRouter()

	router.Handle("/api/1/archive",
		server.ApiF(appContext, server.ArchiveHandler))
	router.Handle("/api/1/event/{id}",
		server.ApiF(appContext, server.GetEventByIdHandler))
	router.Handle("/api/1/config",
		server.ApiF(appContext, server.ConfigHandler))
	router.Handle("/api/1/version",
		server.ApiF(appContext, server.VersionHandler))
	router.Handle("/api/1/eve2pcap", server.ApiF(appContext, server.Eve2PcapHandler))

	router.Handle("/api/1/inbox", server.ApiH(appContext, server.InboxHandler{}))

	router.Handle("/api/1/query", server.ApiF(appContext, server.QueryHandler))

	router.Handle("/api/1/_bulk", server.ApiF(appContext, server.EsBulkHandler))

	//
	// Disable the Elastic Search proxy for now to weed out its usage.
	//

	//// Elastic Search proxy.
	//esProxyHandler, err := elasticsearch.NewElasticSearchProxy(
	//	getElasticSearchUrl(), opts.NoCheckCertificate)
	//if err != nil {
	//	log.Fatal("Failed to initialize Elastic Search proxy: %v", err)
	//}
	//router.PathPrefix("/elasticsearch").Handler(http.StripPrefix("/elasticsearch", esProxyHandler))

	// Static file server, must be last as it serves as the fallback.
	router.PathPrefix("/").Handler(server.StaticHandlerFactory(opts.DevServerUri))

	log.Printf("Listening on %s:%s", opts.Host, opts.Port)
	err = http.ListenAndServe(opts.Host+":"+opts.Port, router)
	if err != nil {
		log.Fatal(err)
	}
}
