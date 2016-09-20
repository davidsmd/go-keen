// same as go-keen
package keen

import (
	"bytes"
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"os"
	"time"
)

const (
	keenAPI = "https://api.keen.io/3.0/projects/"
)

type Connection struct {
	WriteKey   string
	ReadKey    string
	ProjectID  string
	HttpClient http.Client
}

type Client struct {
	*Connection
	CollectionName string
	flushTime      time.Duration
	timeout        time.Duration
	flushChan      chan int
	eventChan      chan *Event
	Addons         []interface{}
}

// this is the format that eventually ends up getting passed to request
type Event struct {
	Event interface{} `json:"event"`
	Keen  interface{} `json:"keen"` // keen object, includes timestamp overwrite
}

type Object struct {
	Timestamp string        `json:"timestamp"` // this needs to be ISO-8601, UTC
	Addons    []interface{} `json:"addons"`
}

type Addon struct {
	Name   string      `json:"name"`
	Input  interface{} `json:"input"`
	Output string      `json:"output"`
}

// NewConnection probably doesn't need to take arguments. Just check the environment for Write/Read/ProjectID.
func NewConnection() (connection *Connection, err error) {
	writekey := os.Getenv("KEEN_WRITE_KEY")
	projectid := os.Getenv("KEEN_PROJECT_ID")

	if writekey == "" {
		return nil, fmt.Errorf("Keen | Environment variable KEEN_WRITE_KEY must be set, found value: %v", writekey)
	}
	if projectid == "" {
		return nil, fmt.Errorf("Keen | Environment variable KEEN_PROJECT_ID must be set, found value: %v", projectid)
	}

	return &Connection{
		WriteKey:  writekey,
		ProjectID: projectid,
	}, nil
}

// NewClient creates a new Client.
func NewClient(collectionName string, flushTime time.Duration, timeout time.Duration) (client *Client, err error) {
	connection, err := NewConnection()
	if err != nil {
		return nil, err
	}

	if flushTime == time.Duration(0) {
		return nil, fmt.Errorf("Keen | flushTime duration must be longer than 0 seconds")
	}
	if timeout == time.Duration(0) {
		return nil, fmt.Errorf("Keen | timeout duration must be longer than 0 seconds")
	}

	client = &Client{
		Connection:     connection,
		CollectionName: collectionName,
		flushTime:      flushTime,
		timeout:        timeout,
		flushChan:      make(chan int),
		eventChan:      make(chan *Event),
		Addons:         make([]interface{}, 0),
	}

	return client, nil
}

// AttachAddon("keen:ua_parser", struct{ UAString string `json:"ua_string"` }{"agent"}, "agent_parsed")
// or something to that effect
func (c *Client) AttachAddon(name string, input interface{}, output string) {
	addon := Addon{
		Name:   name,
		Input:  input, // anonymous struct as passed in during method call
		Output: output,
	}
	c.Addons = append(c.Addons, addon)
}

// CreateEvent takes a userEvent of some type, creates a new baseObject, sets timestamp
// and addons for that baseObject, then creates a baseEvent to hold the baseObject and
// userEvent... the baseEvent is passed through the Client's eventChan and handled in the
// batch loop
func (c *Client) CreateEvent(userTimestamp string, userEvent interface{}) (err error) {
	baseObject := &Object{Timestamp: userTimestamp, Addons: c.Addons}
	baseEvent := &Event{Keen: baseObject, Event: userEvent}

	select {
	case c.eventChan <- baseEvent:
		// maybe a print here later for debug...
	case <-time.After(c.timeout):
		return fmt.Errorf("Keen | Timeout while adding event to batch.")
	}
	return nil
}

func (c *Client) BatchLoop() {
	// straight from the original version... mostly
	go func() {
		for _ = range time.Tick(c.flushTime) {
			c.flushChan <- 1
		}
	}()

	// removed some stuff relative to original
	batch := make(map[string][]interface{})
	for {
		select {
		case evt := <-c.eventChan:
			list, ok := batch[c.CollectionName]
			if !ok {
				list = make([]interface{}, 0)
			}
			batch[c.CollectionName] = append(list, evt)

		case <-c.flushChan:
			if len(batch) == 0 {
				continue
			}
			err := c.PushEvents(batch)

			if err != nil {
				log.Printf("Keen | Error during flush: %v", err)
			}

			batch = make(map[string][]interface{})
		}
	}
}

func (c *Client) PushEvents(events map[string][]interface{}) (err error) {
	resp, err := c.request(events)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("Returned status:", resp.StatusCode)

	return nil
}

// request in the original version took a payload as an interface{} that was
// generally assumed to be a map[string][]interface{}
func (c *Client) request(payload interface{}) (resp *http.Response, err error) {
	// serialize payload
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	log.Println("json formatted body:", string(body))

	// construct url
	url := keenAPI + c.ProjectID + "/events"

	// new request
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	// add auth
	req.Header.Add("Authorization", c.WriteKey)

	// set length/content-type
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
		req.ContentLength = int64(len(body))
	}

	resp, err = c.HttpClient.Do(req)

	return resp, err
}

// come back to this, I may need it to enforce some things w.r.t. how the library is used.
//type Interface interface {
//	Format(events map[string][]interface{}) []*Event
//}
