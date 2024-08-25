package main

import (
	"github.com/sota70/socket_framework/event"
	"github.com/sota70/socket_framework/server"
)


func registerEvents() {
	var orc *event.Orchestrator = event.GetInstance()
	orc.Init()

	var joinEvent event.PlayerJoinEvent = event.PlayerJoinEvent{}
	var listener event.PlayerJoinEventListener = event.PlayerJoinEventListener{
		E: &joinEvent,
	}
	joinEvent.Register(&listener)

	var recvEvent event.ServerRecvMsgEvent = event.ServerRecvMsgEvent{}
	var recvEventListener event.ServerRecvMsgEventListener = event.ServerRecvMsgEventListener{
		E: &recvEvent,
	}
	recvEvent.Register(&recvEventListener)

	var leaveEvent event.PlayerLeaveEvent = event.PlayerLeaveEvent{}
	var leaveEventListener event.PlayerLeaveEventListener = event.PlayerLeaveEventListener{
		E: &leaveEvent,
	}
	leaveEvent.Register(&leaveEventListener)

	var inputEvent event.ServerInputEvent = event.ServerInputEvent{}
	var inputEventListener event.ServerInputEventListener = event.ServerInputEventListener{
		E: &inputEvent,
	}
	inputEvent.Register(&inputEventListener)

	orc.Register("player_join", &joinEvent)
	orc.Register("player_leave", &leaveEvent)
	orc.Register("recv_msg", &recvEvent)
	orc.Register("input", &inputEvent)
}

func main() {
	registerEvents()
	server.Run([4]byte{127, 0, 0, 1}, 1337, event.GetInstance())
}