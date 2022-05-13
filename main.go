package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/maddosaurus/k8spot/pkg/handler/kubelet"
	"github.com/maddosaurus/k8spot/pkg/handler/minikube"
)

var (
	flavorFlag string
)

// Add additional k8s headers
func K8sHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		c.Response().Header().Set("X-Kubernetes-Pf-Flowschema-Uid", "cc104973-c107-4c5a-bfb7-f2f010a678cd")
		c.Response().Header().Set("X-Kubernetes-Pf-Prioritylevel-Uid", "4ee6eada-c507-4632-9ed9-cc837d0acf0b")
		return next(c)
	}
}

func main() {
	flag.StringVar(&flavorFlag, "flavor", "minikube", "Flavor to run. Currently: minikube|kubelet")
	flag.Parse()

	run_port := os.Getenv("K8SPOT_PORT")
	if run_port == "" {
		run_port = "8080"
	}
	bind_string := fmt.Sprintf(":%v", run_port)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Add k8s headers for API Prio & Flow control
	// https://kubernetes.io/docs/concepts/cluster-administration/flow-control/
	e.Use(K8sHeader)

	// Set up logging
	f, err := os.OpenFile("./log/k8spot.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()

	// Set up the logger object as well as the middleware
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetOutput(f)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		//Format: "${uri} - ${method} - ${user_agent} - ${remote_ip}\n",
		Output: f,
		Format: `{"time": "${time_rfc3339_nano}","remote_ip":"${remote_ip}", "host":"k8spot",` +
			`"method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}"}` + "\n",
	}))

	if flavorFlag == "minikube" {
		e.GET("/", minikube.RootHandler)
		e.HEAD("/", minikube.RootHandler)
		e.GET("/apis", minikube.APIsHandler)
		e.GET("/apis/:kind/:version", minikube.APIsGVKHandler)
		e.GET("/api", minikube.APIHandler)
		e.GET("/api/v1", minikube.APIv1Handler)
		e.GET("/api/:version/:resource", minikube.RuntimeResourceHandler)
		e.GET("/api/v1/nodes", minikube.NodeHandler)
		e.GET("/api/v1/namespaces/:namespace/pods", minikube.PodHandler)
		e.GET("/*", minikube.DefaultHandler)
	} else if flavorFlag == "kubelet" {
		// FIXME: Sync the two run modes so they produce coherent output
		e.GET("/pods", kubelet.PodsHandler)
		e.GET("/runningpods", kubelet.RunningPodsHandler)
		e.GET("/configz", kubelet.ConfigzHandler)
		//e.GET("/healthz", kubelet.HealthzHandler)
		e.GET("/metrics", kubelet.MetricsHandler)
		e.GET("/*", kubelet.DefaultHandler)
		e.POST("/run/:namespace/:pod/:container", kubelet.RunHandler)

	}

	//FIXME: This should be different for the kubelet
	go func() {
		if err := e.StartTLS(bind_string, "third_party/cert/apiserver.crt", "third_party/cert/apiserver.key"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
