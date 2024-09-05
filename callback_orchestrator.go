package socket

import "fmt"

var cborc CallbackOrchestrator = CallbackOrchestrator{}

type IEvent interface {}

type EventCallback func (IEvent) string

type CallbackOrchestrator struct {
	events map[string][]EventCallback
	ServerFd int
	Fds []int
}

func GetCBInstance() *CallbackOrchestrator { return &cborc }

func (orc *CallbackOrchestrator) Init() {
	orc.events = make(map[string][]EventCallback)
	orc.Fds = []int{}
}

func (orc *CallbackOrchestrator) Register(key string, cb EventCallback) {
	orc.events[key] = append(orc.events[key], cb)
}

func (orc *CallbackOrchestrator) Call(key string, e Event) {
	var displayMessage string
	for _, cb := range orc.events[key] {
		displayMessage = cb(e)
		if displayMessage != "" {
			fmt.Printf("%s\n", displayMessage)
		}
	}
}