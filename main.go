package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	type Message struct {
		MessageType string
		MessageText string
	}

	var name string

	conn, err := net.Dial("tcp", "localhost:8090")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		name = enterName()
		name = strings.TrimSuffix(name, "\n")

		if name != "" {
			break
		} else {
			fmt.Println("Name cannot be empty")
		}
	}

	msg := &Message{
		MessageType: "name",
		MessageText: strings.TrimSuffix(name, "\n"),
	}

	body, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Fprintf(conn, string(body))
	}

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Print("\033[1A")
			fmt.Print("\n" + message)
		}

		go func() {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print(">")
			text, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
			}

			fmt.Print("\033[1A") // cursor gooes up by 1 line and overwrites current input with message from server

			msg = &Message{
				MessageType: "message",
				MessageText: strings.TrimSuffix(text, "\n"),
			}

			body, err = json.Marshal(msg)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Fprintf(conn, string(body))
			}
		}()
	}
}

func enterName() string {
	nameReader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	name, err := nameReader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}

	return name
}
