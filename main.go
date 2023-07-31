package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"rooms/handler"
	"rooms/model"
	"rooms/repo"
	"rooms/service"
)

func main() {
	repo := repo.NewRepository(
		repo.WithPlayers(map[string]*model.Player{}),
		repo.WithRooms(map[string]*model.Room{}),
		repo.WithWaitingList(map[string]*model.Player{}),
	)
	serv := service.NewService(service.WithRepo(repo))
	handler := handler.NewHandler(handler.WithService(serv))

	mux := mux.NewRouter()
	mux.Handle("/register", handler.Register()).Methods("POST")
	mux.Handle("/stats", handler.Stats()).Methods("GET")

	mux.Handle("/websocket", handler.Websocket())

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	log.Println("server started successfully")

	stopC := make(chan os.Signal)
	signal.Notify(stopC, os.Interrupt)
	<-stopC

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("server stopping ...")
	defer cancel()

	log.Fatal(srv.Shutdown(ctx))
}
