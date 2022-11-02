package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

var MIN_TIMER int = 30
var MAX_TIMER int = 50
var NODES int = 3
var timeDuration int = 0
var term int = 0
var vote_acc int = 0
var isLeader bool = false

var timer *time.Timer
var self Node

type Node struct {
	id      int16
	address string // "IP:Port"
	port 		string
}

type Raft struct {

}

var members = []Node{
	{1, "127.0.0.1:9000", "9000"},
	{2, "127.0.0.1:9001", "9001"},
	{3, "127.0.0.1:9002", "9002"},
}

func (t *Raft) RequestVote(termCandidate int, result *int) error {
	if termCandidate > term {
		*result = 1
	} else {
		*result = 0
	}

	return nil
}

func (t *Raft) AppendEntries(termLeader int, result *int) error {
	term = termLeader
	printTime()
	fmt.Println("Menerima heartbeat dari leader.")
	resetTimer()
	restartCountdownTimer()

	return nil
}

func printTime() {
	dt := time.Now()
	fmt.Println(dt.Format("\n--> 15:04:05"))
}

func resetTimer() {
	if timeDuration > 0 {
		timer.Stop()
		fmt.Println("Timer " + self.port + " distop.")
	}
	timeDuration = rand.Intn(MAX_TIMER - MIN_TIMER) + MIN_TIMER
}

func restartCountdownTimer() {
	printTime()
	fmt.Println("Start timer with duration:", timeDuration)
	timer = time.AfterFunc(time.Duration(timeDuration)*time.Second, func() {
		
		go func() {
			for _, server := range members {
				client, err := rpc.DialHTTP("tcp", server.address)
				handleErr(err)
				term += 1
				var result int
				err = client.Call("Raft.RequestVote", term, &result)
				handleErr(err)
				if result == 1 {
					vote_acc += 1
				}
			}

			printTime()
			fmt.Println("Selesai request vote ke semua node.")
			fmt.Println("Vote acc yang didapat:", vote_acc)

			if vote_acc > NODES / 2 {
				isLeader = true

				fmt.Println("Sekarang", self.port, "adalah leader.")
			}

			go sendHeartBeat()
		}()
	})
}

func sendHeartBeat() {
	for isLeader {
		go func() {
			for _, server := range members {
				if server.id != self.id {
					fmt.Println("Mengirim heartbeat ke", server.port)
					client, err := rpc.DialHTTP("tcp", server.address)
					handleErr(err)
					var result int
					err = client.Call("Raft.AppendEntries", term, &result)
					handleErr(err)
				}
			}
		}()
		printTime()
		time.Sleep(time.Second * 10)
	}
}

func main () {
	rand.Seed(time.Now().UnixNano())

	portNumber := os.Args[1]

	for _, server := range members {
		if server.port == portNumber {
			self = server
		}
	}

	raft_ := &Raft{}
	rpc.Register(raft_)
	rpc.HandleHTTP()
	fmt.Println("Run server di port", portNumber)

	listener, err := net.Listen("tcp", ":"+string(portNumber))
	handleErr(err)

	resetTimer()
	go restartCountdownTimer()

	http.Serve(listener, nil)
}

func handleErr(err error) {
	if err != nil {
		fmt.Println("E:", err.Error())
	}
}