package socket

type Event interface {
	Register(Listener)
	Update(Event)
	Render()
}

type Listener interface {
	Listen()
}