package searcher

import (
	"fmt"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestSearcher_Search(t *testing.T) {
	type fields struct {
		FS fs.FS
	}
	type args struct {
		word string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantFiles []string
		wantErr   bool
	}{
		{
			name: "Ok",
			fields: fields{
				FS: fstest.MapFS{
					"file1.txt": {Data: []byte("World")},
					"file2.txt": {Data: []byte("World1")},
					"file3.txt": {Data: []byte("Hello World")},
				},
			},
			args:      args{word: "World"},
			wantFiles: []string{"file1", "file3"},
			wantErr:   false,
		},
		{
			name: "lil",
			fields: fields{
				FS: fstest.MapFS{
					"file2.txt": {Data: []byte("World1")},
				},
			},
			args:      args{word: "World"},
			wantFiles: nil,
			wantErr:   false,
		},
		{
			name: "Error",
			fields: fields{
				FS: fstest.MapFS{
					"file2.txt": {Data: []byte("World1")},
				},
			},
			args:      args{word: "World"},
			wantFiles: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Searcher{
				FS: tt.fields.FS,
			}
			if tt.wantErr {
				oldDirFilesFS := dirFilesFS
				defer func() { dirFilesFS = oldDirFilesFS }()
				dirFilesFS = func(fsys fs.FS, dir string) ([]string, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			gotFiles, err := s.Search(tt.args.word)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFiles, tt.wantFiles) {
				t.Errorf("Search() gotFiles = %v, want %v", gotFiles, tt.wantFiles)
			}
		})
	}
}
