package socket

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

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

func registerCBEvents() {
	var cborc = GetCBInstance()
	cborc.Register("player_join", func (e IEvent) string {
		if event, ok := e.(*PlayerJoinEvent); ok {
			GetCBInstance().Fds = append(GetCBInstance().Fds, event.NewFd)
			return fmt.Sprintf("%d has joined the server", event.NewFd)
		}
		return ""
	})
	cborc.Register("player_leave", func (e IEvent) string {
		if event, ok := e.(*PlayerLeaveEvent); ok {
			var cborc = GetCBInstance()
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
	cborc.Register("input", func (e IEvent) string {
		if event, ok := e.(*ServerInputEvent); ok {
			var cborc = GetCBInstance()
			if event.Input == "q" || event.Input == "quit" {
				for fd := range cborc.Fds {
					unix.Close(fd)
				}
				unix.Close(cborc.ServerFd)
				os.Exit(0)
				return ""
			}
			for fd := range cborc.Fds {
				unix.Send(fd, []byte(fmt.Sprintf("[server] > %s\n", event.Input)), 0)
			}
			return ""
		}
		return ""
	})
	cborc.Register("recv_msg", func (e IEvent) string {
		if event, ok := e.(*ServerRecvMsgEvent); ok {
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
	// registerEvents()
	registerCBEvents()
	Run([4]byte{127, 0, 0, 1}, 1337, GetCBInstance())
}