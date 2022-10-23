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

// Properti Timer
var heartbeat_counter int
var timer *time.Timer
var MIN_TIMER int = 100
var MAX_TIMER int = 500

var self Node
var isLeader = false

type Message struct {
	index int
	term  int
}

var logs = []Message{}

type SensorRPC struct {
}

type Node struct {
	id      int
	port    string
	address string
	leader  bool
}

var followers = []Node{
	{1, "9000", "127.0.0.1:9000", false},
	{2, "9001", "127.0.0.1:9001", false},
	{3, "9002", "127.0.0.1:9002", false},
	{4, "9003", "127.0.0.1:9003", false},
	{5, "9021", "127.0.0.1:9021", false},
}

func resetCountdownTimer() {
	timer.Stop()
	heartbeat_counter = rand.Intn((MAX_TIMER - MIN_TIMER)) + MIN_TIMER
}

func restartCountdownTimer() {
	fmt.Println("Start timer with duration : ", heartbeat_counter)
	timer = time.AfterFunc((heartbeat_counter)*time.Second, func() {
		// Switch state to Candidate and init election
		fmt.Println("Timer off")
	})
}

func (t *SensorRPC) AddSensorData(args *Message, reply *int) {
	// Reset countdown timer
	resetCountdownTimer()
	if args != nil {
		logs = append(logs, *args)
	}
	*reply = 200
	return nil
}

func (t *SensorRPC) RequestVote(args *Node, reply *int) {
	isLeader = false
	leaderID = *args.id
	for _, server := range members {
		if args == server {
			server.leader = true
		} else {
			server.leader = false
		}
	}
	resetCountdownTimer()
	*reply = 200
	return nil
}

func (t *SensorRPC) ClientRequest(args *Message, reply *int) {
	if isLeader {
		logs = append(logs, *args)
		for _, server := range member {
			if server != self {
				client, err := rpc.DialHTTP("tcp", server.address)
				handleError(err)
				var result int
				err = client.Call("ServerRPC.AppendEntries", args, &result)
				handleError(err)
			}
		}
	} else {
		for _, server := range members {
			if server.id == leaderID {
				client, err := rpc.DialHTTP("tcp", server.address)
				handleError(err)
				var result int
				err = client.Call("ServerRPC.ClientRequest", args, &result)
				handleError(err)
				break
			}
		}
	}
}

func heartbeatCountdown() {
	for heartbeat_counter > 0 {
		heartbeat_counter = heartbeat_counter - 1
		time.Sleep(1 * time.Millisecond)
	}

	if heartbeat_counter == 0 {
		var result int
		err = client.Call("SensorRPC.RequestVote", &self, &result)
		handleError(err)
		isLeader = true
		for _, server := range members {
			if self == server {
				server.leader = true
			} else {
				server.leader = false
			}
		}
		self.leader = true
	}
}

func sendHeartbeatToFollowers(args *Message) {
	logs = append(logs, *args)
	for isLeader {
		for _, server := range members {
			var result int
			fmt.Println("Mengirim AppendEntries RPC ke follower ", server.id)
			err = rpc.Call("SensorRPC.AppendEntries", args, &result)
			handleError(err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	// Inisasi RPC server
	serverID := os.Args[1]
	// Inisasi RPC server
	portNumber := os.Args[2]

	// Inisiasi struct
	service := &SensorRPC{}
	// Registrasikan struct dan method ke RPC
	rpc.Register(service)
	// Deklarasikan bahwa kita menggunakan protokol HTTP sebagai mekanisme pengiriman pesan
	rpc.HandleHTTP()
	fmt.Println("Run server ", serverID, " on port : ", portNumber)
	for _, server := range members {
		if server.port == port {
			self = server
		}
	}

	srv := &Raft{}
	rpc.Register(srv)
	rpc.HandleHTTP()

	// Deklarasikan listerner HTTP dengan layer transport TCP dan Port 1234
	listener, err := net.Listen("tcp", ":"+string(self.port))
	handleError(err)

	// Start timer
	resetCountdownTimer()
	go restartCountdownTimer()

	// Jalankan server HTTP
	go http.Serve(listener, nil)

}

func handleError(err error) {
	if err != nil {
		fmt.Println("Terdapat error : ", err.Error())
	}
}
