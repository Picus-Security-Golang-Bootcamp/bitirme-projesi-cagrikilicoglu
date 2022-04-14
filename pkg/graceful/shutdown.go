package graceful

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Shutdown(srv *http.Server, timeout time.Duration) {
	c := make(chan os.Signal, 1)

	// when there is a interrupt signal, relay it to the channel
	signal.Notify(c, os.Interrupt)

	// block until any signal is received by the channel
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// wait until the timeout deadline and shutdown the server if there is no connections. if there is no connection shutdown immediately
	srv.Shutdown(ctx)

	log.Println("shutting down the server")
	os.Exit(0)
}
