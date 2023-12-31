package block

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var rollappHeightGauge = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "rollapp_height",
	Help: "The height of the local rollapp",
})
var rollappHubHeightGauge = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "rollapp_hub_height",
	Help: "The latest height of the Rollapp that has been synced to the hub.",
})
