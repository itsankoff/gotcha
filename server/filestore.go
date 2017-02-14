package server

import (
    "github.com/itsankoff/gotcha/common"
)

type FileStore struct {

}

func (f FileStore) AddFile(msg *common.Message) string {
    return ""
}

func (f FileStore) RemoveFile(token string) bool {
    return false
}
