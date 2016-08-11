package main

import (
	"fmt"
	"time"

	"github.com/ppg/rosgo/examples/msg"
	"github.com/ppg/rosgo/ros"
)

func main() {
	node := ros.NewNode("/talker")
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelDebug)
	pub := node.NewPublisher("/chatter", msgs.MsgHello)

	for node.OK() {
		node.SpinOnce()
		var msg msgs.Hello
		msg.Data = fmt.Sprintf("hello %s", time.Now().String())
		fmt.Println(msg.Data)
		pub.Publish(&msg)
		time.Sleep(time.Second)
	}
}
