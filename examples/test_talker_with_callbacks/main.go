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
	pub := node.NewPublisherWithCallbacks("/chatter", msgs.MsgHello, onConnect, onDisconnect)

	for node.OK() {
		node.SpinOnce()
		var msg msgs.Hello
		msg.Data = fmt.Sprintf("hello %s", time.Now().String())
		fmt.Println(msg.Data)
		pub.Publish(&msg)
		time.Sleep(time.Second)
	}
}

func onConnect(pub ros.SingleSubscriberPublisher) {
	fmt.Printf("-------Connect callback: node %s topic %s\n", pub.GetSubscriberName(), pub.GetTopic())
	var msg msgs.Hello
	msg.Data = fmt.Sprintf("hello %s", pub.GetSubscriberName())
	pub.Publish(&msg)
}

func onDisconnect(pub ros.SingleSubscriberPublisher) {
	fmt.Printf("-------Disconnect callback: node %s topic %s\n", pub.GetSubscriberName(), pub.GetTopic())
}
