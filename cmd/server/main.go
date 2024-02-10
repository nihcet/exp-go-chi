package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	config "github.com/nihcet/exp-go-chi/configs"
	logutil "github.com/nihcet/go-lib/pkg/util/log"
)

func main() {
	// parse variable from config
	config.LoadConfig()
	serverConfig := config.Get()
	serviceName := serverConfig.ServiceName
	serviceHost := serverConfig.ServiceHost
	servicePort := serverConfig.ServicePort

	// initial log
	logutil.InitializeLog(serviceName)
	log := logutil.GetLogger()

	ctx := context.Background()
	r := chi.NewRouter()

	// add middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// set timeout
	r.Use(middleware.Timeout(10 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	serverAddr := fmt.Sprintf("%s:%s", serviceHost, servicePort)
	server := &http.Server{Addr: serverAddr, Handler: r}

	go func() {
		log.Infof("server is running at: %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf(err, "error to listen on server")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	server.Shutdown(ctx)
	log.Info("shutdown")
}
