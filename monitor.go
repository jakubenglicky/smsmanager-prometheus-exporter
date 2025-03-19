package main

import "github.com/prometheus/client_golang/prometheus"

type Monitor struct {
	Registry   *prometheus.Registry
	CreditInfo *prometheus.GaugeVec
}

func NewMonitor() *Monitor {
	reg := prometheus.NewRegistry()
	monitor := &Monitor{
		Registry: reg,

		CreditInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "smsmanager_credit_info",
			Help: "Status of SMS Manager credit",
		}, []string{"account"}),
	}

	reg.MustRegister(monitor.CreditInfo)

	return monitor
}
