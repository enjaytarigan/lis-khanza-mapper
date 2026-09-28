package medqlab

import "testing"

// simrsNoOrderToMedQLabKey mirrors ApiMEDQLAB Java:
//   body.put("no_order", noorder.substring(4, 14));
//
// Example:
//   SIMRS full:  PK202602250001
//   MedQLab key: 2602250001
//   no_rawat:    2026/02/25/000001
func simrsNoOrderToMedQLabKey(noorder string) string {
	runes := []rune(noorder)
	if len(runes) <= 4 {
		return ""
	}
	end := 14
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[4:end])
}

func TestSimrsNoOrderToMedQLabKey(t *testing.T) {
	cases := []struct {
		full string
		want string
	}{
		{"PK202602250001", "2602250001"}, // production format (prefix PK + YYMMDD + seq)
		{"PL2509230001", "09230001"},     // shorter: from idx 4 to end
		{"PL250923000123", "0923000123"}, // 14+ chars → exactly 10 (idx 4..13)
		{"ABC", ""},
	}
	for _, c := range cases {
		got := simrsNoOrderToMedQLabKey(c.full)
		if got != c.want {
			t.Errorf("simrsNoOrderToMedQLabKey(%q)=%q want %q", c.full, got, c.want)
		}
	}
}

func TestExtractNoRawat(t *testing.T) {
	resp := &Response{
		NoPendaftaran: "2026/02/25/000001",
		Demographics: &Demographics{
			VisitNumber: "2026/02/25/000001",
			NoOrder:     "2602250001",
		},
		NoOrder: "2602250001",
	}
	if got := extractNoRawat(resp); got != "2026/02/25/000001" {
		t.Errorf("extractNoRawat = %q", got)
	}
	// Prefer visitNumber over response.no_pendaftaran
	resp.Demographics.VisitNumber = "2026/02/25/000002"
	if got := extractNoRawat(resp); got != "2026/02/25/000002" {
		t.Errorf("prefer visitNumber: got %q", got)
	}
}
