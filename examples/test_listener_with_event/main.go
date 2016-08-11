package main

import (
	"fmt"

	"github.com/ppg/rosgo/examples/msg"
	"github.com/ppg/rosgo/ros"
)

func callback(msg *msgs.Hello, event ros.MessageEvent) {
	fmt.Printf("Received: %s from %s, header = %v, time = %v\n",
		msg.Data, event.PublisherName, event.ConnectionHeader, event.ReceiptTime)
}

func main() {
	node := ros.NewNode("/listener")
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelDebug)
	node.NewSubscriber("/chatter", msgs.MsgHello, callback)
	node.Spin()
}
