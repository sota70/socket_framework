package socket

import (
	"fmt"

	"golang.org/x/sys/unix"
)

type PlayerLeaveEvent struct {
	DisplayMessage string
	NeedsOutput bool
	LeftFd int
	listeners []Listener
}

type PlayerLeaveEventListener struct {
	E *PlayerLeaveEvent
}

func (e *PlayerLeaveEvent) Register(listener Listener) {
	e.listeners = append(e.listeners, listener)
}

func (e *PlayerLeaveEvent) Update(new Event) {
	if newEvent, ok := new.(*PlayerLeaveEvent); ok {
		e.LeftFd = newEvent.LeftFd
		e.NeedsOutput = newEvent.NeedsOutput
	}
	for i := 0; i < len(e.listeners); i++ {
		e.listeners[i].Listen()
		if e.NeedsOutput {
			e.Render()
		}
	}
}

func (e *PlayerLeaveEvent) Render() {
	fmt.Printf("%s\n", e.DisplayMessage)
	e.DisplayMessage = ""
}

func (listener *PlayerLeaveEventListener) Listen() {
	unix.Close(listener.E.LeftFd)
	listener.E.DisplayMessage = fmt.Sprintf("%d has left the server", listener.E.LeftFd)
	// remove fd from fds
	for i, fd := range GetInstance().Fds {
		if fd == listener.E.LeftFd {
			GetInstance().Fds = append(GetInstance().Fds[:i], GetInstance().Fds[i + 1:]...)
			break
		}
	}
}