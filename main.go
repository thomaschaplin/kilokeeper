package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

type WeightEntry struct {
	Date      string  `json:"date"`
	Kilograms float64 `json:"kilograms"`
}

type WeightResponse struct {
	Date      string  `json:"date"`
	Kilograms float64 `json:"kilograms"`
	Stones    string  `json:"stones"`
	Age       int     `json:"age"`
	Bmi       string  `json:"bmi"`
	BmiStatus string  `json:"bmiStatus"`
	Goal      string  `json:"goal"`
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

func convertKilogramsToStonesPoundsOunces(weightInKilograms float64) string {
	// Convert kilograms to pounds
	weightInPounds := weightInKilograms * 2.20462

	// Calculate the number of stones
	stones := int(math.Floor(weightInPounds / 14))

	// Calculate the remaining pounds
	remainingPounds := int(math.Floor(math.Mod(weightInPounds, 14)))

	// Calculate the remaining ounces
	ounces := int(math.Round(math.Mod(weightInPounds, 1) * 16))

	return fmt.Sprintf("%d st %d lb %d oz", stones, remainingPounds, ounces)
}

func calculateAge(dob string, currentDate string) (int, error) {
	// Parse dates in DD/MM/YYYY format
	today, err := time.Parse("02/01/2006", currentDate)
	if err != nil {
		return 0, fmt.Errorf("invalid current date format: %v", err)
	}

	birthDate, err := time.Parse("02/01/2006", dob)
	if err != nil {
		return 0, fmt.Errorf("invalid date of birth format: %v", err)
	}

	age := today.Year() - birthDate.Year()
	if today.Month() < birthDate.Month() || (today.Month() == birthDate.Month() && today.Day() < birthDate.Day()) {
		age--
	}

	return age, nil
}

func calculateBMI(weightInKilograms, height float64) float64 {
	return weightInKilograms / math.Pow(height, 2)
}

func calculateBmiStatus(bmiValue float64) string {
	if bmiValue < 18.5 {
		return "Underweight"
	} else if bmiValue <= 24.9 {
		return "Normal"
	} else if bmiValue <= 29.9 {
		return "Overweight"
	} else if bmiValue <= 39.9 {
		return "Obese"
	} else if bmiValue >= 40 {
		return "Morbidly Obese"
	}
	return "Invalid data"
}

func getWeightsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dob := os.Getenv("DOB")
	if dob == "" {
		http.Error(w, "Date of birth (DOB) environment variable not set", http.StatusInternalServerError)
		return
	}

	height := os.Getenv("HEIGHT")
	if height == "" {
		http.Error(w, "Height (HEIGHT) environment variable not set", http.StatusInternalServerError)
		return
	}

	goal := os.Getenv("GOAL")
	if height == "" {
		http.Error(w, "Goal weight (GOAL) environment variable not set", http.StatusInternalServerError)
		return
	}

	var response []WeightResponse
	for _, entry := range weights {

		age, err := calculateAge(dob, entry.Date)
		if err != nil {
			http.Error(w, "Error calculating age", http.StatusInternalServerError)
			return
		}

		heightFloat, err := strconv.ParseFloat(height, 64)
		if err != nil {
			http.Error(w, "Invalid height format", http.StatusInternalServerError)
			return
		}

		response = append(response, WeightResponse{
			Date:      entry.Date,
			Kilograms: entry.Kilograms,
			Stones:    convertKilogramsToStonesPoundsOunces(entry.Kilograms),
			Age:       age,
			Bmi:       fmt.Sprintf("%.1f", calculateBMI(entry.Kilograms, heightFloat)),
			BmiStatus: calculateBmiStatus(calculateBMI(entry.Kilograms, heightFloat)),
			Goal:      goal,
		})
	}
	json.NewEncoder(w).Encode(response)
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	fileServer := http.FileServer(http.Dir("./web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
