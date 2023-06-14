package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"geovet-test/models"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

// Define a context key for the user
const userContextKey = "user"

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

type JWTMiddleware struct {
	db *gorm.DB
}

func NewJWTMiddleware(db *gorm.DB) *JWTMiddleware {
	return &JWTMiddleware{db: db}
}

func (m *JWTMiddleware) ValidateJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the JWT token from the request header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Missing authorization header")
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("geovet"), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Invalid token")
			return
		}

		// Check if the token is valid and not expired
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// Retrieve the user from the database
			var user models.User
			result := m.db.First(&user, claims.UserID)
			if result.Error != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Error retrieving user")
				return
			}

			// Add the user to the request context
			ctx := context.WithValue(r.Context(), userContextKey, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Invalid token")
			return
		}
	}
}

func GetUserFromContext(ctx context.Context) *models.User {
	if user, ok := ctx.Value(userContextKey).(*models.User); ok {
		return user
	}
	return nil
}

func GenerateToken(user *models.User, secretKey string) (string, error) {
	// Create the claims for the JWT token
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Set token expiration to 24 hours
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Generate the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
