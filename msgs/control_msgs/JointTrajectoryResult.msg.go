// Code generated by ros-gen-go.
// source: JointTrajectoryResult.msg
// DO NOT EDIT!
package control_msgs

import (
	"io"

	"github.com/ppg/rosgo/ros"
)

type _MsgJointTrajectoryResult struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgJointTrajectoryResult) Text() string {
	return t.text
}

func (t *_MsgJointTrajectoryResult) Name() string {
	return t.name
}

func (t *_MsgJointTrajectoryResult) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgJointTrajectoryResult) NewMessage() ros.Message {
	m := new(JointTrajectoryResult)

	return m
}

var (
	MsgJointTrajectoryResult = &_MsgJointTrajectoryResult{
		`# ====== DO NOT MODIFY! AUTOGENERATED FROM AN ACTION DEFINITION ======
`,
		"control_msgs/JointTrajectoryResult",
		"7ac3b32c97133caf1b14edc99a50c37d",
	}
)

type JointTrajectoryResult struct {
}

func (m *JointTrajectoryResult) Serialize(w io.Writer) (err error) {
	return
}

func (m *JointTrajectoryResult) Deserialize(r io.Reader) (err error) {
	return
}
