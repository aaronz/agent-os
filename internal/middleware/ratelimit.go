package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type RateLimiter struct {
	globalRequests map[string]*clientLimit
	agentRequests  map[string]map[string]*clientLimit
	orgRequests    map[string]map[string]*clientLimit
	mu             sync.RWMutex
	rate           int
	window         time.Duration
}

type clientLimit struct {
	count     int
	resetTime time.Time
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		globalRequests: make(map[string]*clientLimit),
		agentRequests:  make(map[string]map[string]*clientLimit),
		orgRequests:    make(map[string]map[string]*clientLimit),
		rate:           rate,
		window:         window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()

		for key, cl := range rl.globalRequests {
			if now.After(cl.resetTime) {
				delete(rl.globalRequests, key)
			}
		}

		for _, agentMap := range rl.agentRequests {
			for key, cl := range agentMap {
				if now.After(cl.resetTime) {
					delete(agentMap, key)
				}
			}
		}

		for _, orgMap := range rl.orgRequests {
			for key, cl := range orgMap {
				if now.After(cl.resetTime) {
					delete(orgMap, key)
				}
			}
		}

		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) AllowGlobal(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cl, exists := rl.globalRequests[key]

	if !exists || now.After(cl.resetTime) {
		rl.globalRequests[key] = &clientLimit{
			count:     1,
			resetTime: now.Add(rl.window),
		}
		return true
	}

	if cl.count >= rl.rate {
		return false
	}

	cl.count++
	return true
}

func (rl *RateLimiter) AllowAgent(orgID, agentID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if _, ok := rl.agentRequests[orgID]; !ok {
		rl.agentRequests[orgID] = make(map[string]*clientLimit)
	}

	now := time.Now()
	key := agentID
	cl, exists := rl.agentRequests[orgID][key]

	agentRate := rl.rate / 10
	if agentRate < 1 {
		agentRate = 1
	}

	if !exists || now.After(cl.resetTime) {
		rl.agentRequests[orgID][key] = &clientLimit{
			count:     1,
			resetTime: now.Add(rl.window),
		}
		return true
	}

	if cl.count >= agentRate {
		return false
	}

	cl.count++
	return true
}

func (rl *RateLimiter) AllowOrg(orgID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if _, ok := rl.orgRequests[orgID]; !ok {
		rl.orgRequests[orgID] = make(map[string]*clientLimit)
	}

	now := time.Now()
	key := orgID
	cl, exists := rl.orgRequests[orgID][key]

	orgRate := rl.rate
	if orgRate < 1 {
		orgRate = 1
	}

	if !exists || now.After(cl.resetTime) {
		rl.orgRequests[orgID][key] = &clientLimit{
			count:     1,
			resetTime: now.Add(rl.window),
		}
		return true
	}

	if cl.count >= orgRate {
		return false
	}

	cl.count++
	return true
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		orgID := r.Header.Get("X-Org-Id")

		key := apiKey
		if key == "" {
			key = r.RemoteAddr
		}

		if !rl.AllowGlobal(key) {
			http.Error(w, "rate limit exceeded (global)", http.StatusTooManyRequests)
			return
		}

		if apiKey != "" && orgID != "" {
			if !rl.AllowAgent(orgID, apiKey) {
				http.Error(w, "rate limit exceeded (agent)", http.StatusTooManyRequests)
				return
			}
		}

		if orgID != "" {
			if !rl.AllowOrg(orgID) {
				http.Error(w, "rate limit exceeded (organization)", http.StatusTooManyRequests)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func extractOrgID(r *http.Request) string {
	paths := strings.Split(r.URL.Path, "/")
	for i, p := range paths {
		if p == "organizations" && i+1 < len(paths) {
			return paths[i+1]
		}
	}
	return ""
}
