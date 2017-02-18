package server

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type FileStore struct {
	rootFolder  string
	host        string
	networkPath string
}

func NewFileStore(folder string, host string) *FileStore {
	fileStore := &FileStore{
		rootFolder:  folder,
		host:        host,
		networkPath: "/",
	}

	fs := http.FileServer(http.Dir(folder))
	http.Handle(fileStore.networkPath, fs)
	return fileStore
}

func (store FileStore) buildUri(token string) string {
	return store.host + store.networkPath + token
}

func (store FileStore) AddFile(fileContent string) string {
	now := time.Now().UnixNano()
	token := "tmp" + strconv.FormatInt(now, 10)
	filePath := store.rootFolder + "/" + token
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Failed to open and save", filePath, err)
		return ""
	}

	_, err = file.WriteString(string(fileContent))
	if err != nil {
		log.Println("Failed to write file", filePath, err)
		return ""
	}

	uri := store.buildUri(token)
	log.Println("File saved", uri)

	return uri
}

func (store FileStore) RemoveFile(uri string) bool {
	split := strings.Split(uri, "/")
	token := split[len(split)-1]
	filePath := store.rootFolder + "/" + token
	err := os.Remove(filePath)
	if err != nil {
		log.Println("Failed to remove file for uri", uri)
		return false
	}

	log.Printf("File with %s deleted", uri)
	return true
}
