package hello

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type helloActor1 struct{}
type helloActor2 struct{ pid actor.PID }

func (state *helloActor1) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case string:
		fmt.Printf("Hello from actor1: %v\n", msg)
	}
}

func (state *helloActor2) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case string:
		fmt.Printf("Hello from actor2: %v\n", msg)
		state.pid.Tell(fmt.Sprintf("I was called with %v", msg))
	}
}

func Hello() {
	props1 := actor.FromProducer(func() actor.Actor { return &helloActor1{} })
	pid1 := actor.Spawn(props1)
	props2 := actor.FromProducer(func() actor.Actor { return &helloActor2{pid: *pid1} })
	pid2 := actor.Spawn(props2)
	pid2.Tell("Hello world")
}
