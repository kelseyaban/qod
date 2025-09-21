package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

func (a *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer will be called when the stack unwinds
		defer func() {
			// recover() checks for panics
			err := recover()
			if err != nil {
				w.Header().Set("Connection", "close")
				a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Add("Vary", "Origin")
		// Let's check the request origin to see if it's in the trusted list
		origin := r.Header.Get("Origin")

		// Once we have a origin from the request header we need need to check
		if origin != "" {
			for i := range a.config.cors.trustedOrigins {
				if origin == a.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (a *application) rateLimit(next http.Handler) http.Handler {
	//create a rate limiter his rate limiter can initially handle 5 requests. But after that it can only handle 2 per second
	// (1 request every half-second).  If after handling 5 initial requests our server, a quarter-second later, receives another request, that
	// request will be blocked/dropped/queued since we need half-second to be able to process one new request. The bucket is initially full(5) but
	// emptied after 5 requests. We need time to refill it.
	limiter := rate.NewLimiter(2, 5)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this incoming request will be allowed. We will implement rateLimitExceededResponse(w, r) later
		// the Allow() method tries to remove a token from the bucket it returns false if the bucket is empty
		if !limiter.Allow() {
			a.rateLimitExceededResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})

}
