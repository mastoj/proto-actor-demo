package worker

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/hashicorp/consul/api"
	"github.com/unacast/proto-actor-demo/go-node/messages"
)

type workerActor struct {
	workerID  int
	masterPid *actor.PID
}
type workerActorMonitor struct {
	workerCount int
	masterPid   *actor.PID
}

func startCluster(hostName string) {
	cp, err := consul.NewWithConfig(&api.Config{Address: "http://consul:8500"})
	if err != nil {
		log.Fatal(err)
	}
	cluster.Start("FiresideChatCluster", hostName+":12000", cp)
}

func getMasterPid() *actor.PID {
	masterPid, sc := cluster.Get("FiresideChatCluster", "MasterKind")
	for ; sc != remote.ResponseStatusCodeOK; masterPid, sc = cluster.Get("FiresideChatCluster", "MasterKind") {
		fmt.Println("Failed to get master pid, trying again in one second")
		time.Sleep(1 * time.Second)
	}
	return masterPid
}

func requestWork(masterPid *actor.PID, workerPid *actor.PID) {
	message := messages.RequestWork{}
	messagePid := messages.PID{}
	messagePid.Address = workerPid.Address
	messagePid.Id = workerPid.Id
	message.Pid = &messagePid
	masterPid.Tell(&message)
}

func (worker *workerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.SubmitWork:
		fmt.Printf("GoWorker %d, got some work: %v\n", worker.workerID, msg)
		time.Sleep(1500 * time.Millisecond)
		value := float64(msg.Data)
		result := math.Sqrt(value)
		message := messages.SubmitResult{}
		message.Data = msg.Data
		message.Result = float32(result)
		sinkPid := actor.PID{Address: msg.Pid.Address, Id: msg.Pid.Id}
		sinkPid.Tell(&message)
		requestWork(worker.masterPid, context.Self())
	case *actor.Started:
		fmt.Printf("GoWorker %d started, asking for work from master: %v", worker.workerID, worker.masterPid)
		requestWork(worker.masterPid, context.Self())
		fmt.Printf("GoWorker %d Asked for my first piece of work", worker.workerID)
	}
}

func (monitor *workerActorMonitor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		for i := 0; i < monitor.workerCount; i++ {
			childProps := actor.FromProducer(func() actor.Actor { return &workerActor{workerID: i, masterPid: monitor.masterPid} })
			_ = context.Spawn(childProps)
		}
	}
}

func startWorkerMonitor(masterPid *actor.PID, workerCount int) {
	monitorProps := actor.FromProducer(func() actor.Actor { return &workerActorMonitor{workerCount: workerCount, masterPid: masterPid} })
	_ = actor.Spawn(monitorProps)
}

func Run(workerCount int) {
	fmt.Printf("Starting go worker: %d\n", workerCount)
	time.Sleep(3 * time.Second)
	hostName, _ := os.Hostname()

	startCluster(hostName)
	masterPid := getMasterPid()
	fmt.Printf("Will start worker monitor with args: %v, %v, %v\n", hostName, workerCount, *masterPid)
	startWorkerMonitor(masterPid, workerCount)
}
