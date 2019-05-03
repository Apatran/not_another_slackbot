package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var channels map[string]string

func prepareChannels(){
	channels = make(map[string]string)

	channels["general"] = "CBSH247U7"            //Place channel codes here
	channels["direct_test"] = "DHXMS0VHV"
}

func getToken(fileName string) (token string){
	file, err := os.Open(fileName) // just pass the file name that contains the token

	if err != nil {
		fmt.Print(err)
	}
	defer file.Close()

	scanner :=  bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	token = scanner.Text()

	file.Close()
	return
}

func main() {
	fmt.Println("### SLACK BOT INITIALIZING - BEEP BOOP ###")
	fmt.Println(os.Args)

	if len(os.Args) != 2 {
		fmt.Println("Please enter the file name which contains the auth token")
		os.Exit(-1)
	}

	token := getToken(os.Args[1])

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

	SlackWrite(wsocket, message)

	//go SlackQuota(wsocket)

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

		/*
		Catch and reply to your input here
		 */
	}
}
