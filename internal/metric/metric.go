package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const StateCreate = "created"
const StateUpdated = "updated"
const StateRemoved = "removed"

var (
	totalWatersNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "est_water_api_waters_not_found_total",
		Help: "Total number of waters that were not found",
	})

	totalWatersByState = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "est_water_api_waters_by_states_total",
		Help: "Total number of waters CUD operations",
	}, []string{"state"})
)

func IncTotalWaterNotFound() {
	totalWatersNotFound.Inc()
}

func IncTotalWaterState(state string) {
	totalWatersByState.WithLabelValues(state).Inc()
}
