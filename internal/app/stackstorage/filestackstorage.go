package stackstorage

import (
	"fmt"
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
	dirName := filepath.Join("stacks", now.Format("2006-01-02"), model.Namespace)

	err := os.MkdirAll(dirName, os.ModePerm)

	if err != nil {
		return fmt.Errorf("mkdir '%s', error:%v", dirName, err)
	}

	fileName := fmt.Sprintf("%s-%s-%s.log", model.PodName, model.ContainerName, now.Format("15-04-05"))

	filePath := filepath.Join(dirName, fileName)

	err = os.WriteFile(filePath, []byte(model.Stack), 0644)

	if err != nil {
		return fmt.Errorf("write file '%s', error:%v", filePath, err)
	}

	return nil
}
