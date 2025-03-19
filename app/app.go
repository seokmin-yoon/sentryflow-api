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
	// api log
	r.HandleFunc("/api/logs", handler.APILogsHander).Methods("GET")
	r.HandleFunc("/api/logs/namespaces/{namespace}", handler.NamespaceAPILogsHandler).Methods("GET")
	// envoy
	r.HandleFunc("/envoy/metrics", handler.EnvoyMetricsHander).Methods("GET")
	// cluster
	r.HandleFunc("/clusters", handler.ClustersHander).Methods("GET")
	r.HandleFunc("/clusters/{cluster}", handler.ClustersHander).Methods("GET")
	r.HandleFunc("/clusters/{cluster}/namespaces", handler.ClustersHander).Methods("GET")
	//r.HandleFunc("/clusters/{cluster}/namespaces/{namespace}/pods", handler.ClustersHander).Methods("GET")

	log.Println("Server running on port 9090")
	log.Fatal(http.ListenAndServe(":9090", r))
}
