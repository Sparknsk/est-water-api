package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	totalWaterEvents = promauto.NewCounter(prometheus.CounterOpts{
		Name: "retranslator_waters_events_total",
		Help: "Total number of waters that were not found",
	})

	totalWaterEventsNow = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "retranslator_waters_events_now_total",
		Help: "Total number of water events in retranslator",
	})
)

func AddTotalWaterEvents(value uint) {
	totalWaterEvents.Add(float64(value))
}

func AddTotalWaterEventsNow(value uint) {
	totalWaterEventsNow.Add(float64(value))
}

func SubTotalWaterEventsNow(value uint) {
	totalWaterEventsNow.Sub(float64(value))
}
