package event

import (
	"fmt"
	"strings"
)

type ServerRecvMsgEvent struct {
	DisplayMessage string
	Src int
	Dst int
	RecvMsg string
	listeners []Listener
}

type ServerRecvMsgEventListener struct {
	E *ServerRecvMsgEvent
}

func (e *ServerRecvMsgEvent) Register(listener Listener) {
	e.listeners = append(e.listeners, listener)
}

func (e *ServerRecvMsgEvent) Update(new Event) {
	if newEvent, ok := new.(*ServerRecvMsgEvent); ok {
		e.Src = newEvent.Src
		e.Dst = newEvent.Dst
		e.RecvMsg = newEvent.RecvMsg
	}
	for i := 0; i < len(e.listeners); i++ {
		e.listeners[i].Listen()
		e.Render()
	}
}

func (e *ServerRecvMsgEvent) Render() {
	fmt.Printf("%s\n", e.DisplayMessage)
	e.DisplayMessage = ""
}

func (listener *ServerRecvMsgEventListener) Listen() {
	listener.E.DisplayMessage = fmt.Sprintf(
		"[%d] > %s",
		listener.E.Src,
		strings.ReplaceAll(listener.E.RecvMsg, "\n", ""),
	)
}