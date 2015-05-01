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
	EventCollection string `json:"event_collection"`
	Timeframe       string `json:"timeframe,omitempty"` // Need to change this to be a pointer to another struct
	Interval        string `json:"interval,omitempty"`
	GroupBy         string `json:"group_by,omitempty"`
	// Filters
	// Steps
}

func (c *Client) query(path string, params AnalysisParams) (*http.Response, error) {
	// serialize payload
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// construct url
	url := baseUrl + c.ProjectID + "/queries" + path

	// new request
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
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

func (c *Client) storeInInterface(resp *http.Response, dest interface{}) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	return fmt.Errorf("Non 200 reply from keen.io: %s", data)
}

func (c *Client) Count(params AnalysisParams, result interface{}) {
	resp, err := c.query("/count", params)
	if err != nil {
		return err
	}

	return c.respToError(resp)
}

func (c *Client) CountUnique() {}
func (c *Client) Minimum()     {}
func (c *Client) Maximum()     {}
func (c *Client) Average()     {}
func (c *Client) Median()      {}
func (c *Client) Percentile()  {}
func (c *Client) Sum()         {}
