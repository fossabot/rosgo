// Code generated by ros-gen-go.
// source: FollowJointTrajectoryAction.msg
// DO NOT EDIT!
package control_msgs

import (
	"io"

	"github.com/ppg/rosgo/ros"
)

type _MsgFollowJointTrajectoryAction struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgFollowJointTrajectoryAction) Text() string {
	return t.text
}

func (t *_MsgFollowJointTrajectoryAction) Name() string {
	return t.name
}

func (t *_MsgFollowJointTrajectoryAction) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgFollowJointTrajectoryAction) NewMessage() ros.Message {
	m := new(FollowJointTrajectoryAction)

	return m
}

var (
	MsgFollowJointTrajectoryAction = &_MsgFollowJointTrajectoryAction{
		`# ====== DO NOT MODIFY! AUTOGENERATED FROM AN ACTION DEFINITION ======

FollowJointTrajectoryActionGoal action_goal
FollowJointTrajectoryActionResult action_result
FollowJointTrajectoryActionFeedback action_feedback
`,
		"control_msgs/FollowJointTrajectoryAction",
		"bb7bd501b71fb425f1cb609216a06e3b",
	}
)

type FollowJointTrajectoryAction struct {
	ActionGoal     FollowJointTrajectoryActionGoal
	ActionResult   FollowJointTrajectoryActionResult
	ActionFeedback FollowJointTrajectoryActionFeedback
}

func (m *FollowJointTrajectoryAction) Serialize(w io.Writer) (err error) {
	if err = ros.SerializeMessageField(w, "FollowJointTrajectoryActionGoal", &m.ActionGoal); err != nil {
		return err
	}

	if err = ros.SerializeMessageField(w, "FollowJointTrajectoryActionResult", &m.ActionResult); err != nil {
		return err
	}

	if err = ros.SerializeMessageField(w, "FollowJointTrajectoryActionFeedback", &m.ActionFeedback); err != nil {
		return err
	}

	return
}

func (m *FollowJointTrajectoryAction) Deserialize(r io.Reader) (err error) {
	// ActionGoal
	if err = ros.DeserializeMessageField(r, "FollowJointTrajectoryActionGoal", &m.ActionGoal); err != nil {
		return err
	}

	// ActionResult
	if err = ros.DeserializeMessageField(r, "FollowJointTrajectoryActionResult", &m.ActionResult); err != nil {
		return err
	}

	// ActionFeedback
	if err = ros.DeserializeMessageField(r, "FollowJointTrajectoryActionFeedback", &m.ActionFeedback); err != nil {
		return err
	}

	return
}
