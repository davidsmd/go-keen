package keen

import (
	"log"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	flushInterval := time.Duration(10) * time.Second
	timeout := time.Duration(15) * time.Second
	c, err := NewClient("testCollection", flushInterval, timeout)

	if err != nil {
		log.Println("Error creating new Client:", err)
	}

	log.Printf("Created new Client with specifications: %v", c)

}

func TestNewEvent(t *testing.T) {
	flushInterval := time.Duration(5) * time.Second
	timeout := time.Duration(5) * time.Second
	c, err := NewClient("testCollection", flushInterval, timeout)

	if err != nil {
		log.Println("Error creating new Client:", err)
	}

	log.Printf("Created new Client with specifications: %v", c)

	type customEvent struct {
		Timestamp  string   `json:"timestamp"`
		DeviceType int      `json:"devicetype"`
		AppVersion int      `json:"appversion"`
		EventName  string   `json:"eventname"`
		Metadata   []string `json:"metadata"`
	}

	log.Println("starting batch loop")
	go c.BatchLoop()

	testEvent := &customEvent{
		Timestamp:  "2016-09-20T16:42:08.896Z",
		DeviceType: 1,
		AppVersion: 5,
		EventName:  "the Frank Event",
		Metadata:   []string{"Weasel", "Goat", "Badger"},
	}

	log.Printf("Custom event: %v", testEvent)

	testEvent2 := &customEvent{
		Timestamp:  "2016-09-20T16:42:09.896Z",
		DeviceType: 1,
		AppVersion: 5,
		EventName:  "the Dave Event",
		Metadata:   []string{"Capybara", "Rat", "Moose"},
	}

	log.Printf("Custom event: %v", testEvent2)

	err = c.CreateEvent(testEvent.Timestamp, testEvent)
	if err != nil {
		log.Printf("Error creating event: %v", err)
	}

	err = c.CreateEvent(testEvent2.Timestamp, testEvent2)
	if err != nil {
		log.Printf("Error creating event: %v", err)
	}

	log.Printf("Created event, contents of c.eventChan: %v", c.eventChan)

	time.Sleep(20 * time.Second)

}
