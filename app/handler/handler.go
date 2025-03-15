package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"sentryflow-api/app/model"
	"sentryflow-api/config"

	"go.mongodb.org/mongo-driver/bson"
)

// stored data example
var (
	users  = []model.User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
	lastID = 2
	//mu     sync.Mutex
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(users)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func APILogHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	client, err := config.GetMongoClient()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		log.Printf("[MongoDB] Connection error: %v", err)
		return
	}

	collection := client.Database("sentryflow").Collection("APILogs")

	switch r.Method {
	case http.MethodGet:
		// Get API Logs from MongoDB
		var logs []model.APILog
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			http.Error(w, "Failed to fetch API Logs", http.StatusInternalServerError)
			log.Printf("[MongoDB] Query error: %v", err)
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var logEntry model.APILog
			if err := cursor.Decode(&logEntry); err != nil {
				http.Error(w, "Error decoding API Logs", http.StatusInternalServerError)
				return
			}
			logs = append(logs, logEntry)
		}

		json.NewEncoder(w).Encode(logs)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
