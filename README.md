# ATR-NET - The Anonymous Traffic Routing Network

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

> **Reimagining the Internet: A fully decentralized, anonymous, and encrypted peer-to-peer network**

ATR-NET is a revolutionary decentralized networking protocol that provides true anonymity, encryption, and censorship resistance. Built from the ground up in Go, it creates a parallel internet infrastructure that operates independently of traditional centralized systems.

## 🌟 Key Features

### 🔒 **Privacy & Security**
- **End-to-End Encryption**: All communications are encrypted using AES-GCM
- **Onion-Style Routing**: Multi-hop routing through up to 9 nodes for maximum anonymity
- **Ed25519 Signatures**: Cryptographic authentication for all network participants
- **Zero-Knowledge Architecture**: No central authority can monitor or control traffic

### 🌐 **Decentralized Infrastructure**
- **Peer-to-Peer Mesh Network**: Self-organizing network topology
- **Distributed Hash Table (DHT)**: Decentralized data storage and retrieval
- **Blockchain Integration**: Immutable ledger for network integrity
- **DNS-Free Resolution**: Built-in naming system for .atr domains

### 🚀 **Advanced Networking**
- **Chunked Data Transmission**: Splits data across multiple routes for enhanced security
- **Automatic Peer Discovery**: Dynamic network expansion and healing
- **Load Balancing**: Intelligent traffic distribution across network nodes
- **Clearnet Proxy**: Secure bridge to traditional internet

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Bootstrap     │    │   DNS Server    │    │   Web Server    │
│   Server        │    │   (.atr domains)│    │   (HTTP Proxy)  │
│   Port: 7778    │    │   Port: 7779    │    │   Port: 7781    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
    ┌────────────────────────────┼────────────────────────────┐
    │                            │                            │
┌─────────┐              ┌─────────┐              ┌─────────┐
│  Node   │◄────────────►│  Node   │◄────────────►│  Node   │
│ :7777   │              │ :7778   │              │ :7779   │
└─────────┘              └─────────┘              └─────────┘
    │                            │                            │
    └────────────────────────────┼────────────────────────────┘
                                 │
                        ┌─────────────────┐
                        │  Mesh Network   │
                        │  (Decentralized │
                        │   P2P Routing)  │
                        └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Network connectivity (for peer discovery)

### Installation

```bash
# Clone the repository
git clone https://github.com/WolfTech-Innovations/ATR-NET.git
cd ATR-NET

# Build the project
go build -o atr-net main.go

# Run ATR-NET
./atr-net
```

### Usage

Once started, ATR-NET will automatically:
1. Start the bootstrap server on port 7778
2. Launch DNS resolution service on port 7779
3. Initialize web proxy on port 7781
4. Create 5 mesh nodes on ports 7777-7781

Access the network through:
- **Web Portal**: `http://localhost:7781`

## 🌐 Network Services

### Clearnet Proxy
Access traditional websites through ATR-NET's anonymizing proxy:
```
http://localhost:7781/google.com.clear
http://localhost:7781/github.com.clear
```

## 🔧 Configuration

### Network Constants
```go
const (
    NETID     = "ATR-NET-V1"    // Network identifier
    MAXHOPS   = 9               // Maximum routing hops
    CHUNKS    = 7               // Data splitting factor
    NODEPORT  = 7777            // Base node port
    BOOTPORT  = 7778            // Bootstrap port
    DNSPORT   = 7779            // DNS service port
    WEBPORT   = 7780            // Web service port
)
```

### Security Parameters
- **Encryption**: AES-256-GCM with random keys
- **Signatures**: Ed25519 cryptographic signatures
- **Routing**: Multi-path onion routing
- **Key Exchange**: Ephemeral key generation

## 📡 Protocol Specification

### Message Types
- `HELLO` / `HELLO_ACK` - Peer discovery and handshake
- `PING` / `PONG` - Network health monitoring
- `RESOLVE` / `RESOLVED` - DNS resolution
- `PUBLISH` - Domain registration
- `GET` / `PUT` - Data storage operations

### Packet Structure
```go
type PKT struct {
    T    string   // Packet type
    D    []byte   // Data payload
    R    []H256   // Routing path
    L    int      // Hops remaining
    X    XL       // Encryption context
    E    bool     // Encrypted flag
}
```

## 🔍 Technical Deep Dive

### Cryptographic Primitives
- **Hash Function**: SHA-256 for content addressing
- **Digital Signatures**: Ed25519 for authentication
- **Symmetric Encryption**: AES-256-GCM for payload protection
- **Key Derivation**: Secure random key generation

### Network Topology
- **Mesh Structure**: Each node maintains connections to multiple peers
- **Dynamic Routing**: Adaptive path selection based on network conditions
- **Fault Tolerance**: Automatic route recovery and peer replacement
- **Scalability**: Logarithmic lookup complexity with DHT

### Data Persistence
- **Blockchain**: Immutable transaction ledger
- **DHT Storage**: Distributed content-addressable storage
- **Peer Caching**: Local storage for frequently accessed data

## 🛠️ Development

### Project Structure
```
ATR-NET/
├── main.go          # Core network implementation
├── README.md        # This file
├── LICENSE          # MIT License
|_________________
```

### Contributing
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Testing
```bash
# Run unit tests
go test ./...

# Network integration tests
go test -tags=integration ./...
```

## 🔒 Security Considerations

### Threat Model
- **Traffic Analysis**: Mitigated through multi-hop routing
- **Node Compromise**: Limited impact due to decentralized architecture
- **Censorship**: Resistant through distributed infrastructure
- **Surveillance**: Protected by end-to-end encryption

### Best Practices
- Regular key rotation
- Peer diversity maintenance
- Network monitoring
- Security audits

## 📄 License

This project is licensed under the MIT License

## 🙏 Acknowledgments

- The Tor Project for onion routing inspiration
- The Bitcoin community for blockchain concepts
- The Go team for excellent networking libraries
- All contributors and early adopters
- Claude AI for help with the coding and making this awesome README

---

**Built with ❤️ by WolfTech Innovations**

*"Decentralizing the future, one node at a time"*
