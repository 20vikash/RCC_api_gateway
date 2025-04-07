package main

import (
	"sync"

	"golang.org/x/time/rate"
)

var mu sync.Mutex
var visitors = make(map[string]*rate.Limiter)

func (app *Application) getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(5, 5)
		visitors[ip] = limiter
	}
	return limiter
}
