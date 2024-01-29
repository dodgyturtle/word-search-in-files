package searcher

import (
	"bufio"
	"io/fs"
	"sort"
	"strings"
	"sync"
	"word-search-in-files/pkg/internal/dir"
)

type Searcher struct {
	FS fs.FS
}

var dirFilesFS = dir.FilesFS

func findInFiles(fs fs.FS, fileName string, word string) (*string, error) {
	postFile, err := fs.Open(fileName)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(postFile)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		if word == scanner.Text() {
			fileName = strings.TrimRight(fileName, ".txt")
			return &fileName, nil
		}
	}
	return nil, nil
}

func (s *Searcher) Search(word string) (files []string, err error) {
	filesName, err := dirFilesFS(s.FS, "")
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var lock sync.Mutex

	filesNameLen := len(filesName)
	wg.Add(filesNameLen)

	errCh := make(chan error, filesNameLen)

	for _, fileName := range filesName {
		go func(fileName string) {
			defer wg.Done()
			gotFile, err := findInFiles(s.FS, fileName, word)
			if err != nil {
				errCh <- err
			}
			if gotFile != nil {
				lock.Lock()
				files = append(files, *gotFile)
				lock.Unlock()
			}

		}(fileName)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()
	for err := range errCh {
		return nil, err
	}
	if len(files) == 0 {
		return nil, nil
	}
	sort.Strings(files)
	return files, nil
}
