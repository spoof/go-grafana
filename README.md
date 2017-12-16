# go-grafana [![GoDoc](https://godoc.org/github.com/spoof/go-grafana?status.svg)](https://godoc.org/github.com/spoof/go-grafana) [![Go Report Card](https://goreportcard.com/badge/github.com/spoof/go-grafana)](https://goreportcard.com/report/github.com/spoof/go-grafana) [![Build Status](https://travis-ci.org/spoof/go-grafana.svg?branch=master)](https://travis-ci.org/spoof/go-grafana)
go-grafana is a Go client library for [Grafana API](http://docs.grafana.org/http_api/)

## Installation
```
go get -u github.com/spoof/go-grafana
```

## Usage
```go
package main

import (
	"context"
	"log"
	"net/url"

	"github.com/spoof/go-grafana/client"
	"github.com/spoof/go-grafana/grafana"
	"github.com/spoof/go-grafana/grafana/panel"
	"github.com/spoof/go-grafana/grafana/query"
)

func main() {
	const token = "<token>"
	url, _ := url.Parse("http://localhost:3000/")
	client := client.NewClient(url, token, nil)
	ctx := context.Background()

	// Create dashboard
	d := grafana.NewDashboard("Title Demo")
	d.Tags.Set("tag1", "tag2")

	// Add row
	row := grafana.NewRow()
	row.Title = "Row"
	row.ShowTitle = true
 
  // Add Graph panel
	p := panel.NewGraph()
	p.GeneralOptions().Height = "250px"
	p.GeneralOptions().Span = 2
	p.XAxis.Show = true
	p.YAxes.Left.Format = "short"
	p.YAxes.Left.Show = true

  // Graph panel with query to Prometheus datasource
	q := query.NewPrometheus("Prometheus")
	q.Expression = `up{job="job"}`
	q.Interval = "1"
	q.LegendFormat = "{{ instance }}"

	queries := p.Queries()
	*queries = []panel.Query{q}

	row.Panels = append(row.Panels, p)
	d.Rows = append(d.Rows, row)

	overwrite := true
	if err := client.Dashboards.Save(ctx, d, overwrite); err != nil {
		log.Fatalf("Error while saving dashboard %s", err)
	}

}
```

## Current Status

Project is under active development. It's not ready for production use due to high risk of API changes.
Please checkout out [TODO](./TODO.md) for more information what's done and what's going to be done.


## License ##

This library is distributed under the Apache 2.0 license found in the [LICENSE](./LICENSE)
file.
