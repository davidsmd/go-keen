# Keen IO golang client SDK

## API Stability

**The master branch has no API stability guarantees.**


## Writing Events

This is the very beginnings of a Keen IO client SDK in Go. Currently, only adding events to collections is supported.


```go
package main


import (
    "github.com/davidsmd/go-keen"
    "time"
    "log"
)

const keenFlushInterval = 10 * time.Second
const keenTimeout = 10 * time.Second


type exampleEvent struct {
    UserId int  'json:"userid"'
    Amount int  'json:"amount"'
    Type string 'json:"type"'
    Tags []string  'json:"tags"'
    Timestamp string 'json:"Timestamp"'
}

func main() {
    exampleClient, err := keen.NewClient("exampleCollection", keenFlushInterval, keenTimeout)
    if err != nil {
        log.Println("Error creating client:", err)
    }

    // attach desired addons to client before starting batch collection loop
    exampleClient.AttachAddon("keen:ua_parser", struct{ UAString string `json:"ua_string"` }{"agent"}, "agent_parsed")

    go exampleClient.BatchLoop()

    testEvent := &exampleEvent{
        UserId:     11,
        Amount:     25,
        Type:       "Moose",
        Tags:       []string{"Antlers", "Awesome", "Aardvark"},
    }


    // populate timestamp from whatever source makes the most sense in your application
    exampleClient.CreateEvent(testEvent.Timestamp, testEvent)
}
```
