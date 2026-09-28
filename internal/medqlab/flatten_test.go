package medqlab

import (
	"encoding/json"
	"testing"
)

func TestFlattenAndSort(t *testing.T) {
	raw := []byte(`{
	  "examinations": [
	    {
	      "testId": 4,
	      "type": "DEPARTMENT",
	      "position": "2",
	      "testName": "IMMUNOLOGI",
	      "children": [
	        {
	          "testId": 537,
	          "type": "INDIVIDUAL",
	          "position": "2.7.11",
	          "testName": "Anti HCV Total",
	          "localCode": "HCVB",
	          "examValue": "Reactive",
	          "examValueFlag": "Reactive",
	          "normalValueText": "< 0.90",
	          "originalResult": "1.24",
	          "validatedAt": "2025-08-12T09:35:50.064Z",
	          "children": []
	        },
	        {
	          "testId": 248,
	          "type": "INDIVIDUAL",
	          "position": "2.2.20",
	          "testName": "Anti HIV",
	          "children": [
	            {
	              "testId": 249,
	              "type": "SUB_TEST",
	              "position": "2.2.20.1",
	              "testName": "Metode-1",
	              "localCode": "HIV-1",
	              "examValue": "Non-Reaktif",
	              "examValueFlag": "Non-Reaktif",
	              "normalValueText": "< 1.0",
	              "children": []
	            },
	            {
	              "testId": 250,
	              "type": "SUB_TEST",
	              "position": "2.2.20.10",
	              "testName": "Later",
	              "examValue": "x",
	              "children": []
	            }
	          ]
	        }
	      ]
	    }
	  ]
	}`)

	var wrap struct {
		Examinations []Examination `json:"examinations"`
	}
	if err := json.Unmarshal(raw, &wrap); err != nil {
		t.Fatal(err)
	}

	leaves := FlattenLeaves(wrap.Examinations)
	SortByPosition(leaves)

	if len(leaves) != 3 {
		t.Fatalf("want 3 leaves, got %d", len(leaves))
	}
	// Hierarchical sort: 2.2.20.1 < 2.2.20.10 < 2.7.11
	wantOrder := []string{"2.2.20.1", "2.2.20.10", "2.7.11"}
	for i, pos := range wantOrder {
		if leaves[i].Position != pos {
			t.Errorf("index %d: want position %s, got %s (%s)", i, pos, leaves[i].Position, leaves[i].TestName)
		}
	}
	if leaves[0].LisTestID != "249" || leaves[0].Nilai != "Non-Reaktif" {
		t.Errorf("first leaf unexpected: %+v", leaves[0])
	}
	if leaves[2].LisTestID != "537" || leaves[2].Keterangan != "Reactive" {
		t.Errorf("hcv leaf unexpected: %+v", leaves[2])
	}
}

func TestMapLeafKeteranganOnlyExamValueFlag(t *testing.T) {
	flag := "H"
	interp := "High"
	val := "10"
	leaf := mapLeaf(Examination{
		TestID:               json.Number("1"),
		Position:             "1",
		TestName:             "X",
		ExamValue:            &val,
		ExamValueFlag:        &flag,
		ResultInterpretation: &interp,
	})
	if leaf == nil || leaf.Keterangan != "H" {
		t.Fatalf("want flag only, got %+v", leaf)
	}

	leaf2 := mapLeaf(Examination{
		TestID:               json.Number("2"),
		Position:             "2",
		TestName:             "Y",
		ExamValue:            &val,
		ResultInterpretation: &interp,
	})
	if leaf2 == nil || leaf2.Keterangan != "" {
		t.Fatalf("want empty keterangan without flag, got %+v", leaf2)
	}
}

func TestFlattenSkipsHeaderWithoutValue(t *testing.T) {
	exams := []Examination{
		{
			TestID:   json.Number("1"),
			Position: "1",
			TestName: "Header only",
			Children: nil,
		},
	}
	leaves := FlattenLeaves(exams)
	if len(leaves) != 0 {
		t.Fatalf("expected empty, got %+v", leaves)
	}
}

func TestComparePosition(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1.2", "1.10", -1},
		{"1.10", "1.2", 1},
		{"2.2.20.1", "2.2.20.1", 0},
		{"1", "1.0", 0},
	}
	for _, c := range cases {
		got := comparePosition(c.a, c.b)
		if (got < 0 && c.want >= 0) || (got > 0 && c.want <= 0) || (got == 0 && c.want != 0) {
			t.Errorf("comparePosition(%q,%q)=%d want sign of %d", c.a, c.b, got, c.want)
		}
	}
}
