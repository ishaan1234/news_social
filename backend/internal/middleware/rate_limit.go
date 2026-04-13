package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

func RateLimitMiddleware(rps int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		go cleanupClients()

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Unable to determine IP", http.StatusInternalServerError)
				return
			}

			mu.Lock()
			c, exists := clients[ip]
			if !exists {
				limiter := rate.NewLimiter(rate.Limit(rps), rps)
				clients[ip] = &client{
					limiter:  limiter,
					lastSeen: time.Now(),
				}
				c = clients[ip]
			}

			c.lastSeen = time.Now()
			mu.Unlock()

			if !c.limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func cleanupClients() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}