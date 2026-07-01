package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"slices"
	"strconv"
	"time"

	aburos "github.com/Autonomous-Systems-Laboratory-UNIUD/aburos"
	abuagent "github.com/Autonomous-Systems-Laboratory-UNIUD/aburos/agent"
	"github.com/Autonomous-Systems-Laboratory-UNIUD/abusim-goabu-agent/endpoint"
	"github.com/Autonomous-Systems-Laboratory-UNIUD/abusim-goabu-agent/memory"

	"github.com/Autonomous-Systems-Laboratory-UNIUD/abusim-core/schema"
	rosetta "github.com/Autonomous-Systems-Laboratory-UNIUD/goROSetta/ROSetta"

	"log"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "PANIC: %v\n", r)
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	// log every fatal-capable call site, or just add this:
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	arduType := []string{"copter", "plane", "sub", "rover"}
	// I check if a config is present on the Args...
	if len(os.Args) < 2 {
		log.Fatalln("Config not found, exiting")
	}
	// ... and I deserialize it to get its fields
	configStr := os.Args[1]
	agent := schema.AgentConfiguration{}
	err := agent.Deserialize(configStr)
	if err != nil {
		log.Fatalf("Bad config deserialization: %v", err)
	}
	// I create the memory for the agent...
	log.Println("Creating memory")
	mem, err := memory.New(agent.MemoryController, agent.Memory)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("mem: %v\n", mem)
	// ... I create the executer...
	log.Println("Creating agent")
	//logConfig := goabuconfig.LogConfig{
	//	Encoding: "console",
	//	Level:    goabuconfig.LogError,
	//}
	abuAgent, err := abuagent.NewRosAgent()
	if err != nil {
		log.Fatalln(err)
	}
	abuAgent.Zenoh = true
	log.Println("Creating executer")
	//exec := aburos.RosExecuter{}
	exec, err := aburos.NewRosExecuter(mem, agent.Rules, abuAgent, agent.Name, "aburos", "lazy")
	if err != nil {
		log.Fatal(err)
	}
	var rosettaNode *rosetta.ROSettaNode
	bridgeOk := false
	if slices.Contains(arduType, agent.MemoryController) {
		tries := 4
		log.Println("Creating rosetta node")
		for try := range tries {
			rosettaNode, err = rosetta.NewROSettaNode(agent.Name, agent.SimAddr, strconv.Itoa(agent.SimPort), agent.SimID, 100, nil)
			if err != nil {
				if try != tries-1 {
					log.Println(err.Error() + fmt.Sprintf(", for agent %s, with address %s:%s", agent.Name, agent.SimAddr, strconv.Itoa(agent.SimPort)))
					time.Sleep(time.Duration(try) * 100 * time.Millisecond)
				}
			} else {
				bridgeOk = true
				log.Println("Created rosetta node")
				break
			}
		}

	}
	if bridgeOk {
		defer rosettaNode.Close()
	}

	// ... and I create the paused variable
	paused := true
	// I connect to the coordinator...
	time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
	log.Println("Connecting to coordinator")
	end, err := endpoint.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer end.Close()
	end.BridgeOk = bridgeOk

	// ... I send to it the initialization message...
	err = end.SendInit(agent.Name)
	if err != nil {
		log.Fatalln(err)
	}
	// ... and I start the main message loop
	go end.HandleMessages(exec, agent, &paused)
	// Finally, I start the executer loop
	log.Println("Starting main loop")
	for {
		// I execute a command if not paused...
		if !paused {
			exec.Exec()
		}
		// ... and I sleep for a while
		time.Sleep(agent.Tick)
	}
}
