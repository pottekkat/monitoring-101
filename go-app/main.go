package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// Prometheus packages
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Response stores the Message obtained from the python-app
type Response struct {
	Message string `json:"message"`
}

// default language and name
var (
	lang = "en"
	name = "John"
)

// create a new custom Prometheus counter metric
var pingCounter = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "go_app_request_count",
		Help: "No of requests handled by the go-app",
	},
)

// HelloHandler handles requests to the go-app
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	lang = r.URL.String()
	name = r.URL.Query()["name"][0]

	fmt.Println("Request for", lang, "with name", name)
	pingCounter.Inc()

	pUrl := os.Getenv("PYTHON_APP_URL")
	if len(pUrl) == 0 {
		pUrl = "localhost"
	}

	// call the python-app to obtain the translation
	resp, err := http.Get("http://" + pUrl + ":8000" + lang)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	resp.Body.Close()

	var m Response
	json.Unmarshal(body, &m)

	// send back response with "Hello name!" in the specified language
	fmt.Fprintf(w, "%s %s!", m.Message, name)
}

func main() {
	// expose Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", HelloHandler)

	http.ListenAndServe(":8080", nil)
}
