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

func TestDashboardsService_Get(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	baseURL, _ := url.Parse(server.URL + "/")
	client := NewClient(baseURL, "", nil)

	slug := "slug"
	mux.HandleFunc("/api/dashboards/db/"+slug, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"dashboard": {"id": 1}}`)
	})

	d, err := client.Dashboards.Get(context.Background(), slug)
	if err != nil {
		t.Fatalf("Dashboards.Get returned error: %v", err)
	}

	want := &Dashboard{ID: DashboardID(1)}
	if !reflect.DeepEqual(d, want) {
		t.Errorf("Dashboards.Get\nreturned: %+v\nwant: %+v", d, want)
	}

}

func TestDashboardsService_Create(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	baseURL, _ := url.Parse(server.URL + "/")
	client := NewClient(baseURL, "", nil)

	mux.HandleFunc("/api/dashboards/db", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id": 1}`)
	})

	dashboard := &Dashboard{}
	overwrite := false
	d, err := client.Dashboards.Create(context.Background(), dashboard, overwrite)
	if err != nil {
		t.Fatalf("Dashboards.Create returned error: %v", err)
	}

	want := &Dashboard{ID: DashboardID(1)}
	if !reflect.DeepEqual(d, want) {
		t.Errorf("Dashboards.Create\nreturned: %+v\nwant: %+v", d, want)
	}

}

func TestDashboardsService_Search(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	baseURL, _ := url.Parse(server.URL + "/")
	client := NewClient(baseURL, "", nil)

	mux.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		// TODO: check querystring params here
		testMethod(t, r, "GET")
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

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}
