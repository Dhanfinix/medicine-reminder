// Package models contains the data models for the Medicine Reminder API
package models

import (
	"time"
)

// Medicine represents a medication record in the system
type Medicine struct {
	ID        int       `json:"id" db:"id"`                   // Unique identifier for the medicine
	Name      string    `json:"name" db:"name"`               // Name of the medicine
	Dosage    string    `json:"dosage" db:"dosage"`           // Dosage amount (e.g., "500mg")
	Frequency string    `json:"frequency" db:"frequency"`     // How often to take (e.g., "3 times a day")
	TimeOfDay string    `json:"time_of_day" db:"time_of_day"` // JSON array of times to take the medicine
	StartDate time.Time `json:"start_date" db:"start_date"`   // When to start taking the medicine
	EndDate   time.Time `json:"end_date" db:"end_date"`       // When to stop taking the medicine
	Notes     string    `json:"notes" db:"notes"`             // Additional notes or instructions
	CreatedAt time.Time `json:"created_at" db:"created_at"`   // When the record was created
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`   // When the record was last updated
}

// MedicineInput represents the expected input format for creating/updating a medicine
type MedicineInput struct {
	Name      string    `json:"name"`
	Dosage    string    `json:"dosage"`
	Frequency string    `json:"frequency"`
	TimeOfDay []string  `json:"time_of_day"` // Array of times before conversion to JSON
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Notes     string    `json:"notes"`
}
