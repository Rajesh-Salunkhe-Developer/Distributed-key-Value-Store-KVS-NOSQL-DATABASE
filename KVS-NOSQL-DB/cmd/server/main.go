package main

import (
	"bufio"
	"flag"
	"fmt"
	"kvs-nosql-db/internal/store"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type NodeState int
const (
	Follower NodeState = iota
	Leader
)

type Node struct {
	ID            string
	Port          string
	Peers         []string
	Ring          *store.HashRing
	KV            *store.KVStore
	State         NodeState
	CurrentLeader string
	mu            sync.Mutex
	lastHeartbeat time.Time
	logFileName   string 
}

func main() {
	id := flag.String("id", "node1", "Node ID")
	port := flag.String("port", "9001", "Port")
	peersList := flag.String("peers", "", "Peer ports")
	flag.Parse()

	// Seed for random election timeouts
	rand.Seed(time.Now().UnixNano())
	peers := strings.Split(*peersList, ",")
	
	// Inside main()
n := &Node{
    ID:            *id,
    Port:          *port,
    Peers:         peers,
    Ring:          store.NewHashRing(append([]string{*port}, peers...)),
    KV:            store.NewKVStore(),
    State:         Follower,
    lastHeartbeat: time.Now(),
    logFileName:   "/app/data/kvs_" + *id + ".log", // Updated path for Docker
}

	// STAGE 6: Load data from disk before starting
	n.recoverFromDisk()

	ln, _ := net.Listen("tcp", ":"+*port)
	fmt.Printf("✅ [%s] LIVE on %s | Persistence: %s\n", *id, *port, n.logFileName)

	go n.runElectionTimer()

	for {
		conn, _ := ln.Accept()
		go n.handleConnection(conn)
	}
}

// STAGE 6: Persistence Logic (WAL)
func (n *Node) recoverFromDisk() {
	file, err := os.Open(n.logFileName)
	if err != nil { return }
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 2 {
			n.KV.Put(parts[0], parts[1])
			count++
		}
	}
	if count > 0 {
		fmt.Printf("📦 [%s] Recovered %d keys from disk\n", n.ID, count)
	}
}

func (n *Node) writeToLog(key, value string) {
	// Open in append mode to preserve previous writes
	f, err := os.OpenFile(n.logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		f.WriteString(fmt.Sprintf("%s %s\n", key, value))
	}
}

// STAGE 5: Election with Quorum Logic
func (n *Node) runElectionTimer() {
	for {
		// Randomized timeout (2-4 seconds) to prevent collision
		timeout := time.Duration(2000+rand.Intn(2000)) * time.Millisecond
		time.Sleep(500 * time.Millisecond)

		n.mu.Lock()
		if n.State == Leader {
			n.sendHeartbeats()
		} else if time.Since(n.lastHeartbeat) > timeout {
			n.mu.Unlock() 
			n.attemptElection()
			continue
		}
		n.mu.Unlock()
	}
}

func (n *Node) attemptElection() {
	votes := 1 // Vote for self
	for _, peer := range n.Peers {
		resp := n.sendCommandToNode(peer, "VOTE_REQUEST "+n.Port)
		if resp == "GRANTED" { votes++ }
	}

	n.mu.Lock()
	// Quorum for 3 nodes is 2 votes.
	if votes >= 2 { 
		fmt.Printf("👑 [%s] Elected LEADER with %d votes\n", n.ID, votes)
		n.State = Leader
		n.CurrentLeader = n.Port
	}
	n.mu.Unlock()
}

func (n *Node) sendHeartbeats() {
	for _, peer := range n.Peers {
		go n.sendCommandToNode(peer, "HEARTBEAT "+n.Port)
	}
}

func (n *Node) handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		input := scanner.Text()
		parts := strings.Fields(input)
		if len(parts) < 1 { continue }
		cmd := strings.ToUpper(parts[0])

		switch cmd {
		case "HEARTBEAT":
			n.mu.Lock()
			n.lastHeartbeat = time.Now()
			n.State = Follower
			n.CurrentLeader = parts[1]
			n.mu.Unlock()

		case "VOTE_REQUEST":
			n.mu.Lock()
			// Only grant vote if we haven't seen a leader recently
			if time.Since(n.lastHeartbeat) > 1500*time.Millisecond && n.State != Leader {
				conn.Write([]byte("GRANTED\n"))
			} else {
				conn.Write([]byte("DENIED\n"))
			}
			n.mu.Unlock()

		case "SET":
              if len(parts) < 3 { continue }
               key, val := parts[1], parts[2]

              // 1. Persist and Store locally
                 n.writeToLog(key, val)
                    n.KV.Put(key, val)

               // 2. STAGE 7 FIX: Replicate to EVERYONE in the cluster
                   for _, peer := range n.Peers {
                // We run this in a goroutine so the user doesn't wait 
               // if one node is slow
                go n.sendCommandToNode(peer, fmt.Sprintf("REPLICATE %s %s", key, val))
    }
    
    conn.Write([]byte("OK\n"))

		case "GET":
			val, _ := n.KV.Get(parts[1])
			conn.Write([]byte(val + "\n"))

		case "REPLICATE":
			n.writeToLog(parts[1], parts[2])
			n.KV.Put(parts[1], parts[2])
		}
	}
}

func (n *Node) sendCommandToNode(address string, message string) string {
    // address is "node2:9002"
    conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
    if err != nil {
        return "OFFLINE"
    }
    defer conn.Close()
    fmt.Fprintf(conn, message+"\n")
    return "SENT" 
}