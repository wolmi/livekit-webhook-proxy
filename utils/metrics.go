package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromMetrics struct{}

func (m *PromMetrics) StartPrometheus(enable bool, port int) {

	if enable {

		//Disable golang metris
		prometheus.Unregister(collectors.NewGoCollector())

		prometheusMux := http.NewServeMux()
		// register endpoints
		prometheusMux.Handle("/metrics", promhttp.Handler()) //Prometheus endpoint
		prometheusMux.HandleFunc("/hltz", m.healthCheck)     //Health Check endpoint
		prometheusMux.HandleFunc("/", m.healthCheck)         //Health Check endpoint
		prometheusListenPort := strings.Join([]string{":", strconv.Itoa(port)}, "")
		prometheusServer := &http.Server{
			ReadHeaderTimeout: 5 * time.Second,
			Addr:              prometheusListenPort,
			Handler:           prometheusMux,
		}
		go func() {
			//Start server
			fmt.Println("Prometheus endpoint started in port", port)
			if err := prometheusServer.ListenAndServe(); err != nil {
				panic(err)
			}
		}()
	}
}

// Generic function to handle a helath check
func (m *PromMetrics) healthCheck(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}
