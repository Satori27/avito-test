package main

import (
	// "fmt"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gitlab.com/Satori27/avito/internal/config"
	bidfeedback "gitlab.com/Satori27/avito/internal/http-server/handlers/bids/bid_feedback"
	bidsubmitdecision "gitlab.com/Satori27/avito/internal/http-server/handlers/bids/bid_submit_decision"
	bidsrollback "gitlab.com/Satori27/avito/internal/http-server/handlers/bids/bids_rollback"
	getbidstatus "gitlab.com/Satori27/avito/internal/http-server/handlers/bids/get_bid_status"
	getbids "gitlab.com/Satori27/avito/internal/http-server/handlers/bids/get_bids"
	getmybids "gitlab.com/Satori27/avito/internal/http-server/handlers/bids/get_my_bids"
	getreviews "gitlab.com/Satori27/avito/internal/http-server/handlers/bids/get_reviews"
	newbid "gitlab.com/Satori27/avito/internal/http-server/handlers/bids/new"
	patchbid "gitlab.com/Satori27/avito/internal/http-server/handlers/bids/patch_bid"
	putbidstatus "gitlab.com/Satori27/avito/internal/http-server/handlers/bids/put_bid-status"
	"gitlab.com/Satori27/avito/internal/http-server/handlers/ping"
	getmytenders "gitlab.com/Satori27/avito/internal/http-server/handlers/tender/get_my_tenders"
	gettenderstatus "gitlab.com/Satori27/avito/internal/http-server/handlers/tender/get_tender_status"
	gettenders "gitlab.com/Satori27/avito/internal/http-server/handlers/tender/get_tenders"
	new_tender "gitlab.com/Satori27/avito/internal/http-server/handlers/tender/new"
	patchtenderstatus "gitlab.com/Satori27/avito/internal/http-server/handlers/tender/patch_tender_status"
	puttenderstatus "gitlab.com/Satori27/avito/internal/http-server/handlers/tender/put_tender_status"
	tendersrollback "gitlab.com/Satori27/avito/internal/http-server/handlers/tender/tenders_rollback"
	psq "gitlab.com/Satori27/avito/internal/storage/postgres"
)

func main() {
	cfg := config.Load()
	storage := &psq.Storage{}

	ctx, cancel := context.WithCancel(context.Background())
	log := setuplogger()
	go func() {
		err := psq.New(cancel, storage, cfg)
		if err != nil {
			log.Error("failed to init storage", slog.String("error", err.Error()))
		}
	}()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/api", func(r chi.Router) {
		r.Route("/tenders", func(r chi.Router) {
			r.Post("/new", new_tender.New(storage))
			r.Get("/{tenderId}/status", gettenderstatus.New(storage))
			r.Put("/{tenderId}/status", puttenderstatus.New(storage))
			r.Patch("/{tenderId}/edit", patchtenderstatus.New(storage))
			r.Get("/", gettenders.New(storage))
			r.Get("/my", getmytenders.New(storage))
			r.Put("/{tenderId}/rollback/{version}", tendersrollback.New(storage))

		})

		r.Route("/bids", func(r chi.Router) {
			r.Post("/new", newbid.New(storage))
			r.Get("/{bidId}/status", getbidstatus.New(storage))
			r.Put("/{bidId}/status", putbidstatus.New(storage))
			r.Patch("/{bidId}/edit", patchbid.New(storage))
			r.Get("/my", getmybids.New(storage))
			r.Get("/{tenderId}/list", getbids.New(storage))
			r.Put("/{bidId}/submit_decision", bidsubmitdecision.New(storage))
			r.Put("/{bidId}/feedback", bidfeedback.New(storage))
			r.Get("/{tenderId}/reviews", getreviews.New(storage))
			r.Put("/{bidId}/rollback/{version}", bidsrollback.New(storage))

		})

		r.Get("/ping", ping.New(ctx))
	})

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  time.Minute,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancelCtx := context.WithTimeout(serverCtx, 15*time.Second)
		defer cancelCtx()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				panic("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		log.Info("server shut down")
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			panic(err.Error())
		}
		serverStopCtx()
	}()

	// Run the server
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err.Error())
	}

	// Wait for server context to be stopped

	<-serverCtx.Done()

	log.Info("server stopped")

}

func setuplogger() *slog.Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
