package server

import (
	"fmt"

	"golang.org/x/sys/unix"
	"haxtac.org/socket_framework/event"
)

const (
	ZERO_BYTE_ERR = "0 bytes received"
	SOCKET_OP_ERR = "socket operation on non-socket"
	BAD_FD = "bad file descriptor"
	MAX_BUF_SIZE = 1024
)


func recv(clientFd int, max_buf_size int, orc *event.Orchestrator) {
	var buf []byte
	var readLen int
	var err error

	for {
		buf = make([]byte, max_buf_size)
		readLen, _, err = unix.Recvfrom(clientFd, buf, 0)
		if err != nil {
			orc.Call("player_leave", &event.PlayerLeaveEvent{
				LeftFd: clientFd,
				NeedsOutput: true,
			})
			return
		}
		// クライアントが通信を切断した際
		// サイズが0のデータを受信することになる
		// よってそれを検知し、クライアントとの通信を正式に切断する処理を行う
		if readLen < 1 {
			orc.Call("player_leave", &event.PlayerLeaveEvent{
				LeftFd: clientFd,
				NeedsOutput: true,
			})
			return
		}
		orc.Call("recv_msg", &event.ServerRecvMsgEvent{
			Src: clientFd,
			RecvMsg: string(buf),
		})
	}
}

func handleInput() {
	var input string
	var err error
	for {
		input = ""
		_, err = fmt.Scanln(&input)
		if err != nil {
			continue
		}
		event.GetInstance().Call("input", &event.ServerInputEvent{
			Input: input,
			NeedsOutput: true,
		})
	}
}

func runServer(host [4]byte, port int) (int, error) {
	var serverFd int
	var err error
	var addr unix.Sockaddr = &unix.SockaddrInet4{
		Addr: host,
		Port: port,
	}
	serverFd, err = unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_IP)
	if err != nil {
		return -1, err
	}
	// クライアントと通信をしている最中にサーバを閉じた後
	// 再度サーバを立ち上げようとすると、サーバが使用するソケットが使用中となっており
	// ソケットが使用できない
	// それを解決するために、ソケットに再利用フラグを立てている
	unix.SetsockoptInt(serverFd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
	err = unix.Bind(serverFd, addr)
	if err != nil {
		return -1, err
	}
	err = unix.Listen(serverFd, 100)
	if err != nil {
		return -1, err
	}
	return serverFd, nil
}

func Run(host [4]byte, port int, orc *event.Orchestrator) {
	var serverFd int
	var clientFd int
	var err error
	serverFd, err = runServer(host, port)
	if err != nil {
		fmt.Printf("An error occured during starting the server: %v\n", err)
		return
	}
	fmt.Printf("The server is running on port %d\n", port)
	// debug
	fmt.Printf("serverFd: %d\n", serverFd)

	orc.ServerFd = serverFd
	go handleInput()
	for {
		clientFd, _, err = unix.Accept(serverFd)
		if err != nil {
			continue
		}
		orc.Call("player_join", &event.PlayerJoinEvent{
			NewFd: clientFd,
		})
		go recv(clientFd, MAX_BUF_SIZE, orc)
	}
}