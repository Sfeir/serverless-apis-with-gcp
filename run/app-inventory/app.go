package app_inventory

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Cloud Run app inventory received a request.")

	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")

	// [START datastore_ancestor_query]
	ancestor := datastore.NameKey("Environment", "Cloud Run", nil)

	params, ok := r.URL.Query()["year"]
	if ok && len(params[0]) > 1 {
		log.Println("Url Param 'year' is found")
		query := datastore.NewQuery("Microservice").Ancestor(ancestor).Filter("year =", params[0])
	} else {
		query := datastore.NewQuery("Microservice").Ancestor(ancestor)
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
		result += "<li>" + microservice.name + "(" + microservice.year + ")</li>"
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
