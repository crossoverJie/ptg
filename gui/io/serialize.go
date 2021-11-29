package io

import (
	"io/ioutil"
	"os"
	"os/user"
)

func SaveLog(data []byte) error {
	filename, err := initLog()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0666)
}

func initLog() (string, error) {
	home, err := user.Current()
	if err != nil {
		return "", err
	}

	filename := home.HomeDir + "/.ptg/ptg.log"

	if !exist(filename) {
		err := os.MkdirAll(home.HomeDir+"/.ptg/", 0777)
		if err != nil {
			return "", err
		}
		_, err = os.Create(filename)
		if err != nil {
			return "", err
		}
	}
	return filename, nil
}

func LoadLog() ([]byte, error) {
	filename, err := initLog()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(filename)
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
