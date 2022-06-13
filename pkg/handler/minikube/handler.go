package minikube

import (
	"fmt"
	"net/http"

	"github.com/maddosaurus/k8spot/pkg/util"

	"github.com/labstack/echo/v4"
)

var basepath = "third_party/minikube/v1.23.3/"

func APIHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, *util.ReadJson(basepath + "api.json"))
}

func APIv1Handler(c echo.Context) error {
	return c.JSON(http.StatusOK, *util.ReadJson(basepath + "api-v1.json"))
}

func APIsHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, *util.ReadJson(basepath + "apis.json"))
}

func APIsGVKHandler(c echo.Context) error {
	kind := c.Param("kind")
	version := c.Param("version")
	json := util.ReadJson(fmt.Sprintf("%vapis/%v-%v.json", basepath, kind, version))
	if len(*json) == 0 {
		c.Logger().Warn(fmt.Sprintf("Found unimplemented object: %v/%v  -  ", kind, version), json)
	}
	return c.JSON(http.StatusOK, *json)

}

func RootHandler(c echo.Context) error {
	// This is what a secured k8s would return, a 403. When authorized, it will return the paths.
	//return c.JSON(http.StatusForbidden, *util.ReadJson(basepath + "root.json"))
	return c.JSON(http.StatusOK, *util.ReadJson(basepath + "paths.json"))
}

func RuntimeResourceHandler(c echo.Context) error {
	version := c.Param("version")
	resource := c.Param("resource")
	return c.JSON(http.StatusOK, *util.ReadJson(fmt.Sprintf("%v%v-%v.json", basepath, version, resource)))
}

func NodeHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, *util.ReadJson(basepath + "v1-nodes.json"))
}

func PodHandler(c echo.Context) error {
	//namespace := c.Param("namespace")
	return c.JSON(http.StatusOK, *util.ReadJson(basepath + "v1-pods.json"))
}

func DefaultHandler(c echo.Context) error {
	var data map[string]interface{}
	c.Logger().Warnf("Found unimplemented minikube API call! %v", c.Request().URL)
	return c.JSON(http.StatusBadRequest, data)
}
