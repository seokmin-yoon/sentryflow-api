package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"sentryflow-api/app/model"
	"sentryflow-api/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(users)
}

// 1. src, dst cluster name 추가
// 2. Unknown 일 경우 Services 에서 정보 매핑 후 수정
func APILogsHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Printf("[Request] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	client, err := config.GetMongoClient()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		log.Printf("[MongoDB] Connection error: %v", err)
		return
	}

	apiLogsCol := client.Database("SentryFlow").Collection("APILogs")
	servicesCol := client.Database("SentryFlow").Collection("Services")
	podsCol := client.Database("SentryFlow").Collection("Pods")

	// Get API Logs from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // timeout 5s
	defer cancel()

	// 현재 시간 기준 5분 전 타임스탬프 (Unix 기준, string으로 변환)
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute).Unix()
	filter := bson.M{
		"timestamp": bson.M{"$gte": strconv.FormatInt(fiveMinutesAgo, 10)},
	}

	cursor, err := apiLogsCol.Find(ctx, filter)
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

		// cluster name 채우기
		logEntry.SrcCluster = findCluster(ctx, logEntry.SrcType, logEntry.SrcName, logEntry.SrcNamespace, podsCol, servicesCol)
		logEntry.DstCluster = findCluster(ctx, logEntry.DstType, logEntry.DstName, logEntry.DstNamespace, podsCol, servicesCol)

		// Unknown 채우기
		fillMissingSourceInfo(ctx, &logEntry, servicesCol)
		fillMissingDestInfo(ctx, &logEntry, servicesCol)

		// timestamp 변환 (string -> RFC3339)
		logEntry.TimeStamp = convertTimestamp(logEntry.TimeStamp)
		logs = append(logs, logEntry)
	}

	log.Printf("[Response] Returning %d logs", len(logs))

	json.NewEncoder(w).Encode(logs)
}

func findCluster(ctx context.Context, resourceType, name, namespace string, podsCol, servicesCol *mongo.Collection) string {
	var result struct {
		Cluster string `bson:"cluster"`
	}
	var err error

	filter := bson.M{"name": name, "namespace": namespace}

	switch resourceType {
	case "Pod":
		err = podsCol.FindOne(ctx, filter).Decode(&result)
	case "Service":
		err = servicesCol.FindOne(ctx, filter).Decode(&result)
	default:
		log.Printf("[Cluster Lookup] Unknown resource type: %s", resourceType)
		return "cluster1"
	}

	if err != nil {
		log.Printf("[Cluster Lookup] Not found: %s/%s (%s): %v", namespace, name, resourceType, err)
		return "cluster1"
	}

	return result.Cluster
}

func fillMissingSourceInfo(ctx context.Context, logEntry *model.APILog, servicesCol *mongo.Collection) bool {
	if logEntry.SrcNamespace != "Unknown" {
		return true
	}
	cursor, err := servicesCol.Find(ctx, bson.M{"loadbalancerips": bson.M{"$in": []string{logEntry.SrcIP}}})
	if err != nil {
		log.Printf("[MongoDB] Error checking services for SrcIP: %v", err)
		return false
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var service model.Service
		if err := cursor.Decode(&service); err != nil {
			log.Printf("[Decode] Failed to decode service: %v", err)
			continue
		}
		logEntry.SrcType = "Service"
		logEntry.SrcNamespace = service.Namespace
		logEntry.SrcName = service.Name
		return true
	}
	return false
}

func fillMissingDestInfo(ctx context.Context, logEntry *model.APILog, servicesCol *mongo.Collection) bool {
	if logEntry.DstNamespace != "Unknown" {
		return true
	}
	cursor, err := servicesCol.Find(ctx, bson.M{"loadbalancerips": bson.M{"$in": []string{logEntry.DstIP}}})
	if err != nil {
		log.Printf("[MongoDB] Error checking services for DstIP: %v", err)
		return false
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var service model.Service
		if err := cursor.Decode(&service); err != nil {
			log.Printf("[Decode] Failed to decode service: %v", err)
			continue
		}
		logEntry.DstType = "Service"
		logEntry.DstNamespace = service.Namespace
		logEntry.DstName = service.Name
		return true
	}
	return false
}

// convertTimestamp - string timestamp를 RFC3339 형식으로 변환
func convertTimestamp(timestamp interface{}) string {
	switch v := timestamp.(type) {
	case string:
		if tsInt, err := strconv.ParseInt(v, 10, 64); err == nil {
			//return time.Unix(tsInt, 0).UTC().Format(time.RFC3339)
			return time.Unix(tsInt, 0).UTC().Format("2006-01-02 15:04:05")
		}
		log.Printf("[Warning] Invalid timestamp format: %v", v)
		return v
	default:
		log.Printf("[Warning] Unknown timestamp format: %v", v)
		return "unknown"
	}
}

// 3. 로그에서 cluster, namespace 필터링 후 응답
func FilterAPILogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Printf("[Request] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	// 요청 바디 파싱
	var requestBody struct {
		TimeRange  string `json:"timerange"`
		Namespaces []struct {
			Cluster   string `json:"cluster"`
			Namespace string `json:"namespace"`
		} `json:"namespaces"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("[Error] Invalid JSON: %v", err)
		return
	}

	if len(requestBody.Namespaces) == 0 {
		http.Error(w, "Empty request array", http.StatusBadRequest)
		return
	}

	// timerange 파싱 (기본: 5m)
	duration := 5 * time.Minute
	if requestBody.TimeRange != "" {
		parsed, err := time.ParseDuration(requestBody.TimeRange)
		if err != nil {
			http.Error(w, "Invalid timerange format", http.StatusBadRequest)
			log.Printf("[Error] Invalid timerange: %v", err)
			return
		}
		duration = parsed
	}

	client, err := config.GetMongoClient()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		log.Printf("[MongoDB] Connection error: %v", err)
		return
	}

	apiLogsCol := client.Database("SentryFlow").Collection("APILogs")
	servicesCol := client.Database("SentryFlow").Collection("Services")
	podsCol := client.Database("SentryFlow").Collection("Pods")

	// Get API Logs from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // timeout 5s
	defer cancel()

	// 현재 시간 기준 duration 전 타임스탬프 (Unix string)
	startTimestamp := time.Now().Add(-duration).Unix()
	filter := bson.M{
		"timestamp": bson.M{"$gte": strconv.FormatInt(startTimestamp, 10)},
	}

	// 내림차순 정렬 옵션 추가
	findOpts := options.Find()
	findOpts.SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := apiLogsCol.Find(ctx, filter, findOpts)
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

		// cluster name 채우기
		logEntry.SrcCluster = findCluster(ctx, logEntry.SrcType, logEntry.SrcName, logEntry.SrcNamespace, podsCol, servicesCol)
		logEntry.DstCluster = findCluster(ctx, logEntry.DstType, logEntry.DstName, logEntry.DstNamespace, podsCol, servicesCol)

		// Unknown 채우기
		if !fillMissingSourceInfo(ctx, &logEntry, servicesCol) || !fillMissingDestInfo(ctx, &logEntry, servicesCol) {
			continue
		}

		// timestamp 변환 (string -> RFC3339)
		logEntry.TimeStamp = convertTimestamp(logEntry.TimeStamp)

		// 특정 cluster and namespace 필터링
		for _, cond := range requestBody.Namespaces {
			if (logEntry.SrcCluster == cond.Cluster || logEntry.DstCluster == cond.Cluster) &&
				(logEntry.SrcNamespace == cond.Namespace || logEntry.DstNamespace == cond.Namespace) {
				logs = append(logs, logEntry)
				break // 하나라도 만족하면 추가하고 다음 로그로
			}
		}
	}

	log.Printf("[Response] Returning %d logs", len(logs))

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

	client, err := config.GetMongoClient()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		log.Printf("[MongoDB] Connection error: %v", err)
		return
	}

	collection := client.Database("SentryFlow").Collection("Pods")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch pods", http.StatusInternalServerError)
		log.Printf("[MongoDB] Query error: %v", err)
		return
	}
	defer cursor.Close(ctx)

	// 클러스터별 네임스페이스 수집
	clusterNamespaces := make(map[string]map[string]struct{})

	for cursor.Next(ctx) {
		var podEntry model.Pod
		if err := cursor.Decode(&podEntry); err != nil {
			http.Error(w, "Error decoding Pods", http.StatusInternalServerError)
			return
		}

		if _, ok := clusterNamespaces[podEntry.Cluster]; !ok {
			clusterNamespaces[podEntry.Cluster] = make(map[string]struct{})
		}
		clusterNamespaces[podEntry.Cluster][podEntry.Namespace] = struct{}{}
	}

	// Cluster 모델로 변환
	var clusters []model.Cluster
	for clusterName, nsMap := range clusterNamespaces {
		var namespaces []string
		for ns := range nsMap {
			namespaces = append(namespaces, ns)
		}
		clusters = append(clusters, model.Cluster{
			Name:       clusterName,
			Namespaces: namespaces,
		})
	}

	// JSON 응답
	if err := json.NewEncoder(w).Encode(clusters); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("[JSON] Encoding error: %v", err)
	}

	log.Printf("[Response] Returning %d logs", len(clusters))
}
