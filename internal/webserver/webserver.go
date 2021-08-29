package webserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Route struct {
	Pattern string
	Handler http.HandlerFunc
}

type WebServer struct {
	address  string
	shutdown chan bool
}

func NewWebServer(address string, shutdown chan bool) *WebServer {
	server := &WebServer{
		address:  address,
		shutdown: shutdown,
	}
	return server
}

func (webs *WebServer) Run(routes []Route, staticAssetsDir string) {
	// Serve static files: https://golangcode.com/serve-static-assets-using-the-mux-router/
	router := mux.NewRouter().StrictSlash(true)
	router.
		PathPrefix(staticAssetsDir).
		Handler(http.StripPrefix(staticAssetsDir, http.FileServer(http.Dir("."+staticAssetsDir))))

	for _, route := range routes {
		router.HandleFunc(route.Pattern, route.Handler)
	}

	server := &http.Server{
		Addr:         webs.address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	go func() {
		<-webs.shutdown
		var wait time.Duration = time.Second * 15
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		log.Info("Shutting down web server")
		server.Shutdown(ctx)
		// Optionally, you could run srv.Shutdown in a goroutine and block on
		// <-ctx.Done() if your application should wait for other services
		// to finalize based on context cancellation.
	}()
}
