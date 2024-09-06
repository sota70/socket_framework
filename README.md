# Overview

This is a socket server framework that is written in Go.<br>
The framework is a event-driven framework which you don't have to care about where to write code and when to execute code.

# Installation

Write this in go.mod<br>
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

The framework consists of three elements.<br>
(1)Event, (2)Callback, and (3)Orchestrator.<br>
Event is a class that is called in the framework's logic and has information about event.<br>
Callback is a function that handles specific event.<br>
More than one callbacks can be added to one event.<br>
Orchestrator is a class that manages all events, callbacks and common states in the server.<br>
This framework has several pre-built events.<br>
For instance, ServerRecvMsgEvent.<br>
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

First, initialize Orchestrator.<br>
Second, register callback to event.<br>
Finally, register events to Orchestrator with an event name.<br>
All events have DisplayMessage property which is rendered in stdout every time a callback handles an event.<br>
If the event has 3 listeners and they all set some message to DisplayMessage, 3 different messages are rendered in stdout.

## How to implement your own callbacks

In this section, you will learn how to implement your own callbacks.<br>
There are 4 callbacks you must implement.<br>
- PlayerJoin callback<br>
- PlayerLeave callback<br>
- ServerInput callback<br>
- RecvMsg callback<br>

PlayerJoin callback is called when a new client joins the server.<br>
PlayerLeave callback is called when a client leaves the server.<br>
ServerInput callback is called when the server receives user input from server admin.<br>
RecvMsg callback is called when the server receives a message from a client.

### Make your own PlayerJoin callback

The following code is creating new PlayerJoin callback.
```go
var cborc = socket_framework.GetCBInstance()
cborc.Register("player_join", func (e socket_framework.IEvent) string {
	if event, ok := e.(*socket_framework.PlayerJoinEvent); ok {
		socket_framework.GetCBInstance().Fds = append(socket_framework.GetCBInstance().Fds, event.NewFd)
		return fmt.Sprintf("%d has joined the server", event.NewFd)
	}
	return ""
})
```

### Make your own PlayerLeave callback

The following code is creating new PlayerLeave callback.
```go
var cborc = socket_framework.GetCBInstance()
cborc.Register("player_leave", func (e socket_framework.IEvent) string {
	if event, ok := e.(*socket_framework.PlayerLeaveEvent); ok {
		var cborc = socket_framework.GetCBInstance()
		unix.Close(event.LeftFd)
		for i, fd := range cborc.Fds {
			if fd == event.LeftFd {
				cborc.Fds = append(cborc.Fds[:i], cborc.Fds[i + 1:]...)
				break
			}
		}
		return fmt.Sprintf("%d has left the server", event.LeftFd)
	}
	return ""
})
```

### Make your own ServerInput callback

The following code is creating new ServerInput callback.
```go
var cborc = socket_framework.GetCBInstance()
cborc.Register("input", func (e socket_framework.IEvent) string {
	if event, ok := e.(*socket_framework.ServerInputEvent); ok {
		var cborc = socket_framework.GetCBInstance()
		if event.Input == "q" || event.Input == "quit" {
			for _, fd := range cborc.Fds {
				unix.Close(fd)
			}
			unix.Close(cborc.ServerFd)
			os.Exit(0)
			return ""
		}
		for _, fd := range cborc.Fds {
			unix.Send(fd, []byte(fmt.Sprintf("[server] > %s\n", event.Input)), 0)
		}
		return ""
	}
	return ""
}
```

### Make your own RecvMsg callback

The following code is creating new RecvMsg callback.
```go
var cborc = socket_framework.GetCBInstance()
cborc.Register("recv_msg", func (e socket_framework.IEvent) string {
	if event, ok := e.(*socket_framework.ServerRecvMsgEvent); ok {
		return fmt.Sprintf(
			"[%d] > %s",
			event.Src,
			strings.ReplaceAll(event.RecvMsg, "\n", ""),
		)
	}
	return ""
})
```

## Sample Program

This is a sample server program using the framework.
```go
package main

import (
	"fmt"
	"os"
	"strings"

	f "github.com/sota70/socket_framework"
	"golang.org/x/sys/unix"
)

func registerCBEvents() {
	var cborc = f.GetCBInstance()
	cborc.Register("player_join", func (e f.IEvent) string {
		if event, ok := e.(*f.PlayerJoinEvent); ok {
			f.GetCBInstance().Fds = append(f.GetCBInstance().Fds, event.NewFd)
			return fmt.Sprintf("%d has joined the server", event.NewFd)
		}
		return ""
	})
	cborc.Register("player_leave", func (e f.IEvent) string {
		if event, ok := e.(*f.PlayerLeaveEvent); ok {
			var cborc = f.GetCBInstance()
			unix.Close(event.LeftFd)
			for i, fd := range cborc.Fds {
				if fd == event.LeftFd {
					cborc.Fds = append(cborc.Fds[:i], cborc.Fds[i + 1:]...)
					break
				}
			}
			return fmt.Sprintf("%d has left the server", event.LeftFd)
		}
		return ""
	})
	cborc.Register("input", func (e f.IEvent) string {
		if event, ok := e.(*f.ServerInputEvent); ok {
			var cborc = f.GetCBInstance()
			if event.Input == "q" || event.Input == "quit" {
				for _, fd := range cborc.Fds {
					unix.Close(fd)
				}
				unix.Close(cborc.ServerFd)
				os.Exit(0)
				return ""
			}
			for _, fd := range cborc.Fds {
				unix.Send(fd, []byte(fmt.Sprintf("[server] > %s\n", event.Input)), 0)
			}
			return ""
		}
		return ""
	})
	cborc.Register("recv_msg", func (e f.IEvent) string {
		if event, ok := e.(*f.ServerRecvMsgEvent); ok {
			return fmt.Sprintf(
				"[%d] > %s",
				event.Src,
				strings.ReplaceAll(event.RecvMsg, "\n", ""),
			)
		}
		return ""
	})
}

func main() {
	var cborc = f.GetCBInstance()
	cborc.Init()
	registerCBEvents()
	f.Run([4]byte{127, 0, 0, 1}, 1337, cborc)
}
```
