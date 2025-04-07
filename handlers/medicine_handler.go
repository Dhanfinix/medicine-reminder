// Package handlers contains HTTP handlers for the Medicine Reminder API
package handlers

import (
	"encoding/json"
	"fmt"
	"medicine-reminder/database"
	"medicine-reminder/models"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// GetMedicines handles GET /api/medicines
// Returns a list of all medicines
func GetMedicines(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM medicines ORDER BY created_at DESC")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	var medicines []models.Medicine
	for rows.Next() {
		var m models.Medicine
		err := rows.Scan(&m.ID, &m.Name, &m.Dosage, &m.Frequency, &m.TimeOfDay,
			&m.StartDate, &m.EndDate, &m.Notes, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error scanning database result")
			return
		}
		medicines = append(medicines, m)
	}

	respondWithJSON(w, http.StatusOK, medicines)
}

// GetMedicine handles GET /api/medicines/{id}
// Returns a specific medicine by ID
func GetMedicine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var medicine models.Medicine
	err := database.DB.QueryRow("SELECT * FROM medicines WHERE id = $1", id).Scan(
		&medicine.ID, &medicine.Name, &medicine.Dosage, &medicine.Frequency, &medicine.TimeOfDay,
		&medicine.StartDate, &medicine.EndDate, &medicine.Notes, &medicine.CreatedAt, &medicine.UpdatedAt,
	)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Medicine not found")
		return
	}

	respondWithJSON(w, http.StatusOK, medicine)
}

// CreateMedicine handles POST /api/medicines
// Creates a new medicine record
func CreateMedicine(w http.ResponseWriter, r *http.Request) {
	var input models.MedicineInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if err := validateMedicineInput(input); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Convert time_of_day array to JSON string
	timeOfDayJSON, err := json.Marshal(input.TimeOfDay)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error processing time of day")
		return
	}

	query := `
		INSERT INTO medicines (name, dosage, frequency, time_of_day, start_date, end_date, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING *`

	var medicine models.Medicine
	err = database.DB.QueryRow(
		query,
		input.Name,
		input.Dosage,
		input.Frequency,
		string(timeOfDayJSON),
		input.StartDate,
		input.EndDate,
		input.Notes,
		time.Now(),
		time.Now(),
	).Scan(
		&medicine.ID,
		&medicine.Name,
		&medicine.Dosage,
		&medicine.Frequency,
		&medicine.TimeOfDay,
		&medicine.StartDate,
		&medicine.EndDate,
		&medicine.Notes,
		&medicine.CreatedAt,
		&medicine.UpdatedAt,
	)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating medicine")
		return
	}

	respondWithJSON(w, http.StatusCreated, medicine)
}

// UpdateMedicine handles PUT /api/medicines/{id}
// Updates an existing medicine record
func UpdateMedicine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var input models.MedicineInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if err := validateMedicineInput(input); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Convert time_of_day array to JSON string
	timeOfDayJSON, err := json.Marshal(input.TimeOfDay)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error processing time of day")
		return
	}

	query := `
		UPDATE medicines 
		SET name = $1, dosage = $2, frequency = $3, time_of_day = $4, 
			start_date = $5, end_date = $6, notes = $7, updated_at = $8
		WHERE id = $9
		RETURNING *`

	var medicine models.Medicine
	err = database.DB.QueryRow(
		query,
		input.Name,
		input.Dosage,
		input.Frequency,
		string(timeOfDayJSON),
		input.StartDate,
		input.EndDate,
		input.Notes,
		time.Now(),
		id,
	).Scan(
		&medicine.ID,
		&medicine.Name,
		&medicine.Dosage,
		&medicine.Frequency,
		&medicine.TimeOfDay,
		&medicine.StartDate,
		&medicine.EndDate,
		&medicine.Notes,
		&medicine.CreatedAt,
		&medicine.UpdatedAt,
	)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Medicine not found")
		return
	}

	respondWithJSON(w, http.StatusOK, medicine)
}

// DeleteMedicine handles DELETE /api/medicines/{id}
// Deletes a medicine record
func DeleteMedicine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := database.DB.Exec("DELETE FROM medicines WHERE id = $1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting medicine")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error checking delete result")
		return
	}

	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Medicine not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper functions

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func validateMedicineInput(input models.MedicineInput) error {
	if input.Name == "" {
		return fmt.Errorf("name is required")
	}
	if input.Dosage == "" {
		return fmt.Errorf("dosage is required")
	}
	if input.Frequency == "" {
		return fmt.Errorf("frequency is required")
	}
	if len(input.TimeOfDay) == 0 {
		return fmt.Errorf("time of day is required")
	}
	if input.StartDate.IsZero() {
		return fmt.Errorf("start date is required")
	}
	if input.EndDate.IsZero() {
		return fmt.Errorf("end date is required")
	}
	if input.EndDate.Before(input.StartDate) {
		return fmt.Errorf("end date must be after start date")
	}
	return nil
}
