package grafana

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestDashboardsService_Search(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	baseURL, _ := url.Parse(server.URL + "/")
	client := NewClient(baseURL, "", nil)

	mux.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id": 1}, {"id": 2}]`)
	})

	opt := &DashboardSearchOptions{
		Query:     "q",
		Tags:      []string{"tag1", "tag2"},
		IsStarred: true,
		Limit:     10,
	}
	dashboards, err := client.Dashboards.Search(context.Background(), opt)
	if err != nil {
		t.Fatalf("Dashboards.Search returned error: %v", err)
	}

	want := []*DashboardHit{{ID: int64(1)}, {ID: int64(2)}}
	if !reflect.DeepEqual(dashboards, want) {
		t.Errorf("Dashboards.Search returned %+v, want %+v", dashboards, want)
	}
}
