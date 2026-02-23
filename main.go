package main

import (
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/Autonomous-Systems-Laboratory-UNIUD/abusim-goabu-agent/endpoint"
	"github.com/Autonomous-Systems-Laboratory-UNIUD/abusim-goabu-agent/memory"

	aburos "github.com/Autonomous-Systems-Laboratory-UNIUD/aburos"
	"github.com/Autonomous-Systems-Laboratory-UNIUD/abusim-core/schema"
	rosetta "github.com/Autonomous-Systems-Laboratory-UNIUD/gorosetta"

	"log"
)

func main() {
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
	// ... I create the executer...
	log.Println("Creating executer")
	//logConfig := goabuconfig.LogConfig{
	//	Encoding: "console",
	//	Level:    goabuconfig.LogError,
	//}
	abuAgent, err := aburos.NewRosAgent()
	if err != nil {
		log.Fatalln(err)
	}
	exec, err := aburos.NewRosExecuter(mem, agent.Rules, abuAgent, agent.Name, "aburos", "lazy")
	if err != nil {
		log.Fatal(err)
	}
	var rosettaNode *rosetta.ROSettaNode
	if slices.Contains(arduType, agent.MemoryController) {
		rosettaNode, err = rosetta.NewROSettaNode(agent.Name, agent.SimAddr, strconv.Itoa(agent.SimPort), nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer rosettaNode.Close()

	// ... and I create the paused variable
	paused := true
	// I connect to the coordinator...
	log.Println("Connecting to coordinator")
	end, err := endpoint.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer end.Close()

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
