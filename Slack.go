package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"
)

type responseRtmStart struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	Needed string      `json:"needed"`
	Url   string       `json:"url"`
	Self  responseSelf `json:"self"`
}

type responseSelf struct {
	Id string `json:"id"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SlackRTMStart(token string)(respObj responseRtmStart) {
	fmt.Println("### STARTING RTM SESSION ###")

	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		errString := fmt.Sprintf("HTTP GET ERROR: %s", err)
		fmt.Println(errString)
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("ERROR READING RESPONSE BODY")
		return
	}

	json.Unmarshal(body, &respObj)

	if !respObj.Ok {
		err = fmt.Errorf("slack error: %s", respObj.Error)
		needed := fmt.Sprintf("Needed: %s", respObj.Needed)
		fmt.Println(err)
		fmt.Println(needed)
		return
	}

	fmt.Println("### KILL ALL HUMA.... errr... RTM SESSION ESTABLISHED ###")
	return
}

func SlackWebsocketConnect(response responseRtmStart)(wsocket *websocket.Conn, id string){
	wurl := response.Url
	id = response.Self.Id

	requestheader := http.Header{
		"Origin": {"https://api.slack.com/"} ,

		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"} ,
	}

	wsocket, _, err := websocket.DefaultDialer.Dial(wurl,requestheader)

	if err != nil{
		fmt.Println("WEBSOCKET CONNECT ERROR")
		fmt.Println(err)
	}

	return
}

type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

var counter uint64

func SlackRead(ws *websocket.Conn) (m Message, err error) {
	err = ws.ReadJSON(&m)
	return
}

func SlackWrite(wsocket *websocket.Conn, message Message){
	message.Id = atomic.AddUint64(&counter, 1)
	fmt.Println(message)
	wsocket.WriteJSON(message)
	return
}

func SlackQuota(wsocket *websocket.Conn){
	var message Message

	for{
		message.Id = 0
		message.Text = ":eggplant:"
		message.Type = "message"
		message.Channel = channels["general"]

		SlackWrite(wsocket, message)
		time.Sleep(3600000 * time.Millisecond)
	}
}
