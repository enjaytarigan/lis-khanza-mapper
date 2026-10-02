package medqlab

import "testing"

func TestParseCollectDate(t *testing.T) {
	tgl, jam := parseCollectDate("2025-08-07T06:52:48.391Z")
	if tgl == "" || jam == "" {
		t.Fatalf("expected parsed tanggal/jam, got tgl=%q jam=%q", tgl, jam)
	}
	// Local timezone may shift the calendar day; just ensure HH:MM:SS shape.
	if len(jam) != 8 {
		t.Fatalf("jam should be HH:MM:SS, got %q", jam)
	}
	if len(tgl) != 10 {
		t.Fatalf("tgl should be YYYY-MM-DD, got %q", tgl)
	}

	tgl2, jam2 := parseCollectDate("2025-08-07 06:52:48")
	if tgl2 != "2025-08-07" || jam2 != "06:52:48" {
		t.Fatalf("local datetime: got tgl=%q jam=%q", tgl2, jam2)
	}

	tgl3, jam3 := parseCollectDate("")
	if tgl3 != "" || jam3 != "" {
		t.Fatalf("empty collectDate should yield empty, got tgl=%q jam=%q", tgl3, jam3)
	}
}
