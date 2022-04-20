package main

import (
	"bytes"
	"encoding/json"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"net/http"
	"os"
)

var (
	V           = os.Getenv("API_VERSION")
	TOKEN       = os.Getenv("ACCESS_TOKEN")
	SERVER_ADDR = os.Getenv("SERVER_ADDR")
)

func main() {
	if V == "" || TOKEN == "" || SERVER_ADDR == "" {
		println("Please set env vars")
		println("API_VERSION - VK API version")
		println("ACCESS_TOKEN - bot token")
		println("SERVER_ADDR - server address")
		os.Exit(1)
	}

	vk := api.NewVK(TOKEN)
	vk.Version = V

	resp, err := vk.GroupsGetByID(nil)
	if err != nil {
		panic(err)
	}

	groupID := resp[0].ID

	lp, err := longpoll.NewLongPoll(vk, groupID)
	if err != nil {
		panic(err)
	}

	//lp.Goroutine(true)
	lp.FullResponse(func(r longpoll.Response) {
		RedirectToAnotherServer(SERVER_ADDR, r)
	})

	if err := lp.Run(); err != nil {
		panic(err)
	}

}

type Req struct {
	Type    events.EventType `json:"type"`
	Object  json.RawMessage  `json:"object"`
	GroupID int              `json:"group_id"`
	EventID string           `json:"event_id"`
}

// MarshalJSON
// struct this Req to json
func (r *Req) MarshalJSON() ([]byte, error) {
	type Alias Req
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}

func NewReq(e events.GroupEvent) *Req {
	return &Req{
		Type:    e.Type,
		Object:  e.Object,
		GroupID: 0,
		EventID: e.EventID,
	}

}

// RedirectToAnotherServer
// Request to another server with the same data
func RedirectToAnotherServer(addr string, r longpoll.Response) {
	for _, update := range r.Updates {
		update := update
		go func() {
			rawJson, _ := NewReq(update).MarshalJSON()

			resp, _ := http.Post(addr, "application/json", bytes.NewBuffer(rawJson))
			if resp == nil {
				println("Server did not return response")
			} else {
				println("Server response:", resp.Status)
			}
		}()
	}
}
