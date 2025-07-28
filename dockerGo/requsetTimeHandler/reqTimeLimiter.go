package requsetTimeHandler

import (
	"goProject/dockerGo/models"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// In-memory storage for clients
var clients = make(map[string]*models.Client)
var mu sync.Mutex

// Get a client's rate limiter or create one if it doesn't exist
func GetClientLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	// If the client already exists, return the existing limiter
	if client, exists := clients[ip]; exists {
		return client.Limiter
	}

	// Create a new limiter with 10 requests per minute
	limiter := rate.NewLimiter(rate.Every(time.Minute), 1)
	clients[ip] = &models.Client{Limiter: limiter}
	return limiter
}
