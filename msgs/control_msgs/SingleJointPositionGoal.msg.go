// Code generated by ros-gen-go.
// source: SingleJointPositionGoal.msg
// DO NOT EDIT!
package control_msgs

import (
	"io"

	"github.com/ppg/rosgo/ros"
)

type _MsgSingleJointPositionGoal struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgSingleJointPositionGoal) Text() string {
	return t.text
}

func (t *_MsgSingleJointPositionGoal) Name() string {
	return t.name
}

func (t *_MsgSingleJointPositionGoal) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgSingleJointPositionGoal) NewMessage() ros.Message {
	m := new(SingleJointPositionGoal)

	return m
}

var (
	MsgSingleJointPositionGoal = &_MsgSingleJointPositionGoal{
		`# ====== DO NOT MODIFY! AUTOGENERATED FROM AN ACTION DEFINITION ======
float64 position
duration min_duration
float64 max_velocity
`,
		"control_msgs/SingleJointPositionGoal",
		"559179a05fec5a369c077fac92e60034",
	}
)

type SingleJointPositionGoal struct {
	Position    float64
	MinDuration ros.Duration
	MaxVelocity float64
}

func (m *SingleJointPositionGoal) Serialize(w io.Writer) (err error) {
	if err = ros.SerializeMessageField(w, "float64", &m.Position); err != nil {
		return err
	}

	if err = ros.SerializeMessageField(w, "duration", &m.MinDuration); err != nil {
		return err
	}

	if err = ros.SerializeMessageField(w, "float64", &m.MaxVelocity); err != nil {
		return err
	}

	return
}

func (m *SingleJointPositionGoal) Deserialize(r io.Reader) (err error) {
	// Position
	if err = ros.DeserializeMessageField(r, "float64", &m.Position); err != nil {
		return err
	}

	// MinDuration
	if err = ros.DeserializeMessageField(r, "duration", &m.MinDuration); err != nil {
		return err
	}

	// MaxVelocity
	if err = ros.DeserializeMessageField(r, "float64", &m.MaxVelocity); err != nil {
		return err
	}

	return
}
