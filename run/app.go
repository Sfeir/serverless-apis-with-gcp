package main

import (
	"context"
	"fmt"
	"strconv"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type Microservice struct {
	Name        string
	Year        int64
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Cloud Run app inventory received a request.")

	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "serverless-codelab-sandbox")

	// [START datastore_ancestor_query]
	ancestor := datastore.NameKey("Environment", "Cloud Run", nil)
	query := datastore.NewQuery("Microservice").Ancestor(ancestor)	
	
	//var query *datastore.Query
	params, ok := r.URL.Query()["year"]
	if ok && len(params[0]) > 1 {
		log.Println("Url Param 'year' is found")
		year, _ := strconv.Atoi(params[0]) 
   		query = query.Filter("year =", year)
	} 
	// [END datastore_ancestor_query]

	it := client.Run(ctx, query)

	result := "<h1>Cloud Run Applications Inventory</h1><ul>";
	for {
		var microservice Microservice
		_, err := it.Next(&microservice)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching next microservice: %v", err)
		}
		result += "<li>" + microservice.Name + " (" + strconv.FormatInt(microservice.Year, 10) + ")</li>"
	}
	result += "</ul>"

	fmt.Fprintf(w, result)
}

func main() {
	log.Print("Cloud Run app inventory started.")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
