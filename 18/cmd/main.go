package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
	"l2.18/internal/handler"
	"l2.18/internal/repository/memory"
	"l2.18/internal/service/events"
	"l2.18/pkg/server"
)

func main() {
	port := flag.String("port", "8000", "Port to run the server on")
	flag.Parse()

	repo := memory.NewEventsRepository()
	service := events.New(repo)
	eventsHandler := handler.NewEventsHandler(service)

	handlerLogger := slog.New(slog.NewTextHandler(
		os.Stdout, &slog.HandlerOptions{}).WithGroup("handler"))
	middleware := handler.NewMiddleware(handlerLogger)

	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", middleware.Logging(eventsHandler.CreateEvent))
	mux.HandleFunc("/update_event", middleware.Logging(eventsHandler.UpdateEvent))
	mux.HandleFunc("/delete_event", middleware.Logging(eventsHandler.DeleteEvent))
	mux.HandleFunc("/events_for_day", middleware.Logging(eventsHandler.EventsForDay))
	mux.HandleFunc("/events_for_week", middleware.Logging(eventsHandler.EventsForWeek))
	mux.HandleFunc("/events_for_month", middleware.Logging(eventsHandler.EventsForMonth))

	ctx, cancel := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	srv := server.New(*port, mux)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error { return srv.Run() })
	g.Go(func() error {
		<-gCtx.Done()
		return srv.Shutdown(context.Background())
	})

	fmt.Printf("server started on port %s\n", *port)

	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("exit with %v\n", err)
	}

	fmt.Println("server was shut down")
}
