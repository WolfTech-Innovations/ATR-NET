package main
import("bytes";"crypto/aes";"crypto/cipher";"crypto/rand";"crypto/sha256";"encoding/base64";"encoding/json";"fmt";"io";"math/big";"net";"net/http";"sync";"time")
import "os/exec"
func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Command %s %v failed: %v\nOutput:\n%s", name, args, err, output)
	}
	fmt.Printf("Output of %s %v:\n%s\n", name, args, output)
}

func installTools() {
	fmt.Println("Updating package list...")
	runCommand("sudo", "apt-get", "update")

	fmt.Println("Installing I2P router...")
	runCommand("sudo", "apt-get", "install", "-y", "i2prouter")

	fmt.Println("Installing Tor...")
	runCommand("sudo", "apt-get", "install", "-y", "tor")
}

func runsetup() {
	installTools()
	fmt.Println("I2P and Tor installation completed successfully.")
}
type Logger struct {
    mu sync.Mutex
}

func (l *Logger) LogOperation(layer, operation, details string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    timestamp := time.Now().Format("15:04:05.000")
    fmt.Printf("[\033[1;36m%s\033[0m] [\033[1;33m%s\033[0m] [\033[1;32m%s\033[0m] %s\n",
        timestamp, layer, operation, details)
}
type BHttpjRequest struct{ID string `json:"id"`;Method string `json:"method"`;URL string `json:"url"`;Headers map[string]string `json:"headers"`;Body []byte `json:"body,omitempty"`;Timestamp int64 `json:"timestamp"`;AuthToken string `json:"auth_token"`}
type BHttpjResponse struct{ID string `json:"id"`;Status int `json:"status"`;Headers map[string]string `json:"headers"`;Body []byte `json:"body"`;Timestamp int64 `json:"timestamp"`;AuthToken string `json:"auth_token"`}
type BlockchainBlock struct{Index uint64 `json:"index"`;Timestamp int64 `json:"timestamp"`;Data []byte `json:"data"`;PreviousHash string `json:"previous_hash"`;Hash string `json:"hash"`;Nonce uint64 `json:"nonce"`}
type BHttpjBlockchain struct{Blocks []BlockchainBlock `json:"blocks"`;Difficulty int `json:"difficulty"`;mu sync.RWMutex}
func NewBlockchain()*BHttpjBlockchain{genesis:=BlockchainBlock{Index:0,Timestamp:time.Now().Unix(),Data:[]byte("genesis"),PreviousHash:"0",Hash:"0",Nonce:0};return&BHttpjBlockchain{Blocks:[]BlockchainBlock{genesis},Difficulty:4}}
func(bc*BHttpjBlockchain)CalculateHash(block*BlockchainBlock)string{h:=sha256.New();h.Write([]byte(fmt.Sprintf("%d%d%s%s%d",block.Index,block.Timestamp,string(block.Data),block.PreviousHash,block.Nonce)));return base64.StdEncoding.EncodeToString(h.Sum(nil))}
func(bc*BHttpjBlockchain)MineBlock(block BlockchainBlock)BlockchainBlock{target:="";for i:=0;i<bc.Difficulty;i++{target+="0"};for{block.Hash=bc.CalculateHash(&block);if len(block.Hash)>=bc.Difficulty&&block.Hash[:bc.Difficulty]==target{break};block.Nonce++};return block}
func(bc*BHttpjBlockchain)AddBlock(data[]byte){bc.mu.Lock();defer bc.mu.Unlock();prev:=&bc.Blocks[len(bc.Blocks)-1];newBlock:=BlockchainBlock{Index:prev.Index+1,Timestamp:time.Now().Unix(),Data:data,PreviousHash:prev.Hash,Hash:"",Nonce:0};newBlock=bc.MineBlock(newBlock);bc.Blocks=append(bc.Blocks,newBlock)}
func(bc*BHttpjBlockchain)EncryptChain()([]byte,error){bc.mu.RLock();defer bc.mu.RUnlock();chainData,err:=json.Marshal(bc);if err!=nil{return nil,err};h:=sha256.Sum256(chainData);block,err:=aes.NewCipher(h[:]);if err!=nil{return nil,err};gcm,err:=cipher.NewGCM(block);if err!=nil{return nil,err};nonce:=make([]byte,gcm.NonceSize());encrypted:=gcm.Seal(nonce,nonce,chainData,nil);return[]byte(base64.StdEncoding.EncodeToString(encrypted)),nil}
func DecryptChain(encryptedData[]byte)(*BHttpjBlockchain,error){decoded,err:=base64.StdEncoding.DecodeString(string(encryptedData));if err!=nil{return nil,err};h:=sha256.Sum256(decoded);block,err:=aes.NewCipher(h[:]);if err!=nil{return nil,err};gcm,err:=cipher.NewGCM(block);if err!=nil{return nil,err};nonceSize:=gcm.NonceSize();if len(decoded)<nonceSize{return nil,fmt.Errorf("invalid data")};nonce,ciphertext:=decoded[:nonceSize],decoded[nonceSize:];plaintext,err:=gcm.Open(nil,nonce,ciphertext,nil);if err!=nil{return nil,err};var bc BHttpjBlockchain;err=json.Unmarshal(plaintext,&bc);return&bc,err}
type BHttpjPacket struct{ReqType string `json:"req_type"`;Data []byte `json:"data"`;PeerID [20]byte `json:"peer_id"`;InfoHash [20]byte `json:"info_hash"`;Blockchain []byte `json:"blockchain,omitempty"`}
type TorConnection struct{conn net.Conn;key[]byte}
func NewTorConnection(addr string)(*TorConnection,error){conn,err:=net.Dial("tcp",addr);if err!=nil{return nil,err};key:=make([]byte,32);rand.Read(key);return&TorConnection{conn:conn,key:key},nil}
func(tc*TorConnection)Send(data[]byte)error{block,_:=aes.NewCipher(tc.key);gcm,_:=cipher.NewGCM(block);nonce:=make([]byte,gcm.NonceSize());encrypted:=gcm.Seal(nonce,nonce,data,nil);_,err:=tc.conn.Write(encrypted);return err}
type I2pConnection struct {
    conn     net.Conn
    logger   *Logger
    attempts int
}

func NewI2pConnection(addr string) (*I2pConnection, error) {
    ic := &I2pConnection{
        logger:   &Logger{},
        attempts: 0,
    }

    // Try to connect with retries
    for ic.attempts < 3 {
        ic.logger.LogOperation("I2P", "CONNECT", fmt.Sprintf("Attempt %d: Connecting to %s", ic.attempts+1, addr))
        
        conn, err := net.Dial("tcp", addr)
        if err == nil {
            ic.conn = conn
            ic.logger.LogOperation("I2P", "SUCCESS", "Connected to I2P network")
            return ic, nil
        }
        
        // If connection fails, try to start I2P
        if ic.attempts == 0 {
            ic.logger.LogOperation("I2P", "STARTUP", "Connection failed, starting I2P service...")

            cmd := exec.Command("i2prouter", "start", "--port", "8887")
            err = cmd.Run()
            if err != nil {
                ic.logger.LogOperation("I2P", "ERROR", fmt.Sprintf("Failed to start I2P: %v", err))
            } else {
                ic.logger.LogOperation("I2P", "STARTUP", "I2P service started, waiting 10s...")
                time.Sleep(10 * time.Second)
            }
        }
        
        ic.attempts++
        time.Sleep(2 * time.Second)
    }
    
    return nil, fmt.Errorf("failed to establish I2P connection after %d attempts", ic.attempts)
}

func (ic *I2pConnection) Tunnel(data []byte) ([]byte, error) {
    if ic.conn == nil {
        return nil, fmt.Errorf("no active I2P connection")
    }
    
    ic.logger.LogOperation("I2P", "TUNNEL", fmt.Sprintf("Tunneling %d bytes", len(data)))
    
    _, err := ic.conn.Write(data)
    if err != nil {
        ic.logger.LogOperation("I2P", "ERROR", fmt.Sprintf("Write failed: %v", err))
        return nil, err
    }
    
    buf := make([]byte, 8192)
    n, err := ic.conn.Read(buf)
    if err != nil {
        ic.logger.LogOperation("I2P", "ERROR", fmt.Sprintf("Read failed: %v", err))
        return nil, err
    }
    
    ic.logger.LogOperation("I2P", "SUCCESS", fmt.Sprintf("Tunneled response: %d bytes", n))
    return buf[:n], nil
}

func (ic *I2pConnection) Close() error {
    if ic.conn != nil {
        return ic.conn.Close()
    }
    return nil
}
type ObfsConnection struct{}
func NewObfsConnection()*ObfsConnection{return&ObfsConnection{}}
func(oc*ObfsConnection)Obfuscate(input[]byte)[]byte{result:=make([]byte,0,len(input)*3);for _,b:=range input{result=append(result,b^0xAA);randByte:=make([]byte,1);rand.Read(randByte);result=append(result,randByte[0]);result=append(result,(b+0x33)^0x55)};return result}
func(oc*ObfsConnection)Deobfuscate(input[]byte)[]byte{result:=make([]byte,0,len(input)/3);for i:=0;i<len(input);i+=3{if i+2<len(input){result=append(result,input[i]^0xAA)}};return result}
func(oc*ObfsConnection)ObfuscateBlockchain(blockchainData[]byte)[]byte{h:=sha256.Sum256(blockchainData);obfuscated:=oc.Obfuscate(blockchainData);obfuscated=append(obfuscated,h[:]...);return[]byte(base64.StdEncoding.EncodeToString(obfuscated))}
func(oc*ObfsConnection)DeobfuscateBlockchain(obfuscatedData[]byte)([]byte,error){decoded,err:=base64.StdEncoding.DecodeString(string(obfuscatedData));if err!=nil{return nil,err};if len(decoded)<32{return nil,fmt.Errorf("invalid data")};data,hash:=decoded[:len(decoded)-32],decoded[len(decoded)-32:];deobfuscated:=oc.Deobfuscate(data);expectedHash:=sha256.Sum256(deobfuscated);if!bytes.Equal(expectedHash[:],hash){return nil,fmt.Errorf("integrity check failed")};return deobfuscated,nil}
type SnowflakeConnection struct{peers[]string}
func NewSnowflakeConnection()*SnowflakeConnection{return&SnowflakeConnection{peers:[]string{"127.0.0.1:8888"}}}
func(sc*SnowflakeConnection)Relay(data[]byte)([]byte,error){conn,err:=net.Dial("tcp",sc.peers[0]);if err!=nil{return nil,err};defer conn.Close();_,err=conn.Write(data);if err!=nil{return nil,err};buf:=make([]byte,8192);n,err:=conn.Read(buf);return buf[:n],err}
type BitTorrentNode struct{peerID[20]byte;peers map[[20]byte]string;requests map[string]*BHttpjRequest;responses map[string]*BHttpjResponse;sessions map[string]string;blockchain*BHttpjBlockchain;mu sync.RWMutex}
func NewBitTorrentNode()*BitTorrentNode{var peerID[20]byte;rand.Read(peerID[:]);return&BitTorrentNode{peerID:peerID,peers:make(map[[20]byte]string),requests:make(map[string]*BHttpjRequest),responses:make(map[string]*BHttpjResponse),sessions:make(map[string]string),blockchain:NewBlockchain()}}
func(btn*BitTorrentNode)GenerateAuthToken(url string)string{h:=sha256.Sum256([]byte(url));return base64.StdEncoding.EncodeToString(h[:])}
func(btn*BitTorrentNode)ValidateToken(url,token string)bool{return btn.GenerateAuthToken(url)==token}
func(btn*BitTorrentNode)AddToBlockchain(data[]byte){btn.blockchain.AddBlock(data)}
func(btn*BitTorrentNode)GetBlockchainData()([]byte,error){return btn.blockchain.EncryptChain()}
func(btn*BitTorrentNode)SendRequest(req*BHttpjRequest)error{jsonData,err:=json.Marshal(req);if err!=nil{return err};btn.AddToBlockchain(jsonData);blockchainData,err:=btn.GetBlockchainData();if err!=nil{return err};packet:=BHttpjPacket{ReqType:"request",Data:jsonData,PeerID:btn.peerID,InfoHash:[20]byte{},Blockchain:blockchainData};packetData,err:=json.Marshal(packet);if err!=nil{return err};btn.mu.Lock();btn.requests[req.ID]=req;btn.mu.Unlock();for _,addr:=range btn.peers{conn,err:=net.Dial("tcp",addr);if err!=nil{continue};conn.Write(packetData);conn.Close()};return nil}
func(btn*BitTorrentNode)HandlePacket(packetData[]byte)error{var packet BHttpjPacket;err:=json.Unmarshal(packetData,&packet);if err!=nil{return err};if len(packet.Blockchain)>0{obfs:=NewObfsConnection();deobfuscatedChain,err:=obfs.DeobfuscateBlockchain(packet.Blockchain);if err==nil{_,_=DecryptChain(deobfuscatedChain)}};switch packet.ReqType{case"request":var req BHttpjRequest;err:=json.Unmarshal(packet.Data,&req);if err!=nil{return err};if!btn.ValidateToken(req.URL,req.AuthToken){return nil};resp,err:=btn.ProcessHTTPRequest(&req);if err!=nil{return err};respData,_:=json.Marshal(resp);btn.AddToBlockchain(respData);blockchainData,_:=btn.GetBlockchainData();respPacket:=BHttpjPacket{ReqType:"response",Data:respData,PeerID:btn.peerID,InfoHash:packet.InfoHash,Blockchain:blockchainData};respPacketData,_:=json.Marshal(respPacket);for _,addr:=range btn.peers{conn,err:=net.Dial("tcp",addr);if err!=nil{continue};conn.Write(respPacketData);conn.Close();break};case"response":var resp BHttpjResponse;err:=json.Unmarshal(packet.Data,&resp);if err!=nil{return err};if!btn.ValidateToken("response",resp.AuthToken){return nil};btn.mu.Lock();btn.responses[resp.ID]=&resp;btn.mu.Unlock()};return nil}
func(btn*BitTorrentNode)ProcessHTTPRequest(req*BHttpjRequest)(*BHttpjResponse,error){client:=&http.Client{Timeout:30*time.Second};httpReq,err:=http.NewRequest(req.Method,req.URL,bytes.NewBuffer(req.Body));if err!=nil{return nil,err};for k,v:=range req.Headers{httpReq.Header.Set(k,v)};resp,err:=client.Do(httpReq);if err!=nil{return nil,err};defer resp.Body.Close();body,err:=io.ReadAll(resp.Body);if err!=nil{return nil,err};headers:=make(map[string]string);for k,v:=range resp.Header{if len(v)>0{headers[k]=v[0]}};token:=btn.GenerateAuthToken(req.URL);btn.mu.Lock();btn.sessions[req.URL]=token;btn.mu.Unlock();return&BHttpjResponse{ID:req.ID,Status:resp.StatusCode,Headers:headers,Body:body,Timestamp:time.Now().Unix(),AuthToken:token},nil}
type BHttpjProxy struct{btNode*BitTorrentNode;torConn*TorConnection;i2pConn*I2pConnection;obfs*ObfsConnection;snowflake*SnowflakeConnection}
func NewBHttpjProxy()*BHttpjProxy{return&BHttpjProxy{btNode:NewBitTorrentNode(),obfs:NewObfsConnection(),snowflake:NewSnowflakeConnection()}}
func (bp *BHttpjProxy) InitConnections() error {
    logger := &Logger{}
    
    // Initialize I2P first
    logger.LogOperation("PROXY", "INIT", "Initializing I2P connection...")
    i2pConn, err := NewI2pConnection("127.0.0.1:8887")
    if err != nil {
        logger.LogOperation("PROXY", "WARNING", fmt.Sprintf("I2P initialization failed: %v", err))
    } else {
        bp.i2pConn = i2pConn
    }
    
    // Initialize Tor
    logger.LogOperation("PROXY", "INIT", "Initializing Tor connection...")
    torConn, err := NewTorConnection("127.0.0.1:8888")
    if err != nil {
        logger.LogOperation("PROXY", "WARNING", fmt.Sprintf("Tor initialization failed: %v", err))
    } else {
        bp.torConn = torConn
    }
    
    return nil
}
func (bp *BHttpjProxy) ConvertToBHttpj(httpData string) (*BHttpjRequest, error) {
    lines := bytes.Split([]byte(httpData), []byte("\r\n"))
    if len(lines) == 0 {
        return nil, fmt.Errorf("invalid HTTP request: empty request")
    }

    firstLine := string(lines[0])
    parts := bytes.Fields([]byte(firstLine))
    if len(parts) < 2 {
        return nil, fmt.Errorf("invalid HTTP request: missing method or URL")
    }

    method := string(parts[0])
    url := string(parts[1])
    
    headers := make(map[string]string)
    bodyStart := 0
    
    // Parse headers safely
    for i := 1; i < len(lines); i++ {
        line := string(lines[i])
        if line == "" {
            bodyStart = i + 1
            break
        }
        
        headerParts := bytes.SplitN([]byte(line), []byte(":"), 2)
        if len(headerParts) == 2 {
            key := string(bytes.TrimSpace(headerParts[0]))
            value := string(bytes.TrimSpace(headerParts[1]))
            headers[key] = value
        }
    }

    var body []byte
    if bodyStart < len(lines) {
        body = bytes.Join(lines[bodyStart:], []byte("\n"))
    }

    // Generate request ID safely
    idInt, err := rand.Int(rand.Reader, big.NewInt(1000000))
    if err != nil {
        return nil, fmt.Errorf("failed to generate request ID: %v", err)
    }

    token := bp.btNode.GenerateAuthToken(url)

    return &BHttpjRequest{
        ID:        fmt.Sprintf("%d", idInt),
        Method:    method,
        URL:       url,
        Headers:   headers,
        Body:      body,
        Timestamp: time.Now().Unix(),
        AuthToken: token,
    }, nil
}
func(bp*BHttpjProxy)ProcessLayers(data[]byte)([]byte,error){snowflakeData,err:=bp.snowflake.Relay(data);if err!=nil{return nil,err};blockchainData,err:=bp.btNode.GetBlockchainData();if err!=nil{return nil,err};obfsBlockchain:=bp.obfs.ObfuscateBlockchain(blockchainData);obfsData:=bp.obfs.Obfuscate(snowflakeData);combinedData:=append(obfsData,obfsBlockchain...);i2pData:=combinedData;if bp.i2pConn!=nil{i2pData,err=bp.i2pConn.Tunnel(combinedData);if err!=nil{return nil,err}};if bp.torConn!=nil{err=bp.torConn.Send(i2pData);if err!=nil{return nil,err};return[]byte{},nil};return i2pData,nil}
func(bp*BHttpjProxy)HandleWebRequest(httpRequest string)(string,error){bhttpjReq,err:=bp.ConvertToBHttpj(httpRequest);if err!=nil{return"",err};reqJSON,err:=json.Marshal(bhttpjReq);if err!=nil{return"",err};_,err=bp.ProcessLayers(reqJSON);if err!=nil{return"",err};err=bp.btNode.SendRequest(bhttpjReq);if err!=nil{return"",err};time.Sleep(100*time.Millisecond);bp.btNode.mu.RLock();resp,exists:=bp.btNode.responses[bhttpjReq.ID];bp.btNode.mu.RUnlock();if exists{return fmt.Sprintf("HTTP/1.1 %d OK\r\nContent-Type: text/html\r\nContent-Length: %d\r\n\r\n%s",resp.Status,len(resp.Body),string(resp.Body)),nil};return"HTTP/1.1 408 Request Timeout\r\n\r\n",nil}
func (bp *BHttpjProxy) handleConnection(conn net.Conn) {
    defer conn.Close()
    
    buffer := make([]byte, 4096)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Printf("Error reading from connection: %v\n", err)
        return
    }
    
    if n == 0 {
        fmt.Printf("Empty request received\n")
        return
    }

    request := string(buffer[:n])
    response, err := bp.HandleWebRequest(request)
    
    if err != nil {
        errorResponse := fmt.Sprintf("HTTP/1.1 500 Internal Server Error\r\nContent-Type: text/plain\r\n\r\nError: %s", err)
        conn.Write([]byte(errorResponse))
        return
    }

    _, err = conn.Write([]byte(response))
    if err != nil {
        fmt.Printf("Error writing response: %v\n", err)
    }
}
func(bp*BHttpjProxy)StartServer(port int)error{listener,err:=net.Listen("tcp",fmt.Sprintf("127.0.0.1:%d",port));if err!=nil{return err};fmt.Printf("BHTTPJ Proxy listening on port %d\n",port);err=bp.InitConnections();if err!=nil{return err};go bp.startBitTorrentListener();for{conn,err:=listener.Accept();if err!=nil{continue};go bp.handleConnection(conn)}};
func(bp*BHttpjProxy)startBitTorrentListener(){listener,err:=net.Listen("tcp","127.0.0.1:6881");if err!=nil{return};for{conn,err:=listener.Accept();if err!=nil{continue};go func(c net.Conn){defer c.Close();buffer:=make([]byte,8192);n,err:=c.Read(buffer);if err!=nil{return};bp.btNode.HandlePacket(buffer[:n])}(conn)}}
func main() {
    // First run setup to install required tools
    runsetup()
    
    // Create new proxy instance
    proxy := NewBHttpjProxy()
    
    // Initialize logger
    logger := &Logger{}
    logger.LogOperation("MAIN", "STARTUP", "BHTTPJ proxy starting...")
    
    // Start the proxy server
    err := proxy.StartServer(8888)
    if err != nil {
        logger.LogOperation("MAIN", "ERROR", fmt.Sprintf("Failed to start: %v", err))
        return
    }
}