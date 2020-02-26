package main

import (
	"fmt"
	"jongme/app/config"
	"jongme/app/database"
	"jongme/app/routers"
	"os"
	"os/signal"
	"syscall"

	"github.com/lab259/cors"
	"github.com/valyala/fasthttp"
)

func main() {

	mongoURL := "mongodb://" + config.MongoHost + ":" + config.MongoPort
	if config.MongoUserName != "" && config.MongoPassword != "" {
		mongoURL = "mongodb://" + config.MongoUserName + ":" + config.MongoPassword + "@" + config.MongoHost + ":" + config.MongoPort
	}
	db, err := database.New(mongoURL, "Jongme")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := routers.Create(db)
	cor := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Authorization", "X-Requested-With"},
		// OptionsPassthrough: true,
	})

	port := ":" + os.Getenv("PORT")

	srv := &fasthttp.Server{
		Handler: cor.Handler(r.Handler),
		// Handler: cors.Default().Handler(r.Handler),
	}

	go func() {
		if err := srv.ListenAndServe(port); err != nil {
			fmt.Println("ERROR : JONGME is listening on port =", port, "error =", err)
		}
	}()

	fmt.Println("JONGME is ready to listen and serve on port =", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	switch <-sigChan {
	case os.Interrupt:
		fmt.Println("\n Hot SIGINT...")
	case syscall.SIGTERM:
		fmt.Println("\n Got SIGTERM...")
	}

	fmt.Println("The service is shutting down...")
	os.Exit(0)
}
