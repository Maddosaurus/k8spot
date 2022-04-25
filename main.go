package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/maddosaurus/k8spot/pkg/handler"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	//e.Use(middleware.Logger())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		// Format: "${time_rfc3339_nano}${status}/${method} -> ${uri}\n",
		Format: "${uri}\n",
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/apis/:kind/:version", handler.APIsGVKHandler)
	e.GET("/apis", handler.APIsHandler)
	e.GET("/api/:version/:resource", handler.RuntimeResourceHandler)
	e.GET("/*", handler.DefaultHandler)

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
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
