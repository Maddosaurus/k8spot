package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadJson(name string) *map[string]interface{} {
	var data map[string]interface{}
	path, err := filepath.Abs(fmt.Sprintf("./%v", name))
	if err == nil {
		json_file, jerr := os.Open(path)
		byteVal, _ := ioutil.ReadAll(json_file)
		if jerr == nil {
			json.Unmarshal(byteVal, &data)
		}
	}
	return &data
}

func ReadString(name string) string {
	path, err := filepath.Abs(fmt.Sprintf("./%v", name))
	if err == nil {
		textFile, terr := ioutil.ReadFile(path)
		if terr == nil {
			return string(textFile)
		}
	}
	return ""
}
