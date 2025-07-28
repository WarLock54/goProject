package main

import (
	"encoding/json"
	"fmt"
	"goProject/dockerGo/initializers"
	"goProject/dockerGo/middleware"
	"goProject/dockerGo/models"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

var jwtKey = []byte("my_secret_key")

// jwt
func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	// Create the JWT claims, which includes the username and expiry time

	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// Create the Login Endpoint
func login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if creds.Username != "admin" || creds.Password != "password" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := generateToken(creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []models.Book
	if initializers.DB == nil {
		http.Error(w, "Veritabanı bağlantısı kurulamadı (DB is nil)", http.StatusInternalServerError)
		return
	}
	if err := initializers.DB.Find(&books).Error; err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Add a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if initializers.DB == nil {
		http.Error(w, "Veritabanı bağlantısı kurulamadı (DB is nil)", http.StatusInternalServerError)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}
	if err := initializers.DB.Create(&book).Error; err != nil {
		http.Error(w, "Kayıt eklenemedi", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// This middleware checks if the request has a valid JWT token. If not, it returns an unauthorized response.
func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenStr := c.Value
		claims := &models.Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

/*
// kısa yoldan postgre bağlantısı pgx  kütüphanesi ile yapılıyor

	func connectDB() *pgxpool.Pool {
		url := "postgres://postgres:admin@localhost:5432/dockerGo"
		config, err := pgxpool.ParseConfig(url)
		if err != nil {
			log.Fatalf("Unable to parse DB config: %v\n", err)
		}

		dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}
		return dbpool
	}
*/
func TestGetBooksHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getBooks)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
func main() {

	// Initialize the router
	r := mux.NewRouter()
	initializers.LoadEnvVariables()
	initializers.ConnectDB()

	r.Use(middleware.RateLimitingMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.ErrorHandlingMiddleware)

	r.HandleFunc("/login", login).Methods("POST")
	r.Handle("/books", authenticate(http.HandlerFunc(getBooks))).Methods("GET")
	r.Handle("/books", authenticate(http.HandlerFunc(createBook))).Methods("POST")

	fmt.Println("Server started on port :9001")
	log.Fatal(http.ListenAndServe(":9001", r))
}
