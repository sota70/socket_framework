package event

import (
	"fmt"
)

type PlayerJoinEvent struct {
	DisplayMessage string
	NewFd int
	listeners []Listener
}

type PlayerJoinEventListener struct {
	E *PlayerJoinEvent
}

func (e *PlayerJoinEvent) Register(listener Listener) {
	e.listeners = append(e.listeners, listener)
}

func (e *PlayerJoinEvent) Update(new Event) {
	if newEvent, ok := new.(*PlayerJoinEvent); ok {
		e.NewFd = newEvent.NewFd
	}
	for i := 0; i < len(e.listeners); i++ {
		e.listeners[i].Listen()
		e.Render()
	}
}

func (e *PlayerJoinEvent) Render() {
	fmt.Printf("%s\n", e.DisplayMessage)
	e.DisplayMessage = ""
}

func (listener *PlayerJoinEventListener) Listen() {
	listener.E.DisplayMessage = fmt.Sprintf("%d has joined the server", listener.E.NewFd)
	GetInstance().Fds = append(GetInstance().Fds, listener.E.NewFd)
}