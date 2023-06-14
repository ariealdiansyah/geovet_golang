package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"geovet-test/middleware"
	"geovet-test/models"

	"gorm.io/gorm"
)

type AnimalService struct {
	db *gorm.DB
}

func NewAnimalService(db *gorm.DB) *AnimalService {
	return &AnimalService{db: db}
}

func (s *AnimalService) AnimalHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.GetAnimals(w, r)
	case http.MethodPost:
		s.CreateAnimal(w, r)
	case http.MethodDelete:
		s.DeleteAnimal(w, r)
	case http.MethodPut:
		s.UpdateAnimal(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed")
	}
}

func (s *AnimalService) CreateAnimal(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the request context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid token or user not found")
		return
	}

	// Parse the request body
	var animal models.Animal
	err := json.NewDecoder(r.Body).Decode(&animal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error decoding request body")
		return
	}

	// Associate the animal with the user
	animal.ID = user.ID

	// Save the animal to the database
	result := s.db.Create(&animal)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error creating animal")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Animal created successfully")
}

func (s *AnimalService) GetAnimals(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the request context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid token or user not found")
		return
	}

	// Retrieve animals associated with the user
	var animals []models.Animal
	result := s.db.Where("user_id = ?", user.ID).Find(&animals)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error retrieving animals")
		return
	}

	// Convert animals to JSON
	jsonData, err := json.Marshal(animals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error converting animals to JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (s *AnimalService) UpdateAnimal(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the request context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid token or user not found")
		return
	}

	// Retrieve the animal ID from the request URL parameters
	animalID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid animal ID")
		return
	}

	// Retrieve the animal associated with the user
	var animal models.Animal
	result := s.db.Where("id = ? AND user_id = ?", animalID, user.ID).First(&animal)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Animal not found")
		return
	}

	// Parse the request body
	var updatedAnimal models.Animal
	err = json.NewDecoder(r.Body).Decode(&updatedAnimal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error decoding request body")
		return
	}

	// Update the animal
	animal.Name = updatedAnimal.Name
	animal.Species = updatedAnimal.Species
	animal.Species = updatedAnimal.Species
	animal.Species = updatedAnimal.Species

	result = s.db.Save(&animal)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error updating animal")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Animal updated successfully")
}

func (s *AnimalService) DeleteAnimal(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the request context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid token or user not found")
		return
	}

	// Retrieve the animal ID from the request URL parameters
	animalID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid animal ID")
		return
	}

	// Delete the animal associated with the user
	result := s.db.Where("id = ? AND user_id = ?", animalID, user.ID).Delete(&models.Animal{})
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error deleting animal")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Animal deleted successfully")
}
