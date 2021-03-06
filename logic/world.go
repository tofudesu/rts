package logic

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/examples/demo/tadpole/logic/protocol"
	"github.com/lonng/nano/session"
)

// World contains all tadpoles
type Chunk struct {
	component.Base
	*nano.Group
}

// NewWorld returns a world instance
func NewChunk() *Chunk {
	return &Chunk{
		Group: nano.NewGroup(uuid.New().String()),
	}
}

// Init initialize world component
func (w *Chunk) Init() {
	session.Lifetime.OnClosed(func(s *session.Session) {
		w.Leave(s)
		w.Broadcast("leave", &protocol.LeaveWorldResponse{ID: s.ID()})
		log.Println(fmt.Sprintf("session count: %d", w.Count()))
	})
}

// Enter was called when new guest enter
func (w *Chunk) Enter(s *session.Session, msg []byte) error {
	w.Add(s)
	log.Println(fmt.Sprintf("session count: %d", w.Count()))
	return s.Response(&protocol.EnterWorldResponse{ID: s.ID()})
}

// Update refresh tadpole's position
func (w *Chunk) Update(s *session.Session, msg []byte) error {
	return w.Broadcast("update", msg)
}

// Message handler was used to communicate with each other
func (w *Chunk) Message(s *session.Session, msg *protocol.WorldMessage) error {
	msg.ID = s.ID()
	return w.Broadcast("message", msg)
}
