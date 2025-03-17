package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"sentryflow-api/app/model"
	"sentryflow-api/config"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(users)
}

func APILogsHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Printf("[Request] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	client, err := config.GetMongoClient()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		log.Printf("[MongoDB] Connection error: %v", err)
		return
	}

	collection := client.Database("SentryFlow").Collection("APILogs")

	// Get API Logs from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch API Logs", http.StatusInternalServerError)
		log.Printf("[MongoDB] Query error: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var logs []model.APILog
	for cursor.Next(ctx) {
		var logEntry model.APILog
		if err := cursor.Decode(&logEntry); err != nil {
			http.Error(w, "Error decoding API Logs", http.StatusInternalServerError)
			return
		}
		logs = append(logs, logEntry)
	}

	log.Printf("[Response] Returning %d logs", len(logs))

	json.NewEncoder(w).Encode(logs)
}

func NamespaceAPILogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 요청을 받았다는 로그 출력
	log.Printf("[Request] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	// MongoDB 클라이언트 가져오기
	client, err := config.GetMongoClient()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		log.Printf("[MongoDB] Connection error: %v", err)
		return
	}

	collection := client.Database("SentryFlow").Collection("APILogs")

	// 경로 변수에서 네임스페이스 추출
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// 쿼리 파라미터에서 `type` 확인
	queryType := r.URL.Query().Get("type")

	var filter bson.M

	if queryType == "src" {
		filter = bson.M{"srcnamespace": namespace}
		log.Printf("[Filter] Fetching logs where srcnamespace = %s", namespace)
	} else if queryType == "dst" {
		filter = bson.M{"dstnamespace": namespace}
		log.Printf("[Filter] Fetching logs where dstnamespace = %s", namespace)
	} else {
		// 기본적으로 srcnamespace 또는 dstnamespace 중 하나라도 포함된 로그 조회
		filter = bson.M{
			"$or": []bson.M{
				{"srcnamespace": namespace},
				{"dstnamespace": namespace},
			},
		}
		log.Printf("[Filter] Fetching all logs where srcnamespace or dstnamespace = %s", namespace)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Failed to fetch API Logs", http.StatusInternalServerError)
		log.Printf("[MongoDB] Query error: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var logs []model.APILog
	for cursor.Next(ctx) {
		var logEntry model.APILog
		if err := cursor.Decode(&logEntry); err != nil {
			http.Error(w, "Error decoding API Logs", http.StatusInternalServerError)
			log.Printf("[MongoDB] Decoding error: %v", err)
			return
		}
		logs = append(logs, logEntry)
	}

	log.Printf("[Response] Returning %d logs for namespace: %s", len(logs), namespace)

	json.NewEncoder(w).Encode(logs)
}

func EnvoyMetricsHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Printf("[Request] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	client, err := config.GetMongoClient()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		log.Printf("[MongoDB] Connection error: %v", err)
		return
	}

	collection := client.Database("SentryFlow").Collection("EnvoyMetrics")

	// Get API Logs from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch Envoy Metrics", http.StatusInternalServerError)
		log.Printf("[MongoDB] Query error: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var metrics []model.EnvoyMetrics
	for cursor.Next(ctx) {
		var logEntry model.EnvoyMetrics
		if err := cursor.Decode(&logEntry); err != nil {
			http.Error(w, "Error decoding Envoy Metrics", http.StatusInternalServerError)
			return
		}
		metrics = append(metrics, logEntry)
	}

	log.Printf("[Response] Returning %d logs", len(metrics))

	json.NewEncoder(w).Encode(metrics)
}

func ClustersHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Printf("[Request] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	// Mux 변수 추출
	vars := mux.Vars(r)
	clusterName := vars["name"]

	// 요청 URL에 따라 처리
	if clusterName == "" {
		// GET /clusters → 모든 클러스터 반환
		log.Printf("[Response] Returning %d clusters", len(clusters))
		json.NewEncoder(w).Encode(clusters)
		return
	}

	// 특정 클러스터 찾기
	for _, cluster := range clusters {
		if cluster.Name == clusterName {
			// GET /clusters/{name}/namespaces 네임스페이스 정보만 반환
			if r.URL.Path == "/clusters/"+clusterName+"/namespaces" {
				log.Printf("[Response] Returning namespaces for cluster %s", clusterName)
				json.NewEncoder(w).Encode(cluster.Namespaces)
				return
			}

			// GET /clusters/{name} 클러스터 전체 정보 반환
			log.Printf("[Response] Returning details for cluster %s", clusterName)
			json.NewEncoder(w).Encode(cluster)
			return
		}
	}

	// 클러스터를 찾지 못한 경우
	http.Error(w, "Cluster not found", http.StatusNotFound)
}
