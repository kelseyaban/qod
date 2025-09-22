package main

import (
	"fmt"
	"net/http"
	"net"
	"sync"
	"time"
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
	//Define a rat limiter struct
	type client struct {
		limiter *rate.Limiter
		lastSeen time.Time
	}

	var mu sync.Mutex // use to synchronize the map
  	var clients = make(map[string]*client)    // the actual map 
  	// A goroutine to remove stale entries from the map
  	go func() {
      	for {
          	time.Sleep(time.Minute)
          	mu.Lock() // begin cleanup
          	// delete any entry not seen in three minutes
          	for ip, client := range clients {
              	if time.Since(client.lastSeen) > 3 * time.Minute {
                  	delete(clients, ip)
              	}
          	}
        	mu.Unlock()    // finish clean up
        	}
 }()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if a.config.limiter.enabled {
	    // get the IP address
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			a.serverErrorResponse(w, r, err)
			return
		}
   
		mu.Lock()  // exclusive access to the map
		// check if ip address already in map, if not add it
		_, found := clients[ip]
	   if !found {
		   clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(a.config.limiter.rps), a.config.limiter.burst)}
	   }

	    // Update the last seem for the client
 		clients[ip].lastSeen = time.Now()

 		// Check the rate limit status
  		if !clients[ip].limiter.Allow() {
	  		mu.Unlock()        // no longer need exclusive access to the map
	  		a.rateLimitExceededResponse(w, r)
	  		return
  		}
 
  		mu.Unlock()      // others are free to get exclusive access to the map
	}
  		next.ServeHTTP(w, r)
})
 
 }
 