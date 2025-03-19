package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ePort int

func init() {
	flag.IntVar(&ePort, "port", 8080, "Exporter port to listen on")
}

func main() {
	flag.Parse()

	monitor := NewMonitor()
	middleware := func(handlerFor http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accounts := LoadCreditInfo("accounts.json")

			// if err != nil {
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	_, _ = w.Write([]byte("could not scrape GitHub status"))
			// 	return
			// }

			for _, account := range accounts {
				val, err := strconv.ParseFloat(account.Credit, 64)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte("could not scrape GitHub status"))
					return
				}
				monitor.CreditInfo.WithLabelValues(account.Name).Set(val)
			}

			handlerFor.ServeHTTP(w, r)
		})
	}

	http.Handle("/metrics", middleware(promhttp.HandlerFor(monitor.Registry, promhttp.HandlerOpts{Registry: monitor.Registry})))
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`
		<html>
		<head>
		<title>SMS Manager Credit Info Exporter</title>
		</head>
		<body>
		<h1>SMS Manager Credit Info Exporter</h1>
		<p><a href="/metrics">Metrics</a></p>
		</body>
		</html>
		`))
	}))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", ePort), nil))
}
