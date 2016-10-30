package agent

import (
	flag "github.com/spf13/pflag"
	"os"
	"github.com/jasonish/evebox/log"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	golog "log"
	"strings"
	"gopkg.in/yaml.v2"
	"github.com/jasonish/evebox/agent"
	"sync"
)

type Flags struct {
	configFilename string
	verbose bool
}

type AgentConfig map[string]interface{}

func (c AgentConfig) GetMap(key string) AgentConfig {
	return AgentConfig(c[key].(map[string]interface{}))
}

func (c AgentConfig) Get(key string) interface{} {
	return c[key]
}

func (c AgentConfig) GetString(key string) string {
	return c[key].(string)
}

func (c AgentConfig) GetStringList(key string) []string {
	output := make([]string, 0)

	for _, value := range c[key].([]interface{}) {
		output = append(output, value.(string))
	}

	return output
}

func Main(args []string) {

	golog.SetFlags(golog.Llongfile)

	flags := Flags{}

	flagset := flag.NewFlagSet("evebox agent", flag.ExitOnError)

	flagset.StringVarP(&flags.configFilename, "config", "c", "", "Configuration file name")
	flagset.BoolVarP(&flags.verbose, "verbose", "v", false, "Verbose output")

	flagset.Parse(args)

	if flags.verbose {
		log.SetLevel(log.DEBUG)
	}

	config := AgentConfig{}

	if flags.configFilename != "" {
		loadConfig(flags.configFilename, config)
	} else {
		log.Fatal("Please specify a configuration file.")
	}

	var wg sync.WaitGroup

	inputConfig := config.GetMap("input")
	bookmarkDir := inputConfig.GetString("bookmark-dir")
	for _, path := range inputConfig.GetStringList("paths") {
		logReader := agent.NewLogReader(path, bookmarkDir)
		logReader.Init()
		wg.Add(1)
		runReader(wg, logReader)
	}

	log.Println("Waiting for goroutinges to finish.")
	wg.Wait()
}

func loadConfig(filename string, config AgentConfig) {
	configFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasSuffix(filename, ".toml") {
		_, err = toml.Decode(string(configBytes), &config)
		if err != nil {
			log.Error("Failed to load configuration file:")
			log.Fatal(err)
		}
	} else if strings.HasSuffix(filename, ".yaml") {
		var wrapper map[string]interface{}
		yaml.Unmarshal(configBytes, &wrapper)
		for k, v := range wrapper["agent"].(map[interface{}]interface{}) {
			switch k := k.(type) {
			case string:
				config[k] = v
			}
		}
	}

	log.Println(config)
}

func runReader(wg sync.WaitGroup, reader *agent.LogReader) {
	go func() {
		reader.Run()
		wg.Done()
	}()
}