package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/tidwall/buntdb"

	"github.com/patrickvalle/heatmap/cmd/apid/config"
	"github.com/patrickvalle/heatmap/internal/ipv6"
)

// IPV6 represents the IPv6 API method handler set.
type IPV6 struct {
	config config.Config
	db     *buntdb.DB
}

// List returns a list of IPv6 addresses.
func (i *IPV6) List(w http.ResponseWriter, r *http.Request, params map[string]string) {

	// Populate the filters based off the request.
	filters, err := parseFilters(r)
	if err != nil {
		log.Printf("ERROR: Failed while parsing filters: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the list of IPv6 addresses.
	results, err := ipv6.List(i.db, filters)
	if err != nil {
		log.Printf("ERROR: Failed while listing IPv6 addresses: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the response.
	out, err := proto.Marshal(&results)
	if err != nil {
		log.Printf("ERROR: Failed while marshalling response: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.WriteHeader(http.StatusOK)
	w.Write(out)

	return
}

// parseFilters generates a Filters object based off the supplied HTTP request.
func parseFilters(r *http.Request) (ipv6.Filters, error) {

	filters := ipv6.Filters{}

	minLongitude, err := stringToFloat(r.URL.Query().Get("minLongitude"))
	if err != nil {
		return filters, err
	}
	filters.MinLongitude = minLongitude

	maxLongitude, err := stringToFloat(r.URL.Query().Get("maxLongitude"))
	if err != nil {
		return filters, err
	}
	filters.MaxLongitude = maxLongitude

	minLatitude, err := stringToFloat(r.URL.Query().Get("minLatitude"))
	if err != nil {
		return filters, err
	}
	filters.MinLatitude = minLatitude

	maxLatitude, err := stringToFloat(r.URL.Query().Get("maxLatitude"))
	if err != nil {
		return filters, err
	}
	filters.MaxLatitude = maxLatitude

	return filters, nil
}

// populateFloatFromString takes a string value and attempts to populate it
// on the supplied *float32 target. This is a noop if the string value is empty.
func stringToFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}
