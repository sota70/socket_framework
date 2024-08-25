package socket


func registerEvents() {
	var orc *Orchestrator = GetInstance()
	orc.Init()

	var joinEvent PlayerJoinEvent = PlayerJoinEvent{}
	var listener PlayerJoinEventListener = PlayerJoinEventListener{
		E: &joinEvent,
	}
	joinEvent.Register(&listener)

	var recvEvent ServerRecvMsgEvent = ServerRecvMsgEvent{}
	var recvEventListener ServerRecvMsgEventListener = ServerRecvMsgEventListener{
		E: &recvEvent,
	}
	recvEvent.Register(&recvEventListener)

	var leaveEvent PlayerLeaveEvent = PlayerLeaveEvent{}
	var leaveEventListener PlayerLeaveEventListener = PlayerLeaveEventListener{
		E: &leaveEvent,
	}
	leaveEvent.Register(&leaveEventListener)

	var inputEvent ServerInputEvent = ServerInputEvent{}
	var inputEventListener ServerInputEventListener = ServerInputEventListener{
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
	Run([4]byte{127, 0, 0, 1}, 1337, GetInstance())
}