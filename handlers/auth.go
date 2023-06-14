package handlers

import (
	"fmt"
	"net/http"
	"time"

	"geovet-test/middleware"
	"geovet-test/models"
	"geovet-test/utils"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var signingKey = []byte("geovet") // Change this to your own secret key

type AuthService struct {
	db *gorm.DB
}

type SignupService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func NewSignupService(db *gorm.DB) *SignupService {
	return &SignupService{db: db}
}

func (s *SignupService) Signup(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	nameOwner := r.FormValue("name")
	addressOwner := r.FormValue("address")
	phoneOwner := r.FormValue("phone")

	// Check if the username already exists
	var existingUser models.User
	result := s.db.Where("username = ?", username).First(&existingUser)
	if result.Error == nil {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintln(w, "Username already exists")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error hashing password")
		return
	}

	// Create a new user
	newUser := models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	// Save the user to the database
	result = s.db.Create(&newUser)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error creating user")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User created successfully")

	var user models.User
	resultUser := s.db.Where("username = ?", username).First(&user)
	if resultUser.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Owner's Profile Not Found")
		return
	}

	newOwner := models.Owner{
		Name:    nameOwner,
		Address: addressOwner,
		Phone:   phoneOwner,
		UserID:  user.ID,
		User:    user,
	}

	result = s.db.Create(&newOwner)
	if result.Error != nil {
		// Rollback user creation if owner creation fails
		s.db.Delete(&user)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create owner")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Owner created successfully")

}

func (s *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Retrieve user from the database
	var user models.User
	result := s.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid credentials")
		return
	}

	// Compare the provided password with the stored hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid credentials")
		return
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      jwt.TimeFunc().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	// Generate encoded token and send it as response
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error generating token")
		return
	}

	var owner models.Owner
	resultOwner := s.db.Where("user_id = ?", user.ID).First(&owner)
	if resultOwner.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve owner's profile")
		return
	}

	response := map[string]interface{}{
		"user":  user,
		"owner": owner,
	}

	response["token"] = token

	utils.RespondWithJSON(w, http.StatusOK, response)
	fmt.Fprintln(w, tokenString)
}

func (s *AuthService) Home(w http.ResponseWriter, r *http.Request) {
	username := middleware.GetUserFromContext(r.Context())
	if username == nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid token or user not found")
		return
	}

	// Retrieve user from the database
	var user models.User
	result := s.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error retrieving user from database")
		return
	}

	// Authorized, display the username
	fmt.Fprintf(w, "Welcome, %s!", user.Username)
}
