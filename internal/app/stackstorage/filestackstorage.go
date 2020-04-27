package stackstorage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type FileStackStorage struct {
}

func NewFileStackStorage() FileStackStorage {
	return FileStackStorage{}
}

func (s FileStackStorage) Store(model ContainerStackModel) error {
	now := time.Now()
	dirName := filepath.Join("stacks", model.Namespace, now.Format("2006-01-02-15"))

	err := os.MkdirAll(dirName, os.ModePerm)

	if err != nil {
		return fmt.Errorf("mkdir '%s', error:%v", dirName, err)
	}

	fileName := fmt.Sprintf("%s-%s-%s.log", model.PodName, model.ContainerName, now.Format("15-04-05"))

	filePath := filepath.Join(dirName, fileName)

	err = ioutil.WriteFile(filePath, []byte(model.Stack), 0644)

	if err != nil {
		return fmt.Errorf("write file '%s', error:%v", filePath, err)
	}

	return nil
}
