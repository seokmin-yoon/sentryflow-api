package app

import (
	"log"
	"net/http"
	"sentryflow-api/app/handler"
)

func Run() {
	http.HandleFunc("/ping", handler.PingHandler)
	http.HandleFunc("/users", handler.UsersHandler)
	http.HandleFunc("/api/logs", handler.APILogHander)

	log.Println("Server running on port 9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
