package keen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// AnalysisParams is a struct of all possible options for DataAnalysis queries. This
// struct will be marshalled into a POST body
type AnalysisParams struct {
	EventCollection string      `json:"event_collection"`
	Timeframe       interface{} `json:"timeframe,omitempty"`
	Timezone        string      `json:"timezone,omitempty"`
	Interval        string      `json:"interval,omitempty"`
	GroupBy         []string    `json:"group_by,omitempty"`
	MaxAge          int64       `json:"maxAge,omitempty"`
	TargetProperty  string      `json:"target_property,omitempty"`
	Filters         []Filter    `json:"filters,omitempty"`
	// Steps
}

type Analysis struct {
	Metric string
	Params AnalysisParams
}

// Filter is a struct that marshals into a JSON struct matching Keen analysis filters
// https://keen.io/docs/data-analysis/filters/
type Filter struct {
	PropertyName  string      `json:"property_name"`
	Operator      string      `json:"operator"`
	PropertyValue interface{} `json:"property_value"`
}

// Timeframe is a pair of strings denoting an absolute timeframe
// https://keen.io/docs/data-analysis/timeframe/#absolute-timeframes
type Timeframe struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func (c *Client) query(path string, params *AnalysisParams) (*http.Response, error) {
	// serialize payload
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// construct url
	url := baseUrl + c.ProjectID + "/queries" + path

	// new request
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// add auth
	req.Header.Add("Authorization", c.ReadKey)

	// set length/content-type
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
		req.ContentLength = int64(len(body))
	}

	return c.HttpClient.Do(req)
}

func getBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %s", err)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Non 200 reply from keen.io [%d]: %s", resp.StatusCode, data)
	}

	return data, nil
}

func (c *Client) MetricJSON(metric string, params *AnalysisParams) (string, error) {
	resp, _ := c.query("/"+metric, params)
	data, err := getBody(resp)
	return string(data), err
}

func (c *Client) Query(analysis Analysis) (string, error) {
	resp, _ := c.query("/"+analysis.Metric, &analysis.Params)
	data, err := getBody(resp)
	return string(data), err
}
