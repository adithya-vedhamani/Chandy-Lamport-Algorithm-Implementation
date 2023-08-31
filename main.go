package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Process struct {
	ID           int
	Accounts     []int
	Channels     []chan Message
	Snapshot     []int
	State        []int
	mutex        sync.Mutex
	wg           sync.WaitGroup
	markerWg     sync.WaitGroup
	transactions int
}

type Message struct {
	Amount int
	Sender int
	Marker bool
}

var numProcesses int
var processes []Process
var totalAmount int
var totalMutex sync.Mutex

const TransactionsPerSnapshot = 15

func generateTransaction(processID, numProcesses int) (receiverID, amount int) {
	receiverID = rand.Intn(numProcesses)
	amount = rand.Intn(100)
	return receiverID, amount
}

func (p *Process) transaction() {
	defer p.wg.Done()
	for i := 0; i < TransactionsPerSnapshot; i++ {
		time.Sleep(time.Millisecond * 500)
		receiverID, amount := generateTransaction(p.ID, numProcesses)

		if receiverID != p.ID {
			p.mutex.Lock()
			p.Accounts[receiverID] += amount
			p.Accounts[p.ID] -= amount
			p.mutex.Unlock()

			fmt.Printf("Process %d sent $%d to Process %d\n", p.ID, amount, receiverID)
			totalMutex.Lock()
			totalAmount += amount
			totalMutex.Unlock()

			p.transactions++

			if p.transactions%TransactionsPerSnapshot == 0 {
				p.initiateSnapshot()
			}
		}
	}
}

func (p *Process) initiateSnapshot() {
	p.mutex.Lock()
	p.Snapshot = make([]int, numProcesses)
	copy(p.Snapshot, p.Accounts)
	for i := 0; i < numProcesses; i++ {
		if i != p.ID {
			p.Channels[i] <- Message{Sender: p.ID, Marker: true}
		}
	}
	p.mutex.Unlock()
}

func (p *Process) handleMarkerMessages() {
	defer p.markerWg.Done()
	markersReceived := make([]bool, numProcesses)
	for {
		msg, ok := <-p.Channels[p.ID]
		if !ok {
			break
		}
		if msg.Marker {
			p.mutex.Lock()
			if !markersReceived[msg.Sender] {
				p.State[msg.Sender] = p.Accounts[msg.Sender]
				markersReceived[msg.Sender] = true
				p.handleSnapshotMessages()
			}
			p.mutex.Unlock()
		}
	}
}

func (p *Process) handleSnapshotMessages() {
	for i := 0; i < numProcesses; i++ {
		if i != p.ID {
			p.Channels[i] <- Message{Sender: p.ID, Amount: p.Snapshot[p.ID]}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Print("Enter the number of processes: ")
	_, err := fmt.Scan(&numProcesses)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	processes = make([]Process, numProcesses)
	for i := 0; i < numProcesses; i++ {
		processes[i] = Process{
			ID:       i,
			Accounts: make([]int, numProcesses),
			Channels: make([]chan Message, numProcesses),
			Snapshot: make([]int, numProcesses),
			State:    make([]int, numProcesses),
		}
		for j := 0; j < numProcesses; j++ {
			processes[i].Channels[j] = make(chan Message, 100)
		}
		for j := 0; j < numProcesses; j++ {
			processes[i].Accounts[j] = rand.Intn(1000) + 900
		}
		fmt.Printf("Process %d initial amount: $%d\n", i, sum(processes[i].Accounts))
		totalMutex.Lock()
		totalAmount += sum(processes[i].Accounts)
		totalMutex.Unlock()
		processes[i].wg.Add(1)
		processes[i].markerWg.Add(1)
		go processes[i].transaction()
		go processes[i].handleMarkerMessages()
	}
	fmt.Printf("Total amount in the whole system: $%d\n", totalAmount)
	for i := 0; i < numProcesses; i++ {
		processes[i].wg.Wait()
	}
	for i := 0; i < numProcesses; i++ {
		for j := 0; j < numProcesses; j++ {
			close(processes[i].Channels[j])
		}
	}
	for i := 0; i < numProcesses; i++ {
		processes[i].markerWg.Wait()
	}
	fmt.Printf("\nSnapshots and Consistency Check:\n")
	for i, p := range processes {
		fmt.Printf("Snapshot for Process %d: %v\n", i, p.Snapshot)
		isConsistent := true
		for j, snapshotAmount := range p.Snapshot {
			if snapshotAmount != p.State[j] {
				isConsistent = false
				break
			}
		}
		fmt.Printf("Is Snapshot Consistent for Process %d: %v\n", i, isConsistent)
		snapshotTotal := sum(p.Snapshot)
		fmt.Printf("Total amount in snapshot for Process %d: $%d\n", i, snapshotTotal)
		if snapshotTotal == totalAmount {
			fmt.Printf("Snapshot total matches initial total for Process %d\n", i)
		} else {
			fmt.Printf("Snapshot total does not match initial total for Process %d\n", i)
		}
	}
}

func sum(nums []int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}
