package event

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

type ServerInputEvent struct {
	DisplayMessage string
	Input string
	NeedsOutput bool
	listeners []Listener
}

type ServerInputEventListener struct {
	E *ServerInputEvent
}

func (e *ServerInputEvent) Register(listener Listener) {
	e.listeners = append(e.listeners, listener)
}

func (e *ServerInputEvent) Update(new Event) {
	if newEvent, ok := new.(*ServerInputEvent); ok {
		e.Input = newEvent.Input
		e.NeedsOutput = newEvent.NeedsOutput
	}
	for i := 0; i < len(e.listeners); i++ {
		e.listeners[i].Listen()
		if e.NeedsOutput {
			e.Render()
		}
	}
}

func (e *ServerInputEvent) Render() {
	for _, fd := range GetInstance().Fds {
		unix.Send(fd, []byte(e.DisplayMessage), 0)
	}
}

func (listener *ServerInputEventListener) Listen() {
	if listener.E.Input == "q" || listener.E.Input == "quit" {
		for fd := range GetInstance().Fds {
			unix.Close(fd)
		}
		listener.E.NeedsOutput = false
		unix.Close(GetInstance().ServerFd)
		os.Exit(0)
		return
	}
	listener.E.DisplayMessage = fmt.Sprintf("[server] > %s\n", listener.E.Input)
}