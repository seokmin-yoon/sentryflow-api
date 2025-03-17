package app

import (
	"log"
	"net/http"
	"sentryflow-api/app/handler"

	"github.com/gorilla/mux"
)

func Run() {
	// Mux 라우터 생성
	r := mux.NewRouter()

	// 기본 핸들러 등록
	// test
	r.HandleFunc("/ping", handler.PingHandler).Methods("GET")
	r.HandleFunc("/users", handler.UsersHandler).Methods("GET")
	// api
	r.HandleFunc("/api/logs", handler.APILogsHander).Methods("GET")
	r.HandleFunc("/api/logs/{namespace}", handler.NamespaceAPILogsHandler).Methods("GET")
	// envoy
	r.HandleFunc("/envoy/metrics", handler.EnvoyMetricsHander).Methods("GET")
	// cluster
	r.HandleFunc("/clusters", handler.ClustersHander).Methods("GET")
	r.HandleFunc("/clusters/{name}", handler.ClustersHander).Methods("GET")
	r.HandleFunc("/clusters/{name}/namespaces", handler.ClustersHander).Methods("GET")

	log.Println("Server running on port 9090")
	log.Fatal(http.ListenAndServe(":9090", r))
}
