package evereader

import (
	"fmt"
	"syscall"
	"os"
	"encoding/json"
)

type Bookmark struct {
	// The filename.
	Path   string `json:"path"`

	// The offset, for Eve this is the line number.
	Offset uint64 `json:"offset"`

	State  struct {
		       Inode uint64 `json:"inode"`
	       } `json:"state"`
}

type BookmarkedReader struct {
	reader       *EveReader

	bookmarkPath string
}

func NewBookmarkedReader(path string, bookmarkPath string) (*BookmarkedReader, error) {
	bookmarkedReader := &BookmarkedReader{}
	var err error

	bookmarkedReader.reader, err = New(path)
	if err != nil {
		return nil, err
	}

	if bookmarkPath == "" {
		bookmarkPath = fmt.Sprintf("%s.bookmark", path)
	}
	bookmarkedReader.bookmarkPath = bookmarkPath

	// Attempt to read the bookmark.
	currentBookmark, err := bookmarkedReader.ReadBookmark()
	if err == nil {
		err = bookmarkedReader.reader.SkipTo(currentBookmark.Offset)
	}

	// Attempt to write a bookmark,
	err = bookmarkedReader.WriteBookmark()
	if err != nil {
		bookmarkedReader.reader.Close()
		return nil, err
	}

	return bookmarkedReader, nil
}

func (r *BookmarkedReader) Close() {
	r.reader.Close()
}

func (r *BookmarkedReader) Next() (map[string]interface{}, error) {
	return r.reader.Next()
}

func (r *BookmarkedReader) IsActiveFile(bookmark *Bookmark) bool {

	if bookmark.Path != r.reader.path {
		return false;
	}

	fileInfo, err := r.reader.GetFileInfo()
	if err == nil {
		sysStat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if ok {
			if sysStat.Ino != bookmark.State.Inode {
				return false
			}
		}
	}

	// Looks like the same file.
	return true;
}

func (r *BookmarkedReader) ReadBookmark() (*Bookmark, error) {
	file, err := os.Open(r.bookmarkPath)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	var bookmark Bookmark
	err = decoder.Decode(&bookmark)
	if err != nil {
		return nil, err
	}
	return &bookmark, nil
}

func (r *BookmarkedReader) WriteBookmark() error {

	bookmark := Bookmark{}
	bookmark.Path = r.reader.path
	bookmark.Offset = r.reader.Pos()

	fileInfo, err := r.reader.GetFileInfo()
	if err == nil {
		sysStat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if ok {
			bookmark.State.Inode = sysStat.Ino
		}
	}

	file, err := os.Create(r.bookmarkPath)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(bookmark)
	if err != nil {
		return err
	}
	return nil
}