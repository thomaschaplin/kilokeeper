package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type WeightEntry struct {
	Date      string  `json:"date"`
	Kilograms float64 `json:"kilograms"`
}

var weights []WeightEntry

func loadJSON(filename string, target interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(target)
}

func saveJSON(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func getWeightsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weights)
}

func addWeightHandler(w http.ResponseWriter, r *http.Request) {
	var newWeight WeightEntry
	err := json.NewDecoder(r.Body).Decode(&newWeight)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("02/01/2006", newWeight.Date)
	if err != nil {
		http.Error(w, "Invalid date format, use DD/MM/YYYY", http.StatusBadRequest)
		return
	}

	weights = append(weights, newWeight)
	saveJSON("data/weights.json", weights)
	w.WriteHeader(http.StatusCreated)
}

func main() {
	loadJSON("data/weights.json", &weights)

	http.HandleFunc("/weights", getWeightsHandler)
	http.HandleFunc("/weights/add", addWeightHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
