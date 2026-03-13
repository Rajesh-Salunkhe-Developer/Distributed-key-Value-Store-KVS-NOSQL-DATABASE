# Distributed-key-Value-Store-KVS-NOSQL-DATABASE
# Distributed NoSQL Key-Value Store

A high-performance, containerized, distributed Key-Value database built from scratch in **Golang**. This project implements core distributed systems concepts including Leader Election, Replication, Persistence, and Self-Healing via Heartbeats.



##  Key Features

* **Custom TCP Protocol:** Optimized for low latency using raw TCP sockets instead of standard HTTP.
* **Leader-Follower Replication:** Fully automated data synchronization from the Leader node to all Followers.
* **Write-Ahead Logging (WAL):** Industrial-grade persistence. Every transaction is logged to disk before being written to memory, ensuring 0% data loss on crash recovery.
* **Heartbeat Mechanism:** A self-healing "Pulse" system where the Leader monitors peer health and dynamically manages the cluster state.
* **Consistent Hashing:** Built-in logic for balanced data sharding across multiple nodes.
* **High Concurrency:** Engineered with Go's `RWMutex` and `Goroutines` to handle thousands of simultaneous requests without race conditions.
* **Dockerized Orchestration:** Deploy a full 3-node cluster with persistent volumes using a single command.

## Tech Stack

* **Language:** Go (Golang)
* **Infrastructure:** Docker, Docker Compose
* **Networking:** TCP/IP Socket Programming
* **Persistence:** File-based Write-Ahead Logging (WAL)
* **Testing:** Custom Concurrent Benchmarking Tool

##  Getting Started

### Prerequisites
* Docker & Docker Compose
* Go 1.21+ (for running benchmarks locally)
