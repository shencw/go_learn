package main

import (
	"context"
	"go_learn/internal/apiserver"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

func main() {
	insecureServer := &http.Server{
		Addr:         ":8080",
		Handler:      apiserver.Route(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		err := insecureServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Println("insecureServer:", err)
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		log.Println("结束err", err)
	}

	log.Println("结束")
}
