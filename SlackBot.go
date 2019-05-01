package main

import (
	"fmt"
	"strconv"
	"strings"
)

var channels map[string]string

func prepareChannels(){
	channels = make(map[string]string)

	channels["general"] = "CBSH247U7"
	channels["direct_test"] = "DHXMS0VHV"
}

func main() {
	fmt.Println("### SLACK BOT INITIALIZING - BEEP BOOP ###")

	token := "xoxb-399682684981-606209338640-305lJnIbNinj26y8z2RKi7q0"

	respObj := SlackRTMStart(token)
	wsocket, _ := SlackWebsocketConnect(respObj)

	message, err := SlackRead(wsocket)

	if err != nil {
		fmt.Println("ERROR WITH SLACK READ OPERATION")
	}

	if message.Type == "hello"{
		fmt.Println("CONNECTION ESTABLISHED AND VERIFIED BY SLACK")
	}

	prepareChannels()

	message.Text = ":robot_face: KILL ALL HUMANS :robot_face:"
	message.Type = "message"
	message.Channel = channels["direct_test"]

	fmt.Println("WHAT IS MESSAGE: ")

	SlackWrite(wsocket, message)

	go SlackQuota(wsocket)

	for {
		message, err := SlackRead(wsocket)

		if err != nil {
			fmt.Println("ERROR WITH SLACK READ OPERATION")
		}

		if message.Channel != "" {
			fmt.Println("-- Message Details -- ")
			fmt.Println("ID: " + strconv.Itoa(int(message.Id)))
			fmt.Println("Channel: " + message.Channel)
			fmt.Println("Text: " + message.Text)
		}

		if strings.Contains(message.Text, "slacking_off"){
			if strings.Contains(message.Text, "quote"){
				message.Text = "'I would totally beat my employees. Not like hit them with my fists or anything, just smack them over the nose with a rolled up newspaper'"
				message.Type = "message"
				message.Channel = channels["general"]

				SlackWrite(wsocket, message)
			}
		}

		fmt.Println("Slack Message: " + message.Text)
	}
}
