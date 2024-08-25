package socket

var orc Orchestrator = Orchestrator{}

type Orchestrator struct {
	events map[string]Event
	ServerFd int
	Fds []int
}

func GetInstance() *Orchestrator { return &orc }

func (orc *Orchestrator) Init() {
	orc.events = make(map[string]Event)
	orc.Fds = []int{}
}

func (orc *Orchestrator) Register(key string, e Event) {
	orc.events[key] = e
}

func (orc *Orchestrator) Call(key string, e Event) {
	orc.events[key].Update(e)
}