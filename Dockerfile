FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

# Update and install dependencies
RUN apt-get update && apt-get install -y \
    wget \
    curl \
    ca-certificates \
    sudo \
    systemd \
    golang-1.21 \
    software-properties-common \
    gnupg \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Use Go 1.21 from system package (or optionally install manually like before)
ENV PATH="/usr/lib/go-1.21/bin:${PATH}"
ENV GOPATH="/go"
ENV GOBIN="/go/bin"

# Create working directory
WORKDIR /app

# Copy Go module files first
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source
COPY src/ src/

# Build
RUN go build -o main src/main.go

# Create systemd service file
RUN mkdir -p /etc/systemd/system
RUN echo "[Unit]\n\
Description=Go App Service\n\
After=network.target\n\
\n\
[Service]\n\
ExecStart=/app/main\n\
WorkingDirectory=/app\n\
Restart=always\n\
User=root\n\
\n\
[Install]\n\
WantedBy=multi-user.target" > /etc/systemd/system/go-app.service

# Enable the service
RUN systemctl enable go-app.service

STOPSIGNAL SIGRTMIN+3

VOLUME [ "/sys/fs/cgroup" ]

# systemd as init
CMD ["/sbin/init"]
