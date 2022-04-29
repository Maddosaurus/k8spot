package kubelet

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/maddosaurus/k8spot/pkg/util"
)

var basepath = "third_party/kubelet/kubelet-"

func PodsHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, *util.ReadJson(basepath + "pods.json"))
}

func RunningPodsHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, *util.ReadJson(basepath + "runningpods.json"))
}

func ConfigzHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, *util.ReadJson(basepath + "configz.json"))
}

// func HealthzHandler(c echo.Context) error {
//	return c.JSON(http.StatusOK, *util.ReadJson(basepath + "healthz.json"))
// }

func MetricsHandler(c echo.Context) error {
	return c.String(http.StatusOK, util.ReadString(fmt.Sprintf("%vmetrics.txt", basepath)))
}

func RunHandler(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func DefaultHandler(c echo.Context) error {
	var data map[string]interface{}
	c.Logger().Warnf("Found unimplemented kubelet API call! %v", c.Request().URL)
	return c.JSON(http.StatusBadRequest, data)
}
