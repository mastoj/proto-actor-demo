package supervision

import (
	"fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type childActor struct{ childID int }
type parentActor struct{}

func (state *childActor) Receive(context actor.Context) {
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		fmt.Printf("Child %d: playing %d\n", state.childID, i)
	}
	panic("I don't want to play anymore")
}

func (state *parentActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		for i := 0; i < 2; i++ {
			childProps := actor.FromProducer(func() actor.Actor { return &childActor{childID: i} })
			_ = actor.Spawn(childProps)
		}
		fmt.Printf("Hello from actor2: %v\n", msg)
	}
}

func Run() {
	props := actor.FromProducer(func() actor.Actor { return &parentActor{} })
	_ = actor.Spawn(props)
}
