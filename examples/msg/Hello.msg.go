// Code generated by ros-gen-go.
// source: Hello.msg
// DO NOT EDIT!
package msgs

import (
	"io"

	"github.com/ppg/rosgo/ros"
)

type _MsgHello struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgHello) Text() string {
	return t.text
}

func (t *_MsgHello) Name() string {
	return t.name
}

func (t *_MsgHello) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgHello) NewMessage() ros.Message {
	m := new(Hello)

	return m
}

var (
	MsgHello = &_MsgHello{
		`string data
`,
		"msgs/Hello",
		"5b89344b4d080ca499c0648063c40738",
	}
)

type Hello struct {
	Data string
}

func (m *Hello) Type() ros.MessageType {
	return MsgHello
}

func (m *Hello) Serialize(w io.Writer) (err error) {
	if err = ros.SerializeMessageField(w, "string", &m.Data); err != nil {
		return err
	}

	return
}

func (m *Hello) Deserialize(r io.Reader) (err error) {
	// Data
	if err = ros.DeserializeMessageField(r, "string", &m.Data); err != nil {
		return err
	}

	return
}
