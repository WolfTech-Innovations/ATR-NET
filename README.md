![GitHub Repo stars](https://img.shields.io/github/stars/WolfTech-Innovations/BHTTPJ?style=social)
![GitHub last commit](https://img.shields.io/github/last-commit/WolfTech-Innovations/BHTTPJ)
# BHTTPJ - Bittorrent HTTP JSON

A privacy-focused proxy system that combines blockchain technology with multiple privacy layers including I2P, and Snowflake relays.

## Features

- 🔗 **Blockchain Integration** - Immutable request/response logging with PoW mining
- 🧅 **Multi-Layer Privacy**
  - I2P tunneling
  - Snowflake relay system
  - Custom traffic obfuscation
- 🌐 **P2P Communication** - BitTorrent-style peer discovery and data transfer
- 🔐 **Strong Encryption** - AES-GCM encryption with blockchain validation
- 🎯 **Request Authentication** - Token-based request verification

## Installation

```bash
# Clone the repository
git clone https://github.com/WolfTech-Innovations/BHTTPJ
cd BHTTPJ

# Build the proxy
go build -o bhttpj src/main.go

# Run setup (requires sudo for Tor/I2P installation)
sudo ./bhttpj
```

## Requirements

- Go 1.19+
- Linux system
- Superuser privileges (for I2P installation)

## Usage

1. Start the proxy:
```bash
sudo ./bhttpj
```

2. Configure your browser/application to use the proxy:
```
Proxy Address: 127.0.0.1
Port: 8888
```

3. Test the connection:
```bash
curl -v -X GET  http://192.168.0.11:8888/ 'https://example.com'
```

## Architecture

```
Client Request
     ↓
BHTTPJ Proxy
     ↓
Snowflake Relay
     ↓
Blockchain Verification
     ↓
Obfuscation Layer
     ↓
I2P Tunnel
     ↓
Target Server
```

## Security Features

- Request/response pairs stored in blockchain
- Multi-layer encryption
- Traffic obfuscation
- Distributed relay network
- Token-based authentication
- Proof-of-work validation

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License

MIT

## Disclaimer

This software is for educational and research purposes only. Users are responsible for ensuring compliance with local laws and regulations regarding network privacy tools.
