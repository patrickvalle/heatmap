package ipv6_test

import (
	"os"
	"testing"

	"github.com/tidwall/buntdb"

	"github.com/patrickvalle/heatmap/internal/ipv6"
)

func TestLoadData(t *testing.T) {

	// Setup.
	file, err := os.Open("../../GeoLite2-City-Blocks-IPv6.csv")
	if err != nil {
		t.Fatalf("Failed to load input: %s", err.Error())
	}
	db, err := buntdb.Open(":memory:")
	if err != nil {
		t.Fatalf("Failed to open db: %s", err.Error())
	}

	// Call.
	err = ipv6.LoadData(db, file)

	// Verify.
	if err != nil {
		t.Fatalf("LoadData should not produce an error, got: %s", err.Error())
	}
}
