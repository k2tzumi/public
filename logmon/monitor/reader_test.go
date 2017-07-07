package monitor

import (
	"container/ring"
	"net/url"
	"testing"
)

func TestExtractSection(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"http://www.example.org/section-a/subsection-b/", "/section-a"},
		{"http://www.example.org/section-a/subsection-b", "/section-a"},
		{"http://www.example.org/section-a/", "/section-a"},
		{"http://www.example.org/section-a", "/section-a"},
		{"http://www.example.org/", "/"},
	}
	for _, tt := range tests {
		if got := extractSection(mustParseURL(tt.url)); got != tt.want {
			t.Errorf("extractSection() = %v, want %v", got, tt.want)
		}
	}
}

func mustParseURL(u string) *url.URL {
	p, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return p
}

func TestLogReaderAlerts(t *testing.T) {
	type faketraffic struct {
		hits, persecond uint64
	}
	run := func(timespan int, threshold uint64, traffic []faketraffic) []Alert {
		lr := &Reader{
			alerttimespan: ring.New(timespan), // $timespan seconds
			alertthresh:   threshold,          // $threshold hits
		}
		var alarmed bool
		for _, trf := range traffic {
			for i := uint64(0); i < trf.hits; i += trf.persecond {
				lr.totalhits += trf.persecond
				alarmed = lr.checkalert(alarmed)
			}
		}
		return lr.alerts
	}

	tests := []struct {
		name           string
		timespan       int
		threshold      uint64
		traffic        []faketraffic
		expectedalerts int
		expectedopen   int
	}{
		{"normal", 5, 1000, []faketraffic{{10000, 500}}, 0, 0},
		{"peak", 5, 1000, []faketraffic{{5000, 1000}}, 1, 1},
		{"peak-normal", 5, 1000, []faketraffic{{5000, 1000}, {10000, 500}}, 1, 0},
		{"peak-normal-peak", 5, 1000, []faketraffic{{5000, 1000}, {10000, 500}, {5000, 1000}}, 2, 1},
		{"peak-normal-peak-normal", 5, 1000, []faketraffic{{5000, 1000}, {10000, 500}, {5000, 1000}, {10000, 500}}, 2, 0},
	}
	for _, tt := range tests {
		alerts := run(tt.timespan, tt.threshold, tt.traffic)

		if l := len(alerts); tt.expectedalerts != l {
			t.Errorf("%q unexpected number of alerts found. expected: %d. got: %d", tt.name, tt.expectedalerts, l)
			t.Errorf("Alerts:\n%#v", alerts)
		}

		var openAlerts int
		for _, alert := range alerts {
			if alert.End.IsZero() {
				openAlerts++
			}
		}
		if tt.expectedopen != openAlerts {
			t.Errorf("%q unexpected number of open alerts found. expected: %d. got: %d", tt.name, tt.expectedopen, openAlerts)
		}
	}
}
