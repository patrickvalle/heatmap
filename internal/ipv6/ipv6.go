package ipv6

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
)

const collection = "address"
const expectedCSVColumns = 10
const coordinatePrecision = "%f"

var data map[string]int32 // lat,long:count

var buntdbValueRegex = regexp.MustCompile(`\[([\d-\.]+) ([\d-\.]+)\]`) // parses `lat` `long` info from `[lat long]`

// List fetches a list of IPv6 models according to the specified (optional) `Filters`.
func List(db *buntdb.DB, f Filters) (ListResult, error) {

	result := ListResult{}
	var maxCount int32

	err := db.View(func(tx *buntdb.Tx) error {

		// Fetch results from our buntdb index using the supplied filters.
		query := fmt.Sprintf("[%f %f],[%f %f]", f.MinLongitude, f.MinLatitude, f.MaxLongitude, f.MaxLatitude)
		err := tx.Intersects(collection, query, func(key, value string) bool {
			matches := buntdbValueRegex.FindStringSubmatch(value)
			longitude, _ := strconv.ParseFloat(matches[1], 64)
			latitude, _ := strconv.ParseFloat(matches[2], 64)
			count := data[fmt.Sprintf("%f,%f", latitude, longitude)]

			// Keep track of the max count for this result set.
			if count > maxCount {
				maxCount = count
			}

			result.Results = append(result.Results, &Address{
				Latitude:  latitude,
				Longitude: longitude,
				Count:     count,
			})
			return true
		})
		return err
	})
	if err != nil {
		return result, errors.Wrap(err, "failed running intersection query")
	}

	result.MaxCount = maxCount
	return result, err
}

// LoadData parses the data from the supplied source and populates our "persistence" layer.
func LoadData(db *buntdb.DB, source io.Reader) error {

	// Reset our "database".
	data = map[string]int32{}

	// Ensure we have a spatial index on our geo coords.
	db.CreateSpatialIndex(collection, fmt.Sprintf("%s:*:coords", collection), buntdb.IndexRect)

	// Setup a CSV reader on our data source.
	reader := csv.NewReader(source)
	if reader == nil {
		return errors.New("nil reader returned")
	}

	// Throwaway the first header line.
	_, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed reading from CSV reader")
	}

	// Populate our `data` with a map of `latitude,longitude : count`
	for {

		// Read the record.
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed reading from CSV reader")
		}
		if len(record) != expectedCSVColumns {
			return errors.Wrapf(fmt.Errorf("failed while reading record"), "expected %d columns in record, found %d", expectedCSVColumns, len(record))
		}

		// Grab the lat/long from the record.
		latitude, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			return errors.Wrap(err, "failed while parsing latitude float")
		}
		longitude, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			return errors.Wrap(err, "failed while parsing longitude float")
		}

		// Increment the count for this lat/long combo.
		key := fmt.Sprintf(coordinatePrecision+","+coordinatePrecision, latitude, longitude)
		data[key] = data[key] + 1
	}

	// Now that we have our `data` populated, index the coordinates in buntdb for fast lookups.
	err = db.Update(func(tx *buntdb.Tx) error {
		i := 0
		for coords := range data {
			key := fmt.Sprintf("%s:%d:coords", collection, i)
			latitude := strings.Split(coords, ",")[0]
			longitude := strings.Split(coords, ",")[1]
			value := fmt.Sprintf("[%s %s]", longitude, latitude)
			_, _, err := tx.Set(key, value, nil)
			if err != nil {
				return errors.Wrapf(err, "failed while indexing %s : %s", key, value)
			}
			i++
		}
		return nil
	})

	return nil
}
