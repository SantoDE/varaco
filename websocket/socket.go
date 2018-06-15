package websocket

import (
	"github.com/SantoDE/varaco/types"
	"golang.org/x/net/websocket"
	"fmt"
)

func DoSocketCall(cmd types.ExecuteCommand) {
	_, err := websocket.Dial(fmt.Sprintf("%s?token=%s", cmd.Url, cmd.Token), "", cmd.Url)

	if err != nil {
		fmt.Printf("Error Dialing WS %s \n", err.Error())
	}

	fmt.Printf("Dialed success! \n")
}
