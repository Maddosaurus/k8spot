package main

import (
	"context"
	"flag"
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

func main() {
	flag.StringVar(&flavorFlag, "flavor", "minikube", "Flavor to run. Currently: minikube|kubelet")
	flag.Parse()

	// FIXME: We need to set additional headers by adding a custom middleware
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		//Format: "${time_rfc3339_nano}${status}/${method} -> ${uri}\n",
		Format: "${uri} - ${method}\n",
	}))

	if flavorFlag == "minikube" {
		e.GET("/", minikube.RootHandler)
		e.HEAD("/", minikube.RootHandler)
		e.GET("/apis", minikube.APIsHandler)
		e.GET("/apis/:kind/:version", minikube.APIsGVKHandler)
		e.GET("/api", minikube.APIHandler)
		e.GET("/api/v1", minikube.APIv1Handler)
		e.GET("/api/:version/:resource", minikube.RuntimeResourceHandler)
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
	//FIXME: HTTPS breaks the kubectl config
	go func() {
		//if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		if err := e.StartTLS(":8080", "third_party/cert/apiserver.crt", "third_party/cert/apiserver.key"); err != nil && err != http.ErrServerClosed {
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
