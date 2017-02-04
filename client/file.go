package client

import "errors"

type File struct {
    Name string
    data []byte
}

func (f *File) Read(count int64) ([]byte, error) {
    return []byte{}, errors.New("Not Implemented")
}

func (f *File) Save(filePath string) error {
    return errors.New("Not Implemented")
}

func (f *File) IsEOF() bool {
    return true
}

func (f *File) Seek(position int64) error {
    return errors.New("Not Implemented")
}
