package utils

import (
	"io/ioutil"
	"os"
)

func OpenJson(name string) []byte {
	jsonFile, _ := os.Open(name)
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}
