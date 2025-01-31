package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dyrakavalyanie/Clothes_shop/services/catalog/internal/config"
	"github.com/Dyrakavalyanie/Clothes_shop/services/catalog/internal/handlers"
	//"github.com/Dyrakavalyanie/Clothes_shop/services/catalog/internal/models"
	"github.com/Dyrakavalyanie/Clothes_shop/services/catalog/internal/scripts"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg := config.LoadConfig()

	connPool, err := scripts.CreateDBConnection(cfg)
	if err != nil {
		fmt.Println("error in creating DB connection:", err)
		return
	}

	err = scripts.CreateCatalogTables(connPool)
	if err != nil {
		fmt.Println("error in creating catalog tables:", err)
		return
	}

	mux := http.NewServeMux()

	handler := handlers.Handler{ConnPool : connPool}
	mux.HandleFunc("/item", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handler.AddItem(w, r)
		case "PATCH":
			handler.UpdateItem(w, r)
		case "DELETE":
			handler.DeleteItem(w, r)
		case "GET":
			handler.GetItem(w, r)
		}
	})

	httpServer := http.Server{
		Addr: ":" + cfg.CatalogPort,
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("err in listen: %s\n", err)
			return fmt.Errorf("failed to serve http server: %w", err)
		}
		fmt.Println("after listening")

		return nil
	})

	group.Go(func() error {
		fmt.Println("before ctx done")
		<-ctx.Done()
		fmt.Println("after ctx done")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			return err
		}
		fmt.Println("after server shutdown")

		return nil
	})

	err = group.Wait()
	if err != nil {
		fmt.Printf("after wait: %s\n", err)
		return
	}
}