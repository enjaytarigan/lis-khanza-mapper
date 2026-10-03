package medqlab

import (
	"testing"
	"time"
)

func TestParseCollectDateIn(t *testing.T) {
	wib, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Fatal(err)
	}
	wita, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		t.Fatal(err)
	}

	tgl, jam := parseCollectDateIn("2025-08-07T06:52:48.391Z", wib)
	if tgl != "2025-08-07" || jam != "13:52:48" {
		t.Fatalf("WIB: got tgl=%q jam=%q want 2025-08-07 / 13:52:48", tgl, jam)
	}

	tgl2, jam2 := parseCollectDateIn("2025-08-07T06:52:48.391Z", wita)
	if tgl2 != "2025-08-07" || jam2 != "14:52:48" {
		t.Fatalf("WITA: got tgl=%q jam=%q want 2025-08-07 / 14:52:48", tgl2, jam2)
	}

	tgl3, jam3 := parseCollectDateIn("2025-08-07 06:52:48", wib)
	if tgl3 != "2025-08-07" || jam3 != "06:52:48" {
		t.Fatalf("naive datetime in WIB: got tgl=%q jam=%q", tgl3, jam3)
	}

	tgl4, jam4 := parseCollectDateIn("", wib)
	if tgl4 != "" || jam4 != "" {
		t.Fatalf("empty collectDate should yield empty, got tgl=%q jam=%q", tgl4, jam4)
	}
}
