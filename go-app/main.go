package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	port    = ":8080"
	version = "v1.0"
)

var pingCounter = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "go_app_request_count",
		Help: "No of requests handled by the go-app",
	},
)

func root(w http.ResponseWriter, r *http.Request) {
	pingCounter.Inc()
	fmt.Fprintf(w, "Hello from API %s!\n", version)
	fmt.Println("Endpoint Hit: root")
}
func handleRequests() {
	http.HandleFunc("/", root)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(port, nil)
}

func main() {
	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
		if len(os.Args) > 2 {
			version = os.Args[2]
		}
	}
	fmt.Printf("API version: %s\n", version)
	fmt.Println("Listening on port", port)
	handleRequests()
}
