package main

import (
	"flag"
	"log"
	"net/http"
)

// create a simple HTML page with some JS added. Obviously in a professional
// setting, we would have the JS code in a script file
const html = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
</head>
<body>
    <h1>CORS Implementation</h1>
    <div id="output"></div>
    <script>
    document.addEventListener('DOMContentLoaded', function() {
        fetch("http://localhost:4000/v1/healthcheck")
        .then(response => response.json())
        .then(data => {
            document.getElementById("output").innerHTML =
                "<pre>" + JSON.stringify(data, null, 2) + "</pre>";
        })
        .catch(err => {
            document.getElementById("output").textContent = err;
        });
    });
</script>
  </body>
  </html>`
  
  // A very simple HTTP server
  func main() {
	  addr := flag.String("addr", ":9000", "Server address")
	  flag.Parse()
  
	  log.Printf("starting server on %s", *addr)
	  err := http.ListenAndServe(*addr, 
			 http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(html))
		 }))
	  log.Fatal(err)
  }
  
