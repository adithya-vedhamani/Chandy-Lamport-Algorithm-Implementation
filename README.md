# Chandy-Lamport-Algorithm-Implementation

## Algorithm 


* Initialization:

Take input 'n' for the number of processes in the distributed system.
Create 'n' process objects, each representing a process in the system.
Initialize 3 bank accounts with some large amount for each process.
Create communication channels (Go channels) for every process to communicate with others.
Initialize space to store the local snapshot and state of the communication channel for each process.
Print the total amount in the whole system and store it as 'x'.

* Transaction Function:

Each process has a function called 'transaction()' that generates random credit/debit transactions between accounts of different processes.
Invoke the 'transaction()' function for all processes concurrently using Goroutines.

* Transaction Handling:

For every process 'pn', initiate the Chandy Lamport algorithm after every 'nth' transaction it initiates.
Upon initiating the 'nth' transaction, the process 'pn' records its local state (snapshot) by storing the current amounts in its accounts.
It then sends marker messages to all other processes through the communication channels.
Marker messages are used to identify the point at which the snapshot is taken.

* Marker Message Handling:

Each process listens for incoming messages on its communication channels.
When a marker message is received from another process, the receiving process records the state of the incoming channel (amount in the sender's accounts).
This recorded state represents the state of the communication channel at the point of receiving the marker message.

* Snapshot Consistency Check:

After the marker message handling is complete, each process compares its local snapshot with the recorded states of incoming channels.
If all recorded states match the corresponding amounts in the snapshot, the snapshot is considered consistent.
Otherwise, the snapshot is inconsistent.

* Print Snapshot Results:

Print the snapshot of each process, indicating the amounts in its accounts at the point of the snapshot.
Print whether the snapshot is consistent or not for each process.
Calculate the total amount in each snapshot and compare it with the initial total 'x'.
Print whether the snapshot total matches the initial total for each process.
