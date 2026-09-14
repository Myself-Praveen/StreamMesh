# 🌐 StreamMesh

**StreamMesh** is a distributed real-time event broker and WebSocket gateway designed to handle large-scale concurrent connections (50k+). It provides a seamless pub/sub system across multiple nodes using Redis Streams and gRPC.

[![Go Version](https://img.shields.io/badge/go-1.26-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

## Features
- **High-Concurrency WebSocket Gateway:** Built with Go to handle 50k+ active connections efficiently.
- **Distributed Pub/Sub:** Uses Redis Streams and gRPC to fan out messages across a cluster of StreamMesh nodes.
- **Token-Bucket Rate Limiting:** Prevents flooding and DDoS via atomic-level token bucket rate limiting.
- **Heartbeat & Zombie Pruning:** Automatically tracks active connections and prunes dead ones to save memory.
- **Developer Dashboard:** Real-time metrics dashboard built with Next.js, featuring Recharts for live throughput monitoring.

## Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        WS1["WebSocket Clients"]
        WS2["Dashboard UI"]
        API["REST API Clients"]
    end

    subgraph "Gateway Layer (Go)"
        LB["Load Balancer"]
        N1["StreamMesh Node 1"]
        N2["StreamMesh Node 2"]
    end

    subgraph "Node Internals"
        CM["Connection Manager"]
        RL["Rate Limiter"]
        PS["PubSub Engine"]
    end

    subgraph "Data Layer"
        RS["Redis Streams"]
    end

    WS1 --> LB --> N1 & N2
    N1 --> CM & RL & PS
    PS --> RS
    RS --> N2
    DASH["Live Dashboard"] --> N1
```

## Getting Started

1. Clone the repository.
2. Run `docker-compose up -d` to spin up the required Redis instances.
3. Start the server: `make run`.
4. Open the dashboard at `http://localhost:3000`.
