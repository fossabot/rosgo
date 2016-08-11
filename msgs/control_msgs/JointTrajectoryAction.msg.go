// Code generated by ros-gen-go.
// source: JointTrajectoryAction.msg
// DO NOT EDIT!
package control_msgs

import (
	"io"

	"github.com/ppg/rosgo/ros"
)

type _MsgJointTrajectoryAction struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgJointTrajectoryAction) Text() string {
	return t.text
}

func (t *_MsgJointTrajectoryAction) Name() string {
	return t.name
}

func (t *_MsgJointTrajectoryAction) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgJointTrajectoryAction) NewMessage() ros.Message {
	m := new(JointTrajectoryAction)

	return m
}

var (
	MsgJointTrajectoryAction = &_MsgJointTrajectoryAction{
		`# ====== DO NOT MODIFY! AUTOGENERATED FROM AN ACTION DEFINITION ======

JointTrajectoryActionGoal action_goal
JointTrajectoryActionResult action_result
JointTrajectoryActionFeedback action_feedback
`,
		"control_msgs/JointTrajectoryAction",
		"42ed1b0cbf2b5e2d85cea28d528e5774",
	}
)

type JointTrajectoryAction struct {
	ActionGoal     JointTrajectoryActionGoal
	ActionResult   JointTrajectoryActionResult
	ActionFeedback JointTrajectoryActionFeedback
}

func (m *JointTrajectoryAction) Serialize(w io.Writer) (err error) {
	if err = ros.SerializeMessageField(w, "JointTrajectoryActionGoal", &m.ActionGoal); err != nil {
		return err
	}

	if err = ros.SerializeMessageField(w, "JointTrajectoryActionResult", &m.ActionResult); err != nil {
		return err
	}

	if err = ros.SerializeMessageField(w, "JointTrajectoryActionFeedback", &m.ActionFeedback); err != nil {
		return err
	}

	return
}

func (m *JointTrajectoryAction) Deserialize(r io.Reader) (err error) {
	// ActionGoal
	if err = ros.DeserializeMessageField(r, "JointTrajectoryActionGoal", &m.ActionGoal); err != nil {
		return err
	}

	// ActionResult
	if err = ros.DeserializeMessageField(r, "JointTrajectoryActionResult", &m.ActionResult); err != nil {
		return err
	}

	// ActionFeedback
	if err = ros.DeserializeMessageField(r, "JointTrajectoryActionFeedback", &m.ActionFeedback); err != nil {
		return err
	}

	return
}