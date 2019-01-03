package tests

import (
	"../db"
	"../indexer"
	"os"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {
	ex, err := os.Executable()
	if err != nil {
		t.Fatal("Not able to get the executable path")
	}

	dir := filepath.Dir(ex)
	filepath := dir + "/TestDbQuery.data"

	f, err := os.Create(filepath)
	if err != nil {
		t.Fatal("Not able to create test file")
	}
	defer f.Close()
	_, err = f.WriteString(fileContent)
	if err != nil {
		t.Fatal("Not able to create test file")
	}
	f.Sync()

	var d db.Data
	OK, err := indexer.ReadFile(filepath, &d)
	if !OK {
		if err != nil {
			t.Fatalf("indexer.ReadFile() failed: %s", err)
		} else {
			t.Fatal("indexer.ReadFile() failed to parse the test file")
		}

	}

	if err != nil {
		t.Fatalf("indexer.ReadFile() failed: %s", err)
	}

	c, OK := d.Cities4Search[5882142]
	if !OK || c.Name != "Acton Vale" {
		t.Fatal("indexer.ReadFile() failed for geoname '5882142' in test file")
	}

	c, OK = d.Cities4Search[5882873]
	if !OK || c.Name != "Ajax" {
		t.Fatal("indexer.ReadFile() failed for geoname '5882873' in test file")
	}

}
