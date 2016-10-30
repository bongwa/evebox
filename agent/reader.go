package agent

import (
	"crypto/md5"
	"fmt"
	"github.com/jasonish/evebox/evereader"
	"github.com/jasonish/evebox/log"
	"path"
	"io"
	"time"
)

type LogReader struct {
	Path        string
	BookmarkDir string

	bookmarkPath string
	bookmarker   *evereader.Bookmarker
	reader       *evereader.EveReader
}

func NewLogReader(path string, bookmarkDir string) *LogReader {
	return &LogReader{
		Path:        path,
		BookmarkDir: bookmarkDir,
	}
}

func (lr *LogReader) Init() error {

	var err error

	log.Info("Initializing reader for %s.", lr.Path)

	if lr.BookmarkDir == "" {
		lr.bookmarkPath = fmt.Sprintf("%s.bookmark", lr.Path)
	} else {
		lr.bookmarkPath = path.Join(lr.BookmarkDir,
			fmt.Sprintf("%x.bookmark", md5.Sum([]byte(lr.Path))))
	}

	log.Info("Using bookmark path of %s", lr.bookmarkPath)

	lr.reader, err = evereader.New(lr.Path)
	if err != nil {
		log.Fatal(err)
	}

	lr.bookmarker = &evereader.Bookmarker{
		Filename: lr.bookmarkPath,
		Reader: lr.reader,
	}
	err = lr.bookmarker.Init(true)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (lr *LogReader) Run() error {

	for {
		eof := false
		event, err := lr.reader.Next()
		if err != nil {
			if err == io.EOF {
				eof = true
			} else {
				return err
			}
		}

		if event != nil {
			log.Println(event)

			bookmark := lr.bookmarker.GetBookmark()
			lr.bookmarker.WriteBookmark(bookmark)

		} else if eof {
			log.Println("Got eof, sleeping.")
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}
