package keen

import (
	"bytes"
	"encoding/json"
	"errors"
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
	MaxAge          int64  `json:"maxAge,omitempty"`
	TargetProperty  string `json:"target_property,omitempty"`
	// Filters
	// Steps
}

type AnalysisResult struct {
	Result int64 `json:"result,omitempty"`
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

func (c *Client) storeInInterface(resp *http.Response, dest *interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
	json.Unmarshal(data, &dest)

	if err != nil {
		return err
	}

	return fmt.Errorf("Non 200 reply from keen.io: %s", data)
}

func (c *Client) metric(metric string, params *AnalysisParams, dest interface{}) error {
	resp, err := c.query("/"+metric, params)
	if err != nil {
		return err
	}

	c.storeInInterface(resp, &dest)

	return nil
}

func (c *Client) Count(params *AnalysisParams, dest interface{}) error {
	return c.metric("count", params, dest)
}

func (c *Client) CountUnique(params *AnalysisParams, dest interface{}) error {
	if params.TargetProperty == "" {
		return errors.New("TargetProperty must be supplied")
	}
	return c.metric("count_unique", params, dest)
}
func (c *Client) Minimum(params *AnalysisParams, dest interface{}) error {
	if params.TargetProperty == "" {
		return errors.New("TargetProperty must be supplied")
	}
	return c.metric("minimum", params, dest)
}
func (c *Client) Maximum(params *AnalysisParams, dest interface{}) error {
	if params.TargetProperty == "" {
		return errors.New("TargetProperty must be supplied")
	}
	return c.metric("maximum", params, dest)
}
func (c *Client) Average(params *AnalysisParams, dest interface{}) error {
	if params.TargetProperty == "" {
		return errors.New("TargetProperty must be supplied")
	}
	return c.metric("average", params, dest)
}
func (c *Client) Median(params *AnalysisParams, dest interface{}) error {
	if params.TargetProperty == "" {
		return errors.New("TargetProperty must be supplied")
	}
	return c.metric("median", params, dest)
}
func (c *Client) Percentile(params *AnalysisParams, dest interface{}) error {
	if params.TargetProperty == "" {
		return errors.New("TargetProperty must be supplied")
	}
	return c.metric("percentile", params, dest)
}
func (c *Client) Sum(params *AnalysisParams, dest interface{}) error {
	if params.TargetProperty == "" {
		return errors.New("TargetProperty must be supplied")
	}
	return c.metric("sum", params, dest)
}
