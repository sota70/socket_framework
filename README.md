# Overview

This is a socket server framework that is written in Go.
The framework is a event-driven framework which you don't have to care about where to write code and when to execute code.

# Installation

Write this in go.mod
Execute this command in the root folder of your project.
```
go get github.com/sota70/socket_framework
```

To use this, you can import in import section
```go
import (
	f "github.com/sota70/socket_framework"
)
```

# How this framework works

The framework consists of three elements.
(1)Event, (2)Callback, and (3)Orchestrator.
Event is a class that is called in the framework's logic and has information about event.
Callback is a function that handles specific event.
More than one callbacks can be added to one event.
Orchestrator is a class that manages all events, callbacks and common states in the server.
This framework has several pre-built events.
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
Second, register callback to event.
Finally, register events to Orchestrator with an event name.
All events have DisplayMessage property which is rendered in stdout every time a callback handles an event.
If the event has 3 listeners and they all set some message to DisplayMessage, 3 different messages are rendered in stdout.

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

### Make your own PlayerJoin listener

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

### Make your own PlayerLeave listener

The following code is creating new PlayerLeave listener class and is registering it to PlayerLeaveEvent.
```go
package main

import (
	"fmt"

	f "github.com/sota70/socket_framework"
)

type SamplePlayerLeaveListener struct {
	Event *f.PlayerLeaveEvent
}

func (listener *SamplePlayerLeaveListener) Listen() {
        listener.Event.DisplayMessage = fmt.Sprintf("[LOG]Client %d has left the server", listener.Event.LeftFd)
}

func main() {
  f.GetInstance().Init()
  var leaveEvent f.PlayerLeaveEvent = f.PlayerLeaveEvent {}
  var leaveListener SamplePlayerLeaveListener = SamplePlayerLeaveListener {
    Event: &leaveEvent,
  }
  leaveEvent.Register(&leaveListener)

  f.GetInstance().Register("leave", &leaveEvent)
}
```

### Make your own ServerRecvMsgEvent listener

The following code is creating new ServerRecvMsgEvent listener class and is registering it to ServerRecvMsgEvent.
```go
package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	f "github.com/sota70/socket_framework"
)

type SampleServerRecvMsgListener struct {
	Event *f.ServerRecvMsgEvent
}

func (listener *SampleServerRecvMsgListener) Listen() {
	listener.Event.DisplayMessage = fmt.Sprintf(
		"[Sample][%d] > %s",
		listener.Event.Src,
		strings.ReplaceAll(listener.Event.RecvMsg, "\n", ""),
	)
}

func main() {
	f.GetInstance().Init()
	var recvEvent f.ServerRecvMsgEvent = f.ServerRecvMsgEvent{}
	var sampleRecvListener = SampleServerRecvMsgListener{
		Event: &recvEvent,
	}
	recvEvent.Register(&sampleRecvListener)
	orc.Register("recv_msg", &recvEvent)
}
```

### Make your own ServerInputEvent listener

The following code is creating new ServerInputEvent listener class and is registering it to ServerInputEvent.
```go
package main

import (
	"fmt"
	"os"
	"golang.org/x/sys/unix"
	f "github.com/sota70/socket_framework"
)

type SampleInputListener {
	Event *f.ServerInputEvent
}

func (listener *SampleInputListener) Listen() {
	if listener.Event.Input == "s" || listener.E.Input == "stop" {
		for fd := range f.GetInstance().Fds {
			unix.Close(fd)
		}
		listener.Event.NeedsOutput = false
		unix.Close(f.GetInstance().ServerFd)
		os.Exit(0)
		return
	}
	listener.Event.DisplayMessage = fmt.Sprintf("[Sample][server] > %s\n", listener.Event.Input)
}

func main() {
	var inputEvent f.ServerInputEvent = f.ServerInputEvent{}
	var sampleInputListener f.SampleInputListener = SampleInputListener{
		Event: &inputEvent,
	}
	inputEvent.Register(&sampleInputListener)
	f.GetInstance().Register("input", &inputEvent)
}
```
