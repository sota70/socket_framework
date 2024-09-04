# Overview

This is a socket server framework that is written in Go.
The framework is a event-driven framework which you don't have to care about where to write code and when to execute code.

# Installation

Write this in go.mod
```
require (
  github.com/sota70/socket_framework v0.0.1
)
```

# How this framework works

The framework consists of three elements.
(1)Event, (2)Listener, and (3)Orchestrator.
Event is a class that is called in the framework's logic and has information about event.
It can have several listeners.
Listener is a class that handles specific event.
It can only have one event.
Orchestrator is a class that manages all events, listeners and common states in the server.
This framework has several pre-built events and listeners.
For instance, ServerRecvMsgEvent.
For example, the process that is receiving message from a client calls ServerRecvMsgEvent in the following code.
```go
buf = make([]byte, max_buf_size)
readLen, _, err = unix.Recvfrom(clientFd, buf, 0)

...

orc.Call("recv_msg", &ServerRecvMsgEvent{
  Src: clientFd,
  RecvMsg: string(buf),
})
```

# How to use this framework

First, initialize Orchestrator.
Second, make listeners.
After that, register them to event.
Finally, register events to Orchestrator with an event name.

## Use this framework with Pre-built events

You can use pre-built events and listeners.
Register them to Orchestrator in the following code.
```go
package main

import (
	f "github.com/sota70/socket_framework"
)


f.GetInstance().Init()
var joinEvent f.PlayerJoinEvent = f.PlayerJoinEvent{}
var joinListener f.PlayerJoinEventListener = f.PlayerJoinEventListener{
	E: &joinEvent,
}
joinEvent.Register(&joinListener)

var leaveEvent f.PlayerLeaveEvent = f.PlayerLeaveEvent{}
var leaveListener f.PlayerLeaveEventListener = f.PlayerLeaveEventListener{
	E: &leaveEvent,
}
leaveEvent.Register(&leaveListener)

var recvEvent f.ServerRecvMsgEvent = f.ServerRecvMsgEvent{}
var recvListener f.ServerRecvMsgEventListener = f.ServerRecvMsgEventListener{
	E: &recvEvent,
}
recvEvent.Register(&recvListener)

var inputEvent f.ServerInputEvent = f.ServerInputEvent{}
var inputListener f.ServerInputEventListener = f.ServerInputEventListener{
	E: &inputEvent,
}
inputEvent.Register(&inputListener)

f.GetInstance().Register("player_join", &joinEvent)
f.GetInstance().Register("player_leave", &leaveEvent)
f.GetInstance().Register("recv_msg", &recvEvent)
f.GetInstance().Register("input", &inputEvent)
```
After you register them, run the server with following code.
```go
var host [4]byte = [4]byte{127, 0, 0, 1}
var port int = [port];
f.Run(host, port, f.GetInstance())
```

## Use your own listener

In the previous section, you use pre-built listeners.
However, you can also make your own listeners.
The following code is creating new PlayerJoin listener class and is registering it to PlayerJoinEvent.
```go
package main

import (
	"fmt"

	f "github.com/sota70/socket_framework"
)

type SamplePlayerJoinListener struct {
	Event *f.PlayerJoinEvent
}

func (listener *SamplePlayerJoinListener) Listen() {
	listener.Event.DisplayMessage = fmt.Sprintf("[LOG]Client %d has joined the server", listener.Event.NewFd)
}

func main() {
  f.GetInstance().Init()
  var joinEvent f.PlayerJoinEvent = f.PlayerJoinEvent {}
  var joinListener SamplePlayerJoinListener = SamplePlayerJoinListener {
    Event: &joinEvent,
  }
  joinEvent.Register(&joinListener)

  f.GetInstance().Register("join", &joinEvent)
}
```
