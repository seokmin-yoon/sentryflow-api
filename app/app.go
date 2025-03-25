package app

import (
	"log"
	"net/http"
	"sentryflow-api/app/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Run() {
	// Mux 라우터 생성
	r := mux.NewRouter()

	// 기본 핸들러 등록
	r.HandleFunc("/ping", handler.PingHandler).Methods("GET")
	r.HandleFunc("/users", handler.UsersHandler).Methods("GET")
	// api log
	r.HandleFunc("/api/logs", handler.APILogsHander).Methods("GET")
	r.HandleFunc("/api/logs", handler.FilterAPILogsHandler).Methods("POST")
	// envoy
	r.HandleFunc("/envoy/metrics", handler.EnvoyMetricsHander).Methods("GET")
	// cluster
	r.HandleFunc("/clusters", handler.ClustersHander).Methods("GET")
	//r.HandleFunc("/clusters/{cluster}", handler.ClustersHander).Methods("GET")
	//r.HandleFunc("/clusters/{cluster}/namespaces", handler.ClustersHander).Methods("GET")

	// CORS 미들웨어 적용
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                                       // 모든 도메인 허용
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // 허용할 HTTP 메서드
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),           // 허용할 헤더
	)

	log.Println("Server running on port 9090")
	log.Fatal(http.ListenAndServe(":9090", corsHandler(r))) // CORS 적용
}
