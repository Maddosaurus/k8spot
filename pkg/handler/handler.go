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
	path, err := filepath.Abs(fmt.Sprintf("./third_party/minikube/v1.23.3/%v", name))
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
	return c.JSON(http.StatusOK, *readJson("apis.json"))
}

func APIsGVKHandler(c echo.Context) error {
	kind := c.Param("kind")
	version := c.Param("version")
	if kind == "node.k8s.io" {
		return c.JSON(http.StatusOK, *readJson(fmt.Sprintf("apis/%v-%v.json", kind, version)))
	}

	return c.JSON(http.StatusOK, "")
}

func RuntimeResourceHandler(c echo.Context) error {
	version := c.Param("version")
	resource := c.Param("resource")
	return c.JSON(http.StatusOK, *readJson(fmt.Sprintf("%v-%v.json", version, resource)))
}

func DefaultHandler(c echo.Context) error {
	var data map[string]interface{}

	return c.JSON(http.StatusOK, data)
}
