package main

import (
	"fmt"

	"github.com/ppg/rosgo/examples/msg"
	"github.com/ppg/rosgo/ros"
)

func callback(msg *msgs.Hello) {
	fmt.Printf("Received: %s\n", msg.Data)
}

func main() {
	node := ros.NewNode("/listener")
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelDebug)
	node.NewSubscriber("/chatter", msgs.MsgHello, callback)
	node.Spin()
}
