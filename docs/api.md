# ATR-NET API Documentation

## Overview

ATR-NET (Anonymous Traffic Routing Network) is a decentralized peer-to-peer network implementation that provides anonymous communication, distributed storage, and web proxying capabilities. The network uses onion-style routing, blockchain-based data integrity, and distributed hash tables for content discovery.

## Network Constants

```go
const (
    NETID    = "ATR-NET-V1"  // Network identifier
    MAXHOPS  = 9             // Maximum routing hops
    CHUNKS   = 7             // Data chunking factor
    PKTSIZE  = 65536         // Packet size (64KB)
    BLKSIZE  = 1048576       // Block size (1MB)
    NODEPORT = 7777          // Default node port
    BOOTPORT = 7778          // Bootstrap server port
    DNSPORT  = 7779          // DNS server port
    WEBPORT  = 7780          // Web server port
)
```

## Core Types

### XL (XOR Layer)
Stream cipher for data encryption/decryption.

```go
type XL struct {
    k []byte  // Key material
    p int     // Position counter
}
```

**Methods:**
- `NewXL() XL` - Creates new XL cipher with random 512-byte key
- `E(data []byte) []byte` - Encrypts data using XOR stream
- `D(data []byte) []byte` - Decrypts data (same as encrypt for XOR)

### H256 (Hash256)
32-byte SHA256 hash representation.

```go
type H256 [32]byte
```

**Methods:**
- `S() string` - Returns short string representation (first 8 bytes as hex)
- `NH(data []byte) H256` - Creates new hash from data

### NK (Node Key)
Ed25519 cryptographic key pair for node identity.

```go
type NK struct {
    pr ed25519.PrivateKey  // Private key
    pk ed25519.PublicKey   // Public key  
    h  H256                // Hash of public key (node ID)
}
```

**Methods:**
- `NewNK() (NK, error)` - Generates new key pair
- `S(data []byte) []byte` - Signs data
- `V(data, sig []byte) bool` - Verifies signature

### NAddr (Node Address)
Network address information for nodes.

```go
type NAddr struct {
    H    H256      // Node hash/ID
    IP   string    // IP address
    Port int       // Port number
    T    time.Time // Timestamp
}
```

**Methods:**
- `S() string` - Returns formatted string "hash@ip:port"
- `A() string` - Returns address string "ip:port"

## Blockchain Components

### BLK (Block)
Blockchain block containing data and metadata.

```go
type BLK struct {
    H H256    // Block hash
    D []byte  // Block data
    P H256    // Previous block hash
    T int64   // Timestamp
    N int64   // Nonce
    S []byte  // Signature
}
```

**Methods:**
- `NewBLK(data []byte, prev H256, nonce int64, nk NK) BLK` - Creates new signed block

### BC (Blockchain)
Thread-safe blockchain implementation.

```go
type BC struct {
    b  []BLK           // Block slice
    h  map[H256]BLK    // Hash to block mapping
    mu sync.RWMutex    // Read-write mutex
    nk NK              // Node key for signing
}
```

**Methods:**
- `NewBC(nk NK) BC` - Creates new blockchain
- `A(block BLK) bool` - Adds block (returns false if exists)
- `G(hash H256) BLK` - Gets block by hash
- `L() H256` - Returns hash of latest block

## DHT (Distributed Hash Table)

```go
type DHT struct {
    d  map[string]NAddr     // Domain to address mapping
    k  map[H256][]byte      // Hash to data mapping
    mu sync.RWMutex        // Read-write mutex
}
```

**Methods:**
- `NewDHT() DHT` - Creates new DHT
- `RD(name string, addr NAddr)` - Registers domain
- `FD(name string) (NAddr, bool)` - Finds domain
- `SK(hash H256, data []byte)` - Stores key-value pair
- `GK(hash H256) ([]byte, bool)` - Gets value by key

## Messaging System

### MSG (Message)
Network message structure.

```go
type MSG struct {
    T  string  // Message type
    F  H256    // From (sender hash)
    To H256    // To (recipient hash)
    D  []byte  // Data payload
    TS int64   // Timestamp
    S  []byte  // Signature
    H  H256    // Message hash
}
```

**Methods:**
- `NewMSG(type string, from, to H256, data []byte, nk NK) MSG` - Creates signed message
- `EN() []byte` - Encodes message to JSON

**Message Types:**
- `HELLO` - Initial connection request
- `HELLO_ACK` - Connection acknowledgment
- `PING` - Keep-alive ping
- `PONG` - Ping response
- `RESOLVE` - Domain name resolution request
- `RESOLVED` - Domain resolution response
- `PUBLISH` - Publish content/service
- `GET` - Request data by hash
- `PUT` - Store data
- `DATA` - Data response
- `GETPEERS` - Request peer list

### PKT (Packet)
Routing packet with onion-style layered encryption.

```go
type PKT struct {
    T string   // Packet type
    D []byte   // Data payload
    R []H256   // Route (node hashes)
    L int      // Remaining hops
    X XL       // XOR cipher instance
    E bool     // Encrypted flag
}
```

**Methods:**
- `NewPKT(type string, data []byte, route []H256) PKT` - Creates new packet
- `EN() []byte` - Encodes and encrypts packet
- `DEPKT(data []byte) (PKT, error)` - Decodes packet

## Node Implementation

### PEER
Peer connection information.

```go
type PEER struct {
    A  NAddr     // Peer address
    K  []byte    // Shared key
    S  int       // Score/reputation
    L  time.Time // Last seen
    BC BC        // Peer's blockchain state
}
```

### NODE
Core network node implementation.

```go
type NODE struct {
    nk      NK                 // Node key pair
    addr    NAddr             // Node address
    peers   map[H256]PEER     // Connected peers
    dht     DHT               // Distributed hash table
    bc      BC                // Blockchain
    msgs    chan MSG          // Message channel
    pkts    chan PKT          // Packet channel
    mu      sync.RWMutex      // Mutex
    xl      XL                // Encryption layer
    running bool              // Running status
}
```

**Methods:**

#### Network Operations
- `NewNODE(ip string, port int) (NODE, error)` - Creates new node
- `BOOT(bootAddrs []string) error` - Bootstraps to network
- `LISTEN() error` - Starts listening for connections
- `ADDPEER(addr NAddr, key []byte)` - Adds peer connection
- `GETPEERS(count int) []PEER` - Gets active peers

#### Routing & Communication
- `ROUTE(data []byte, hops int) ([]byte, error)` - Routes data through network
- `SEND(type string, to H256, data []byte) error` - Sends message
- `TRANSMIT(data []byte) error` - Transmits data to peers
- `PROC(data []byte) error` - Processes received data

#### Maintenance
- `MAINTAIN()` - Performs periodic maintenance
- `DISCOVER()` - Discovers new peers
- `HANDLEMSG(msg MSG)` - Handles incoming messages

## Network Services

### DNS Service
Decentralized domain name resolution.

```go
type DNS struct {
    node    NODE                // Underlying node
    domains map[string]string   // Domain mappings
}
```

**Methods:**
- `NewDNS(ip string) (DNS, error)` - Creates DNS server
- `REG(domain, target string)` - Registers domain
- `RES(domain string) string` - Resolves domain
- `START() error` - Starts DNS service

**Pre-registered Domains:**
- `search.atr` → `encrypted-search-engine.onion`
- `social.atr` → `decentralized-social.mesh`
- `code.atr` → `distributed-git.p2p`
- `news.atr` → `anonymous-news.net`
- `market.atr` → `private-marketplace.dark`

### WEB Server
HTTP server with proxy capabilities.

```go
type WEB struct {
    node      NODE                    // Underlying node
    pages     map[string][]byte       // Static pages
    templates map[string]string       // Page templates
}
```

**Methods:**
- `NewWEB(ip string) (WEB, error)` - Creates web server
- `SERVE(path string, content []byte)` - Serves static content
- `GET(path string) []byte` - Gets page content
- `START() error` - Starts web server
- `HTTP(conn net.Conn)` - Handles HTTP requests
- `PROXY(conn net.Conn, req string)` - Handles proxy requests

**Features:**
- Serves `.atr` domains through the mesh network
- Proxies `.clear` domains through regular internet
- Automatic routing through ATR-NET for anonymity

### BR (Browser Request)
HTTP request structure for proxy functionality.

```go
type BR struct {
    ID     string              // Request ID
    Method string              // HTTP method
    URL    string              // Target URL
    H      map[string]string   // Headers
    B      []byte              // Body
    M      map[string]string   // Metadata
    Time   int64               // Timestamp
}
```

### BOOT (Bootstrap Server)
Network bootstrap server for peer discovery.

```go
type BOOT struct {
    nodes []NAddr        // Known nodes
    mu    sync.RWMutex   // Mutex
}
```

**Methods:**
- `NewBOOT() BOOT` - Creates bootstrap server
- `ADD(addr NAddr)` - Adds node to list
- `LIST() []string` - Returns node addresses
- `START() error` - Starts bootstrap server
- `HANDLE(conn net.Conn)` - Handles bootstrap requests

## Complete Network

### NET
Main network orchestrator.

```go
type NET struct {
    boot  BOOT    // Bootstrap server
    dns   DNS     // DNS server
    web   WEB     // Web server
    nodes []NODE  // Network nodes
}
```

**Methods:**
- `NewNET() (NET, error)` - Creates complete network
- `ADDNODE(ip string, port int) error` - Adds node to network
- `START() error` - Starts all network services

## Usage Example

```go
// Create and start the network
net, err := NewNET()
if err != nil {
    log.Fatal(err)
}

// Start all services
err = net.START()
if err != nil {
    log.Fatal(err)
}

// Network is now running on:
// Bootstrap: :7778
// DNS: :7779  
// Web: :7781
// Nodes: :7777-7781
```

## Security Features

- **Ed25519 Signatures**: All messages and blocks are cryptographically signed
- **Onion Routing**: Multi-hop routing with layered encryption
- **XOR Stream Cipher**: Fast symmetric encryption for data streams
- **Hash-based Identity**: Nodes identified by public key hashes
- **Blockchain Integrity**: Data integrity through blockchain storage
- **Anonymous Proxying**: Internet access through distributed routing

## Network Topology

The network operates as a mesh where:
1. Nodes discover peers through bootstrap servers
2. Messages route through multiple hops for anonymity
3. Data is chunked and distributed across multiple paths
4. DHT provides decentralized content discovery
5. Blockchain ensures data integrity and persistence

## Port Configuration

- **7777-7781**: Node communication ports
- **7778**: Bootstrap server
- **7779**: DNS resolution service  
- **7780**: Internal web service
- **7781**: HTTP proxy interface

## HTTP API Reference (curl/wget)

### Bootstrap Server API (Port 7778)

#### Get Node List
```bash
# Using curl
curl -X GET localhost:7778 -d "LIST"

# Using wget
echo "LIST" | wget --post-data=- -O- localhost:7778

# Response format: comma-separated addresses
# Example: 127.0.0.1:7777,127.0.0.1:7778,127.0.0.1:7779
```

### Web Server API (Port 7781)

#### Access ATR-NET Sites (.atr domains)
```bash
# Access the main portal
curl http://localhost:7781/

# Access specific .atr domains
curl http://localhost:7781/search.atr
curl http://localhost:7781/social.atr
curl http://localhost:7781/code.atr
curl http://localhost:7781/news.atr
curl http://localhost:7781/market.atr

# Using wget
wget http://localhost:7781/
wget http://localhost:7781/search.atr
```

#### Proxy Regular Internet (.clear suffix)
```bash
# Access regular websites anonymously through ATR-NET
curl http://localhost:7781/google.com.clear
curl http://localhost:7781/github.com.clear
curl http://localhost:7781/stackoverflow.com.clear

# HTTPS sites
curl http://localhost:7781/https://api.github.com.clear

# Using wget
wget http://localhost:7781/example.com.clear
wget -O response.html http://localhost:7781/https://httpbin.org/get.clear
```

#### POST Requests Through Proxy
```bash
# POST data through proxy
curl -X POST http://localhost:7781/httpbin.org/post.clear \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello from ATR-NET"}'

# Form data
curl -X POST http://localhost:7781/httpbin.org/post.clear \
  -d "key1=value1&key2=value2"

# Using wget for POST
wget --post-data='{"test": "data"}' \
  --header="Content-Type: application/json" \
  http://localhost:7781/httpbin.org/post.clear
```

#### Custom Headers Through Proxy
```bash
# Send custom headers
curl http://localhost:7781/httpbin.org/headers.clear \
  -H "X-Custom-Header: ATR-NET-Client" \
  -H "User-Agent: ATR-NET-Browser/1.0"

# With authentication
curl http://localhost:7781/api.example.com/data.clear \
  -H "Authorization: Bearer your-token-here"
```

### DNS Resolution API (Port 7779)

The DNS service operates through the node messaging system. Direct HTTP access is not available, but you can test DNS resolution through the web interface:

```bash
# Test domain resolution (returns the resolved content)
curl http://localhost:7781/search.atr

# The DNS resolution happens automatically when accessing .atr domains
```

### Node Communication Examples

#### Raw TCP Communication (Advanced)
```bash
# Connect to a node directly (requires proper message formatting)
# This is typically handled by the ATR-NET protocol, not HTTP

# Example of sending a raw message to node
echo '{"T":"PING","F":"...","To":"...","D":"test","TS":...,"S":"...","H":"..."}' | nc localhost 7777
```

### API Response Formats

#### Bootstrap Server Response
```
# Success: comma-separated node addresses
127.0.0.1:7777,127.0.0.1:7778,127.0.0.1:7779

# Empty response if no nodes available
```

#### Web Server Responses
```html
<!-- ATR-NET Portal (/) -->
<!DOCTYPE html>
<html>
<head>
    <title>ATR-NET Portal</title>
    <style>body{background:#000;color:#0f0;font-family:monospace;padding:20px}h1{color:#f00}a{color:#0ff}</style>
</head>
<body>
    <h1>WELCOME TO ATR-NET</h1>
    <p>The Anonymous Traffic Routing Network</p>
    <p>Node ID: [node-hash]</p>
    <p>Peers: Connected to decentralized mesh</p>
    <p>Status: SECURE • ANONYMOUS • ENCRYPTED</p>
</body>
</html>
```

```html
<!-- 404 Response -->
404 NOT FOUND
```

#### Proxy Responses
Proxy responses maintain the original HTTP response format from the target server, including:
- Status codes (200, 404, 500, etc.)
- Original headers
- Response body

### Common Use Cases

#### 1. Check Network Status
```bash
# Check if ATR-NET is running
curl -f http://localhost:7781/ && echo "ATR-NET is online"

# Get bootstrap nodes
NODES=$(curl -s localhost:7778 -d "LIST")
echo "Available nodes: $NODES"
```

#### 2. Anonymous Web Browsing
```bash
# Browse anonymously
curl -s http://localhost:7781/ipinfo.io.clear | jq '.'

# Check your apparent location through ATR-NET
curl http://localhost:7781/httpbin.org/ip.clear
```

#### 3. Access Decentralized Services
```bash
# Access ATR-NET native services
curl http://localhost:7781/search.atr
curl http://localhost:7781/social.atr

# Check service availability
for service in search social code news market; do
  echo "Testing $service.atr..."
  curl -f http://localhost:7781/$service.atr > /dev/null && echo "✓ Online" || echo "✗ Offline"
done
```

#### 4. API Testing Through Proxy
```bash
# Test APIs anonymously
curl -H "Accept: application/json" \
  http://localhost:7781/api.github.com/users/octocat.clear

# POST to APIs
curl -X POST http://localhost:7781/httpbin.org/post.clear \
  -H "Content-Type: application/json" \
  -d '{"data": "anonymous request via ATR-NET"}'
```

### Error Handling

#### Common Error Responses
```bash
# Network unreachable
HTTP/1.1 502 Bad Gateway
Bad Gateway, Sorry :(

# Service unavailable  
HTTP/1.1 500 Internal Server Error
Proxy Error

# Invalid .atr domain
404 NOT FOUND
```

#### Connection Testing
```bash
# Test connectivity to each service
services=("7778:Bootstrap" "7779:DNS" "7781:Web")
for service in "${services[@]}"; do
  port=$(echo $service | cut -d: -f1)
  name=$(echo $service | cut -d: -f2)
  nc -z localhost $port && echo "$name: ✓" || echo "$name: ✗"
done
```

### Publishing Content (PUBLISH Message)

The PUBLISH message allows you to register domains and content on the ATR-NET. This happens through the node messaging system.

#### Publishing API Structure

The PUBLISH message format in the ATR-NET protocol:
```json
{
  "T": "PUBLISH",
  "F": "sender_node_hash",
  "To": "target_node_hash", 
  "D": "domain:target_address",
  "TS": timestamp,
  "S": "signature",
  "H": "message_hash"
}
```

#### Direct TCP Publishing (Advanced)

Since PUBLISH operates at the node messaging level, you need to send properly formatted messages to nodes:

```bash
# Example PUBLISH message structure (requires proper signing)
# Format: "domain:target_address"

# Publishing a service
PUBLISH_DATA="myservice.atr:192.168.1.100"

# Note: Direct TCP requires proper Ed25519 signing
# This is typically done through ATR-NET client software
```

#### Simulated Publishing via Web Interface

Since the web server automatically publishes content when you access it, you can trigger publishing by:

```bash
# Method 1: Access a domain to trigger auto-registration
curl http://localhost:7781/mynewsite.atr

# Method 2: Use the proxy to publish external content
curl http://localhost:7781/myexternalsite.com.clear
```

#### Publishing Through Node Integration

For proper publishing, you would typically:

1. **Create a publishing client:**
```go
// Example client code (Go)
node, err := NewNODE("127.0.0.1", 8888)
if err != nil {
    log.Fatal(err)
}

// Connect to ATR-NET
err = node.BOOT([]string{"127.0.0.1:7778"})
if err != nil {
    log.Fatal(err)
}

// Publish your service
domain := "myapp.atr"
target := "192.168.1.100:8080"
publishData := fmt.Sprintf("%s:%s", domain, target)

err = node.SEND("PUBLISH", H256{}, []byte(publishData))
if err != nil {
    log.Fatal(err)
}
```

2. **Publishing workflow:**
   - Your service starts on a specific IP:PORT
   - Send PUBLISH message with "domain:ip:port" format
   - ATR-NET nodes propagate the registration
   - Domain becomes accessible via .atr resolution

#### Content Publishing Examples

```bash
# Publishing different types of services:

# Web service
# Data: "webapp.atr:192.168.1.10:8080"

# API service  
# Data: "api.atr:10.0.0.5:3000"

# File server
# Data: "files.atr:192.168.1.20:8000"

# Database service
# Data: "db.atr:10.0.0.10:5432"
```

#### Publishing via Raw Socket (Advanced)

```bash
# Create a raw TCP connection to publish
# Note: This requires proper message signing and formatting

cat << 'EOF' > publish_message.json
{
  "T": "PUBLISH",
  "F": "your_node_hash_here",
  "To": "",
  "D": "mydomain.atr:192.168.1.100:8080",
  "TS": $(date +%s),
  "S": "signature_here",
  "H": "message_hash_here"
}
EOF

# Send to node (requires proper cryptographic signing)
nc localhost 7777 < publish_message.json
```

#### Automated Publishing Script

```bash
#!/bin/bash
# publish_service.sh - Automated service publishing

DOMAIN="$1"
TARGET="$2"
NODE_PORT="${3:-7777}"

if [ -z "$DOMAIN" ] || [ -z "$TARGET" ]; then
    echo "Usage: $0 <domain.atr> <target_ip:port> [node_port]"
    echo "Example: $0 myapp.atr 192.168.1.100:8080"
    exit 1
fi

echo "Publishing $DOMAIN -> $TARGET"

# Method 1: Trigger through web access
echo "Triggering publication via web interface..."
curl -s "http://localhost:7781/$DOMAIN" > /dev/null

# Method 2: Direct node communication (requires proper client)
echo "Service published! Access via: http://localhost:7781/$DOMAIN"

# Verify publication
echo "Testing resolution..."
sleep 2
curl -f "http://localhost:7781/$DOMAIN" && echo "✓ Successfully published" || echo "✗ Publication failed"
```

#### Publishing Verification

```bash
# Check if your domain is published
curl -f http://localhost:7781/yourdomain.atr && echo "Published successfully"

# Test domain resolution through the network
for node in 7777 7778 7779; do
    echo "Testing domain resolution on node $node..."
    # Domain resolution test would go here
done
```

#### Common Publishing Patterns

1. **Web Application:**
```bash
# Start your web app on port 8080
python3 -m http.server 8080 &

# Publish to ATR-NET
./publish_service.sh mywebapp.atr 127.0.0.1:8080

# Access via: http://localhost:7781/mywebapp.atr
```

2. **API Service:**
```bash
# Start API server
node server.js & # Runs on port 3000

# Publish API
./publish_service.sh myapi.atr 127.0.0.1:3000

# Test API through ATR-NET
curl http://localhost:7781/myapi.atr/endpoint
```

3. **File Sharing:**
```bash
# Start file server
cd /path/to/files
python3 -m http.server 9000 &

# Publish file service
./publish_service.sh files.atr 127.0.0.1:9000

# Access files
curl http://localhost:7781/files.atr/document.pdf
```

#### Publishing Limitations

- **Node Connectivity**: Your service must be reachable by ATR-NET nodes
- **Port Accessibility**: Published ports must be open and accessible
- **Network Propagation**: Domain registration takes time to propagate
- **Security**: Services are published openly on the mesh network
- **Persistence**: Publications may need renewal depending on node uptime

#### Publishing Security Notes

- **Domain Squatting**: First to publish a domain owns it during session
- **Service Validation**: No automatic validation of published services
- **Access Control**: Published services are accessible to all mesh users
- **Data Integrity**: Use HTTPS/TLS for sensitive services even within ATR-NET

### Performance Considerations

- **Proxy requests** go through multiple hops, expect 2-5x latency
- **ATR-NET domains** (.atr) have faster resolution as they're mesh-native
- **Large downloads** are chunked across multiple routes for anonymity
- **Connection pooling** is not supported; each request creates new connections
- **Published services** should handle increased latency from mesh routing