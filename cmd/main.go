package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type ConvertRequest struct {
	Value float64 `json:"value"`
	From  string  `json:"from"`
	To    string  `json:"to"`
}

type ConvertResponse struct {
	Result float64 `json:"result"`
}

func main() {
	// Route handler for homepage
	http.HandleFunc("/", home)
	http.HandleFunc("/api/convert", convertHandler)

	fmt.Println("Server is running on http://localhost:8080 🚀")

	// Start server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.html")
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	result, err := convert(req.Value, req.From, req.To)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ConvertResponse{Result: result})
}

func convert(value float64, from, to string) (float64, error) {
	lengthToM := map[string]float64{
		"m": 1, "km": 1000, "cm": 0.01, "mm": 0.001,
		"ft": 0.3046, "in": 0.0254, "yd": 0.9144, "mi": 1609.344,
	}
	weightToKg := map[string]float64{
		"kg": 1, "g": 0.001, "lb": 0.453592, "oz": 0.0283495, "t": 1000,
	}

	// length
	if fv, ok := lengthToM[from]; ok {
		tv, ok := lengthToM[to]
		if !ok {
			return 0, fmt.Errorf("unknown unit: %s", to)
		}
		return value * fv / tv, nil
	}

	// Weight
	if fv, ok := weightToKg[from]; ok {
		tv, ok := weightToKg[to]
		if !ok {
			return 0, fmt.Errorf("unknown unit: %s", to)
		}
		return value * fv / tv, nil
	}

	var celsius float64
	switch from {
	case "c":
		celsius = value
	case "f":
		celsius = (value - 32) * 5 / 9
	case "k":
		celsius = value - 273.15
	default:
		return 0, fmt.Errorf("unknown unit: %s", from)
	}

	switch to {
	case "c":
		return celsius, nil
	case "f":
		return celsius*9/5 + 32, nil
	case "k":
		return celsius + 273.15, nil
	default:
		return 0, fmt.Errorf("unknown unit: %s", to)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	// Parsing the specified template file being passes as input
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t.Execute(w, nil)
}
