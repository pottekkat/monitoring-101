package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Response struct {
	Message string `json:"message"`
}

var (
	lang = "en"
	name = "John"
)

var pingCounter = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "go_app_request_count",
		Help: "No of requests handled by the go-app",
	},
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	lang = r.URL.String()
	name = r.URL.Query()["name"][0]

	fmt.Println("Request for", lang, "with name", name)
	pingCounter.Inc()

	pUrl := os.Getenv("PYTHON_APP_URL")
	if len(pUrl) == 0 {
		pUrl = "localhost"
	}

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

	fmt.Fprintf(w, "%s %s!", m.Message, name)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", HelloHandler)

	http.ListenAndServe(":8080", nil)
}
