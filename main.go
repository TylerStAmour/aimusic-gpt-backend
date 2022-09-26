package main

import (
	"aimusic_gpt_backend/router"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr: "http://127.0.0.1:8080",
		Handler: router.Get(),
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
