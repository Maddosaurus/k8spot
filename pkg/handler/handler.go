package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func readJson(name string) *map[string]interface{} {
	var data map[string]interface{}
	path, err := filepath.Abs(fmt.Sprintf("./third_party/kubernetes/v1.23.3/%v", name))
	if err == nil {
		json_file, jerr := os.Open(path)
		byteVal, _ := ioutil.ReadAll(json_file)
		if jerr == nil {
			json.Unmarshal(byteVal, &data)
		}
	}
	return &data
}

func APIsHandler(c echo.Context) error {
	c.Logger().Info("I've been called!")
	return c.JSON(http.StatusOK, *readJson("apis.json"))
}

func DefaultHandler(c echo.Context) error {
	var data map[string]interface{}

	return c.JSON(http.StatusOK, data)
}
