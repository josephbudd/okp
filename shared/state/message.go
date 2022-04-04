package state

import "time"

// Message informs the front process of state changes.
type Message struct {
	Time              int64
	NewCourse         bool
	PassedCopying     bool
	PassedKeying      bool
	CompletedCopying  bool
	CompletedKeying   bool
	CompletedHomeWork bool
	CompletedCourse   bool
}

type StateWith func(msg *Message)

// NewMessage constructs a new State message.
func NewMessage(with ...StateWith) (msg *Message) {
	msg = &Message{
		Time: time.Now().UnixMicro(),
	}
	for _, w := range with {
		w(msg)
	}
	return
}

// StateWithNewCourse signals that the current course is new or completed.
// The backend will use with NewState.
func StateWithNewCourse() (f func(msg *Message)) {
	f = func(msg *Message) {
		msg.NewCourse = true
	}
	return
}

// StateWithPassedCopying signals that in the current lesson, copying has leveled up.
// The backend will use with NewState.
func StateWithPassedCopying() (f func(msg *Message)) {
	f = func(msg *Message) {
		msg.PassedCopying = true
	}
	return
}

// StateWithPassedKeying signals that in the current lesson, keying has leveled up.
// The backend will use with NewState.
func StateWithPassedKeying() (f func(msg *Message)) {
	f = func(msg *Message) {
		msg.PassedKeying = true
	}
	return
}

// StateWithCompletedCopying signals that in the current lesson, copying has leveled up.
// The backend will use with NewState.
func StateWithCompletedCopying() (f func(msg *Message)) {
	f = func(msg *Message) {
		msg.CompletedCopying = true
	}
	return
}

// StateWithCompletedKeying signals that in the current lesson, copying has leveled up.
// The backend will use with NewState.
func StateWithCompletedKeying() (f func(msg *Message)) {
	f = func(msg *Message) {
		msg.CompletedKeying = true
	}
	return
}

// StateWithCompletedHomeWork signals that in the current course, the homework has leveled up to a new lesson.
// The backend will use with NewState.
func StateWithCompletedHomeWork() (f func(msg *Message)) {
	f = func(msg *Message) {
		msg.CompletedHomeWork = true
	}
	return
}

// StateWithCompletedCourse signals that in the current course, the homework has leveled up to a new lesson.
// The backend will use with NewState.
func StateWithCompletedCourse() (f func(msg *Message)) {
	f = func(msg *Message) {
		msg.CompletedCourse = true
	}
	return
}
