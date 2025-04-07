package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"medicine-reminder/database"
	"medicine-reminder/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) {
	// Initialize test database
	database.InitDB()

	// Clear the medicines table
	_, err := database.DB.Exec("DELETE FROM medicines")
	assert.NoError(t, err)
}

func TestGetMedicines(t *testing.T) {
	setupTestDB(t)

	// Create a test medicine
	medicine := createTestMedicine(t)

	// Create request
	req, err := http.NewRequest("GET", "/api/medicines", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(GetMedicines)
	handler.ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var medicines []models.Medicine
	err = json.Unmarshal(rr.Body.Bytes(), &medicines)
	assert.NoError(t, err)

	// Verify response
	assert.Equal(t, 1, len(medicines))
	assert.Equal(t, medicine.Name, medicines[0].Name)
	assert.Equal(t, medicine.Dosage, medicines[0].Dosage)
}

func TestGetMedicine(t *testing.T) {
	setupTestDB(t)

	// Create a test medicine
	medicine := createTestMedicine(t)

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/medicines/%d", medicine.ID), nil)
	assert.NoError(t, err)

	// Add URL parameters to request
	vars := map[string]string{
		"id": fmt.Sprintf("%d", medicine.ID),
	}
	req = mux.SetURLVars(req, vars)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(GetMedicine)
	handler.ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Medicine
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response
	assert.Equal(t, medicine.Name, response.Name)
	assert.Equal(t, medicine.Dosage, response.Dosage)
}

func TestCreateMedicine(t *testing.T) {
	setupTestDB(t)

	// Create test input
	input := models.MedicineInput{
		Name:      "Test Medicine",
		Dosage:    "100mg",
		Frequency: "Once daily",
		TimeOfDay: []string{"09:00"},
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 7),
		Notes:     "Test notes",
	}

	// Convert input to JSON
	body, err := json.Marshal(input)
	assert.NoError(t, err)

	// Create request
	req, err := http.NewRequest("POST", "/api/medicines", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(CreateMedicine)
	handler.ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Parse response
	var response models.Medicine
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response
	assert.Equal(t, input.Name, response.Name)
	assert.Equal(t, input.Dosage, response.Dosage)
	assert.NotZero(t, response.ID)
}

func TestUpdateMedicine(t *testing.T) {
	setupTestDB(t)

	// Create a test medicine
	medicine := createTestMedicine(t)

	// Create update input
	input := models.MedicineInput{
		Name:      "Updated Medicine",
		Dosage:    "200mg",
		Frequency: "Twice daily",
		TimeOfDay: []string{"09:00", "21:00"},
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 14),
		Notes:     "Updated notes",
	}

	// Convert input to JSON
	body, err := json.Marshal(input)
	assert.NoError(t, err)

	// Create request
	req, err := http.NewRequest("PUT", fmt.Sprintf("/api/medicines/%d", medicine.ID), bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Add URL parameters to request
	vars := map[string]string{
		"id": fmt.Sprintf("%d", medicine.ID),
	}
	req = mux.SetURLVars(req, vars)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(UpdateMedicine)
	handler.ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Medicine
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response
	assert.Equal(t, input.Name, response.Name)
	assert.Equal(t, input.Dosage, response.Dosage)
	assert.Equal(t, medicine.ID, response.ID)
}

func TestDeleteMedicine(t *testing.T) {
	setupTestDB(t)

	// Create a test medicine
	medicine := createTestMedicine(t)

	// Create request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/medicines/%d", medicine.ID), nil)
	assert.NoError(t, err)

	// Add URL parameters to request
	vars := map[string]string{
		"id": fmt.Sprintf("%d", medicine.ID),
	}
	req = mux.SetURLVars(req, vars)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(DeleteMedicine)
	handler.ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Verify medicine is deleted
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM medicines WHERE id = $1", medicine.ID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestCreateMultipleMedicines(t *testing.T) {
	setupTestDB(t)

	// Create 10 medicines
	medicines := createMultipleTestMedicines(t, 10)

	// Verify all medicines were created
	assert.Equal(t, 10, len(medicines))

	// Verify each medicine has valid data
	for _, medicine := range medicines {
		assert.NotZero(t, medicine.ID)
		assert.NotEmpty(t, medicine.Name)
		assert.NotEmpty(t, medicine.Dosage)
		assert.NotEmpty(t, medicine.Frequency)
		assert.NotEmpty(t, medicine.TimeOfDay)
		assert.NotZero(t, medicine.StartDate)
		assert.NotZero(t, medicine.EndDate)
		assert.NotEmpty(t, medicine.Notes)
		assert.NotZero(t, medicine.CreatedAt)
		assert.NotZero(t, medicine.UpdatedAt)

		// Verify the medicine exists in the database
		var count int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM medicines WHERE id = $1", medicine.ID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	}
}

// Helper function to create a test medicine
func createTestMedicine(t *testing.T) models.Medicine {
	input := models.MedicineInput{
		Name:      "Test Medicine",
		Dosage:    "100mg",
		Frequency: "Once daily",
		TimeOfDay: []string{"09:00"},
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 7),
		Notes:     "Test notes",
	}

	timeOfDayJSON, err := json.Marshal(input.TimeOfDay)
	assert.NoError(t, err)

	var medicine models.Medicine
	err = database.DB.QueryRow(`
		INSERT INTO medicines (name, dosage, frequency, time_of_day, start_date, end_date, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING *`,
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
	assert.NoError(t, err)

	return medicine
}

// Helper function to create multiple test medicines
func createMultipleTestMedicines(t *testing.T, count int) []models.Medicine {
	medicines := make([]models.Medicine, count)

	medicineNames := []string{
		"Aspirin",
		"Paracetamol",
		"Ibuprofen",
		"Amoxicillin",
		"Metformin",
		"Lisinopril",
		"Amlodipine",
		"Atorvastatin",
		"Omeprazole",
		"Metoprolol",
	}

	medicineDosages := []string{
		"100mg",
		"500mg",
		"200mg",
		"250mg",
		"500mg",
		"10mg",
		"5mg",
		"20mg",
		"20mg",
		"25mg",
	}

	frequencies := []string{
		"Once daily",
		"Every 6 hours",
		"Every 8 hours",
		"Twice daily",
		"Twice daily",
		"Once daily",
		"Once daily",
		"Once daily",
		"Once daily",
		"Twice daily",
	}

	timeOfDayOptions := [][]string{
		{"08:00"},
		{"08:00", "14:00", "20:00", "02:00"},
		{"08:00", "16:00", "00:00"},
		{"09:00", "21:00"},
		{"08:00", "20:00"},
		{"08:00"},
		{"08:00"},
		{"20:00"},
		{"08:00"},
		{"08:00", "20:00"},
	}

	for i := 0; i < count; i++ {
		input := models.MedicineInput{
			Name:      medicineNames[i],
			Dosage:    medicineDosages[i],
			Frequency: frequencies[i],
			TimeOfDay: timeOfDayOptions[i],
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 30),
			Notes:     fmt.Sprintf("Notes for %s", medicineNames[i]),
		}

		timeOfDayJSON, err := json.Marshal(input.TimeOfDay)
		assert.NoError(t, err)

		err = database.DB.QueryRow(`
			INSERT INTO medicines (name, dosage, frequency, time_of_day, start_date, end_date, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING *`,
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
			&medicines[i].ID,
			&medicines[i].Name,
			&medicines[i].Dosage,
			&medicines[i].Frequency,
			&medicines[i].TimeOfDay,
			&medicines[i].StartDate,
			&medicines[i].EndDate,
			&medicines[i].Notes,
			&medicines[i].CreatedAt,
			&medicines[i].UpdatedAt,
		)
		assert.NoError(t, err)
	}

	return medicines
}
