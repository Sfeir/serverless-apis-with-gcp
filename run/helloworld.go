package main

import (
        "fmt"
        "log"
        "net/http"
        "os"
)

func handler(w http.ResponseWriter, r *http.Request) {
        log.Print("Hello Cloud Night from Cloud Run received a request.")
        
	params, ok := r.URL.Query()["name"]
	name := "World";
	if ok && len(params[0]) > 1 {
        	log.Println("Url Param 'name' is found")
                name = params[0];
        }
        fmt.Fprintf(w, "<h1>Cloud Run</h1><h2>Hello %s!</h2>", name)
}

func main() {
        log.Print("Hello world sample started.")

        http.HandleFunc("/", handler)

        port := os.Getenv("PORT")
        if port == "" {
                port = "8080"
        }

        log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
