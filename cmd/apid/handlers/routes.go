package handlers

import (
	"net/http"

	"github.com/dimfeld/httptreemux"

	"github.com/patrickvalle/heatmap/cmd/apid/config"
)

// API returns a handler for a set of routes.
func API(c config.Config) http.Handler {

	router := httptreemux.New()

	ipV6 := IPV6{c}
	router.GET("/v1/ipv6", ipV6.List)

	return router
}
