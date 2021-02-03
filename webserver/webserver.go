package webserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
	router := webs.newRouter(staticAssetsDir)
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
			log.Println(err)
		}
	}()

	go func() {
		<-webs.shutdown
		var wait time.Duration = time.Second * 15
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		fmt.Println("Shutting down web server")
		server.Shutdown(ctx)
		// Optionally, you could run srv.Shutdown in a goroutine and block on
		// <-ctx.Done() if your application should wait for other services
		// to finalize based on context cancellation.
		return
	}()
}

// Serve static files: https://golangcode.com/serve-static-assets-using-the-mux-router/
// staticAssetsDir: e.g.: "/static/"
func (webs *WebServer) newRouter(staticAssetsDir string) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.
		PathPrefix(staticAssetsDir).
		Handler(http.StripPrefix(staticAssetsDir, http.FileServer(http.Dir("."+staticAssetsDir))))
	return router
}

// https://golang-examples.tumblr.com/post/99458329439/get-local-ip-addresses
func (webs *WebServer) ListNetworkInterfaces() error {
	list, err := net.Interfaces()
	if err != nil {
		return err
	}

	for i, iface := range list {
		str := fmt.Sprintf(" %d %s: ", i+1, iface.Name)
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for j, addr := range addrs {
			// fmt.Printf(" %d %v\n", j, addr)
			str += fmt.Sprintf("%v", addr)
			if j < len(addrs) {
				str += ", "
			}
		}
		fmt.Println(str)
	}
	return nil
}
