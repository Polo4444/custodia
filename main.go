package main

//go:generate go-localize -input localizations_src -output localizations

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	muxHandlers "github.com/gorilla/handlers"
	"polo.gamesmania.io/custodia/aws_sdk"
	"polo.gamesmania.io/custodia/db"
	"polo.gamesmania.io/custodia/gp"
	"polo.gamesmania.io/custodia/handlers"
	"polo.gamesmania.io/custodia/spa"

	"bitbucket.org/polo44/goutilities"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const ConfigFileName = "config.yaml"
const AWSConfigFileName = "aws.yaml"

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// ─── INIT ────────────────────────────────────────────────────────────────────────
	gp.LogFatalIfErr(gp.Init(ConfigFileName))
	goutilities.DebugMode = gp.PConfig.DebugMode
	handlers.VersionType = &gp.PConfig.Environment

	// ─── Aws ─────────────────────────────────────────────────────────────
	gp.LogFatalIfErr(aws_sdk.Init(AWSConfigFileName))

	// ─── DB ────────────────────────────────────────────────────────────────────────
	gp.LogFatalIfErr(db.Connect())
	gp.LogFatalIfErr(db.Init(), "Can't init mongo database")

	// We check reconnection
	dbReconnectTicker := time.NewTicker(gp.PConfig.DBConnectionRefreshTicker)
	defer dbReconnectTicker.Stop()
	go func() {
		for range dbReconnectTicker.C {
			db.ReconnectCheck()
		}
	}()

	// ─── ROUTER ─────────────────────────────────────────────────────────────────────
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/version", handlers.VersionHandler)
	router.HandleFunc("/health", handlers.HealthHandler)

	// ─── SITE-DATA PANEL──────────────────────────────────────────────────────────────────
	subRouterAPIV1 := router.PathPrefix(handlers.APIV1Endpoint).Subrouter()
	handlers.InitAPIV1Routes(subRouterAPIV1)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "x-requested-with",
			"apikey", "user-token", "lang",
		},
		ExposedHeaders:   []string{},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: false,
	})

	spa := spa.Handler{StaticPath: "ui", IndexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)
	handler := c.Handler(router)

	srv := &http.Server{
		Handler:      muxHandlers.LoggingHandler(os.Stdout, handler),
		Addr:         ":" + gp.PConfig.HTTPPort,
		WriteTimeout: 10 * time.Minute,
		ReadTimeout:  2 * time.Minute,
	}

	go func() {
		log.Printf("%s (%s scope) Http Server running on port :%s", gp.PConfig.ProjectName, *handlers.VersionType, gp.PConfig.HTTPPort)
		log.Fatal(srv.ListenAndServe())
	}()

	<-sigs

	// ─── Shutdown Program ────────────────────────────────────────────────
	// Stop db reconnect ticker
	dbReconnectTicker.Stop()

	// Shutdown server
	log.Println("Shutting down the server...")
	srv.Shutdown(context.TODO())
}
