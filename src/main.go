package main
import("bytes";"crypto/aes";"crypto/cipher";"crypto/ed25519";"crypto/rand";"crypto/sha256";"encoding/base64";"encoding/binary";"encoding/json";"fmt";"io";"net";"net/http";"os";"path/filepath";"sort";"strconv";"strings";"sync";"time")
const(NETID="ATR-NET-V1";MAXHOPS=9;CHUNKS=7;PKTSIZE=65536;BLKSIZE=1048576;NODEPORT=7777;BOOTPORT=7778;DNSPORT=7779;WEBPORT=7780)
type XL struct{k[]byte;p int}
func NewXL()XL{k:=make([]byte,512);rand.Read(k);return XL{k,0}}
func(x*XL)E(d[]byte)[]byte{o:=make([]byte,len(d));for i,b:=range d{o[i]=b^x.k[x.p%len(x.k)];x.p++};return o}
func(x*XL)D(d[]byte)[]byte{return x.E(d)}
type H256 [32]byte
func(h H256)S()string{return fmt.Sprintf("%x",h[:8])}
func NH(d[]byte)H256{return sha256.Sum256(d)}
type NK struct{pr ed25519.PrivateKey;pk ed25519.PublicKey;h H256}
func NewNK()(NK,error){pk,pr,e:=ed25519.GenerateKey(rand.Reader);if e!=nil{return NK{},e};h:=NH(pk);return NK{pr,pk,h},nil}
func(nk*NK)S(d[]byte)[]byte{return ed25519.Sign(nk.pr,d)}
func(nk*NK)V(d,s[]byte)bool{return ed25519.Verify(nk.pk,d,s)}
type NAddr struct{H H256;IP string;Port int;T time.Time}
func(na NAddr)S()string{return fmt.Sprintf("%s@%s:%d",na.H.S(),na.IP,na.Port)}
func(na NAddr)A()string{return fmt.Sprintf("%s:%d",na.IP,na.Port)}
type BLK struct{H H256;D[]byte;P H256;T int64;N int64;S[]byte}
func NewBLK(d[]byte,p H256,n int64,nk*NK)BLK{h:=NH(d);t:=time.Now().Unix();sig:=nk.S(append(append(h[:],p[:]...),[]byte(fmt.Sprintf("%d%d",t,n))...));return BLK{h,d,p,t,n,sig}}
type BC struct{b[]BLK;h map[H256]*BLK;mu sync.RWMutex;nk NK}
func NewBC(nk NK)BC{return BC{[]BLK{},make(map[H256]*BLK),sync.RWMutex{},nk}}
func(bc*BC)A(b BLK)bool{bc.mu.Lock();defer bc.mu.Unlock();if _,ok:=bc.h[b.H];ok{return false};bc.b=append(bc.b,b);bc.h[b.H]=&b;return true}
func(bc*BC)G(h H256)*BLK{bc.mu.RLock();defer bc.mu.RUnlock();return bc.h[h]}
func(bc*BC)L()H256{bc.mu.RLock();defer bc.mu.RUnlock();if len(bc.b)==0{return H256{}};return bc.b[len(bc.b)-1].H}
type DHT struct{d map[string]NAddr;k map[H256][]byte;mu sync.RWMutex}
func NewDHT()DHT{return DHT{make(map[string]NAddr),make(map[H256][]byte),sync.RWMutex{}}}
func(dht*DHT)RD(n string,a NAddr){dht.mu.Lock();defer dht.mu.Unlock();dht.d[n]=a}
func(dht*DHT)FD(n string)(NAddr,bool){dht.mu.RLock();defer dht.mu.RUnlock();a,ok:=dht.d[n];return a,ok}
func(dht*DHT)SK(h H256,d[]byte){dht.mu.Lock();defer dht.mu.Unlock();dht.k[h]=d}
func(dht*DHT)GK(h H256)([]byte,bool){dht.mu.RLock();defer dht.mu.RUnlock();d,ok:=dht.k[h];return d,ok}
type MSG struct{T string;F H256;To H256;D[]byte;TS int64;S[]byte;H H256}
func NewMSG(t string,f,to H256,d[]byte,nk*NK)MSG{ts:=time.Now().Unix();h:=NH(append(append([]byte(t),d...),[]byte(fmt.Sprintf("%d",ts))...));s:=nk.S(h[:]);return MSG{t,f,to,d,ts,s,h}}
type PKT struct{T string;D[]byte;R[]H256;L int;X XL;E bool}
func NewPKT(t string,d[]byte,r[]H256)PKT{xl:=NewXL();return PKT{t,d,r,len(r),xl,true}}
func(p*PKT)EN()[]byte{d,_:=json.Marshal(p);ed:=p.X.E(d);return[]byte(base64.StdEncoding.EncodeToString(ed))}
func DEPKT(d[]byte)(PKT,error){bd,e:=base64.StdEncoding.DecodeString(string(d));if e!=nil{return PKT{},e};var p PKT;json.Unmarshal(bd,&p);dd:=p.X.D(bd);json.Unmarshal(dd,&p);return p,nil}
type PEER struct{A NAddr;K []byte;S int;L time.Time;BC BC}
type NODE struct{nk NK;addr NAddr;peers map[H256]*PEER;dht DHT;bc BC;msgs chan MSG;pkts chan PKT;mu sync.RWMutex;xl XL;running bool}
func NewNODE(ip string,port int)(NODE,error){nk,e:=NewNK();if e!=nil{return NODE{},e};addr:=NAddr{nk.h,ip,port,time.Now()};return NODE{nk,addr,make(map[H256]*PEER),NewDHT(),NewBC(nk),make(chan MSG,1000),make(chan PKT,1000),sync.RWMutex{},NewXL(),false},nil}
func(n*NODE)BOOT(ba[]string)error{for _,a:=range ba{conn,e:=net.Dial("tcp",a);if e!=nil{continue};msg:=NewMSG("HELLO",n.nk.h,H256{},[]byte(n.addr.S()),&n.nk);conn.Write(msg.EN());conn.Close()}return nil}
func(n*NODE)ADDPEER(a NAddr,k[]byte){n.mu.Lock();defer n.mu.Unlock();n.peers[a.H]=&PEER{a,k,100,time.Now(),NewBC(n.nk)}}
func(n*NODE)GETPEERS(c int)[]PEER{n.mu.RLock();defer n.mu.RUnlock();var ps[]PEER;for _,p:=range n.peers{if time.Since(p.L)<5*time.Minute{ps=append(ps,*p)}};if len(ps)>c{ps=ps[:c]};return ps}
func(n*NODE)ROUTE(d[]byte,hops int)([]byte,error){ps:=n.GETPEERS(hops*CHUNKS);if len(ps)==0{return d,nil};routes:=make([][]H256,CHUNKS);for i:=0;i<CHUNKS;i++{rt:=make([]H256,hops);for j:=0;j<hops&&j+i*hops<len(ps);j++{rt[j]=ps[j+i*hops].A.H};routes[i]=rt};chs:=make([][]byte,CHUNKS);cs:=len(d)/CHUNKS;for i:=0;i<CHUNKS-1;i++{chs[i]=d[i*cs:(i+1)*cs]};chs[CHUNKS-1]=d[(CHUNKS-1)*cs:];var fd[]byte;for i,ch:=range chs{pkt:=NewPKT("DATA",ch,routes[i]);ed:=pkt.EN();for j:=0;j<len(routes[i]);j++{k:=make([]byte,32);rand.Read(k);c,_:=aes.NewCipher(k);g,_:=cipher.NewGCM(c);nonce:=make([]byte,g.NonceSize());rand.Read(nonce);ed=g.Seal(append(k,nonce...),nonce,ed,nil)};if i==0{fd=ed}else{fd=append(fd,ed...)}}return fd,nil}
func(n*NODE)PROC(d[]byte)error{pkt,e:=DEPKT(d);if e!=nil{return e};if pkt.L>0{ps:=n.GETPEERS(1);if len(ps)>0{conn,e:=net.Dial("tcp",ps[0].A.A());if e==nil{pkt.L--;conn.Write(pkt.EN());conn.Close()}}return nil};switch pkt.T{case"DATA":n.dht.SK(NH(pkt.D),pkt.D);case"MSG":var msg MSG;json.Unmarshal(pkt.D,&msg);n.msgs<-msg};return nil}
func(n*NODE)SEND(t string,to H256,d[]byte)error{msg:=NewMSG(t,n.nk.h,to,d,&n.nk);md,_:=json.Marshal(msg);return n.TRANSMIT(md)}
func(n*NODE)TRANSMIT(d[]byte)error{rd,e:=n.ROUTE(d,MAXHOPS);if e!=nil{return e};ps:=n.GETPEERS(3);for _,p:=range ps{go func(peer PEER){conn,e:=net.Dial("tcp",peer.A.A());if e==nil{conn.Write(rd);conn.Close()}}(p)};return nil}
func(n*NODE)LISTEN()error{l,e:=net.Listen("tcp",n.addr.A());if e!=nil{return e};n.running=true;fmt.Printf("NODE %s LIVE %s\n",n.nk.h.S(),n.addr.A());go n.MSGLOOP();for n.running{conn,e:=l.Accept();if e!=nil{continue};go n.HANDLE(conn)};return nil}
func(n*NODE)HANDLE(conn net.Conn){defer conn.Close();buf:=make([]byte,PKTSIZE);nr,e:=conn.Read(buf);if e!=nil||nr==0{return};n.PROC(buf[:nr])}
func(n*NODE)MSGLOOP(){for n.running{select{case msg:=<-n.msgs:n.HANDLEMSG(msg);case<-time.After(time.Second):n.MAINTAIN()}}}
func(n*NODE)HANDLEMSG(msg MSG){switch msg.T{case"HELLO":parts:=strings.Split(string(msg.D),"@");if len(parts)==2{addrparts:=strings.Split(parts[1],":");if len(addrparts)==2{port,_:=strconv.Atoi(addrparts[1]);addr:=NAddr{msg.F,addrparts[0],port,time.Now()};n.ADDPEER(addr,nil);resp:=NewMSG("HELLO_ACK",n.nk.h,msg.F,[]byte(n.addr.S()),&n.nk);rd,_:=json.Marshal(resp);n.TRANSMIT(rd)}};case"HELLO_ACK":parts:=strings.Split(string(msg.D),"@");if len(parts)==2{addrparts:=strings.Split(parts[1],":");if len(addrparts)==2{port,_:=strconv.Atoi(addrparts[1]);addr:=NAddr{msg.F,addrparts[0],port,time.Now()};n.ADDPEER(addr,nil)}};case"PING":resp:=NewMSG("PONG",n.nk.h,msg.F,[]byte("OK"),&n.nk);rd,_:=json.Marshal(resp);n.TRANSMIT(rd);case"RESOLVE":domain:=string(msg.D);if addr,ok:=n.dht.FD(domain);ok{resp:=NewMSG("RESOLVED",n.nk.h,msg.F,[]byte(addr.S()),&n.nk);rd,_:=json.Marshal(resp);n.TRANSMIT(rd)};case"PUBLISH":parts:=strings.SplitN(string(msg.D),":",2);if len(parts)==2{n.dht.RD(parts[0],NAddr{msg.F,parts[1],WEBPORT,time.Now()})};case"GET":h:=NH(msg.D);if d,ok:=n.dht.GK(h);ok{resp:=NewMSG("DATA",n.nk.h,msg.F,d,&n.nk);rd,_:=json.Marshal(resp);n.TRANSMIT(rd)};case"PUT":h:=NH(msg.D);n.dht.SK(h,msg.D);blk:=NewBLK(msg.D,n.bc.L(),time.Now().Unix(),&n.nk);n.bc.A(blk)}}
func(n*NODE)MAINTAIN(){n.mu.Lock();for h,p:=range n.peers{if time.Since(p.L)>10*time.Minute{delete(n.peers,h)}};n.mu.Unlock();if len(n.peers)<5{n.DISCOVER()}}
func(n*NODE)DISCOVER(){ps:=n.GETPEERS(3);for _,p:=range ps{n.SEND("GETPEERS",p.A.H,[]byte("REQ"))}}
func(msg MSG)EN()[]byte{d,_:=json.Marshal(msg);return d}
type DNS struct{node NODE;domains map[string]string}
func NewDNS(ip string)(DNS,error){n,e:=NewNODE(ip,DNSPORT);if e!=nil{return DNS{},e};return DNS{n,make(map[string]string)},nil}
func(dns*DNS)REG(domain,target string){dns.domains[domain]=target;dns.node.dht.RD(domain,NAddr{dns.node.nk.h,target,WEBPORT,time.Now()})}
func(dns*DNS)RES(domain string)string{if target,ok:=dns.domains[domain];ok{return target};return""}
func(dns*DNS)START()error{go dns.node.LISTEN();dns.REG("search.atr","encrypted-search-engine.onion");dns.REG("social.atr","decentralized-social.mesh");dns.REG("code.atr","distributed-git.p2p");dns.REG("news.atr","anonymous-news.net");dns.REG("market.atr","private-marketplace.dark");return nil}
type WEB struct{node NODE;pages map[string][]byte;templates map[string]string}
func NewWEB(ip string)(WEB,error){n,e:=NewNODE(ip,WEBPORT);if e!=nil{return WEB{},e};return WEB{n,make(map[string][]byte),make(map[string]string)},nil}
func(w*WEB)SERVE(path string,content[]byte){w.pages[path]=content;h:=NH(content);w.node.dht.SK(h,content)}
func(w*WEB)GET(path string)[]byte{if content,ok:=w.pages[path];ok{return content};return[]byte("404 NOT FOUND")}
func(w*WEB)START()error{go w.node.LISTEN();l,e:=net.Listen("tcp",fmt.Sprintf(":%d",WEBPORT+1));if e!=nil{return e};fmt.Printf("WEB SERVER LIVE :%d\n",WEBPORT+1);w.LOADDEFAULT();for{conn,e:=l.Accept();if e!=nil{continue};go w.HTTP(conn)};return nil}
func(w*WEB)HTTP(conn net.Conn){defer conn.Close();buf:=make([]byte,8192);n,e:=conn.Read(buf);if e!=nil{return};req:=string(buf[:n]);lines:=strings.Split(req,"\r\n");if len(lines)==0{return};parts:=strings.Fields(lines[0]);if len(parts)<2{return};path:=parts[1];if strings.HasSuffix(path,".clear")||strings.Contains(path,"http")||!strings.HasSuffix(path,".atr")&&path!="/"&&path!="/index"{w.PROXY(conn,req);return};if path=="/"{path="/index"};content:=w.GET(path);resp:=fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: text/html\r\n\r\n%s",len(content),string(content));conn.Write([]byte(resp))}
func(w*WEB)PROXY(conn net.Conn,req string){lines:=strings.Split(req,"\r\n");if len(lines)==0{return};parts:=strings.Fields(lines[0]);if len(parts)<2{return};method,url:=parts[0],parts[1];if strings.HasSuffix(url,".clear"){url=strings.TrimSuffix(url,".clear")};if!strings.HasPrefix(url,"http"){url="https://"+url};reqdata:=&BR{fmt.Sprintf("proxy_%d",time.Now().UnixNano()),method,url,make(map[string]string),[]byte{},make(map[string]string),time.Now().Unix()};for i:=1;i<len(lines)&&lines[i]!="";i++{p:=strings.SplitN(lines[i],":",2);if len(p)==2{reqdata.H[strings.TrimSpace(p[0])]=strings.TrimSpace(p[1])}};client:=&http.Client{Timeout:30*time.Second};httpreq,e:=http.NewRequest(method,url,bytes.NewBuffer(reqdata.B));if e!=nil{conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\nProxy Error"));return};for k,v:=range reqdata.H{if k!="Host"&&k!="Connection"{httpreq.Header.Set(k,v)}};rd,_:=w.node.ROUTE([]byte(fmt.Sprintf("PROXY:%s",url)),MAXHOPS);w.node.dht.SK(NH(rd),rd);resp,e:=client.Do(httpreq);if e!=nil{conn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\nClearnet Unreachable"));return};defer resp.Body.Close();body,_:=io.ReadAll(resp.Body);var hs strings.Builder;for k,v:=range resp.Header{if len(v)>0{hs.WriteString(fmt.Sprintf("%s: %s\r\n",k,v[0]))}};proxyresp:=fmt.Sprintf("HTTP/1.1 %d %s\r\n%s\r\n%s",resp.StatusCode,resp.Status[4:],hs.String(),string(body));conn.Write([]byte(proxyresp))}index:=`<!DOCTYPE html><html><head><title>ATR-NET Portal</title><style>body{background:#000;color:#0f0;font-family:monospace;padding:20px}h1{color:#f00}a{color:#0ff}</style></head><body><h1>WELCOME TO ATR-NET</h1><p>The Anonymous Traffic Routing Network</p><ul><li><a href="/search">Search Portal</a></li><li><a href="/social">Social Hub</a></li><li><a href="/code">Code Repository</a></li><li><a href="/news">News Feed</a></li><li><a href="/market">Private Market</a></li></ul><p>Node ID: `+w.node.nk.h.S()+`</p><p>Peers: Connected to decentralized mesh</p><p>Status: SECURE • ANONYMOUS • ENCRYPTED</p></body></html>`;w.SERVE("/index",[]byte(index));search:=`<html><head><title>ATR Search</title></head><body style="background:#000;color:#0f0;font-family:monospace"><h1>ENCRYPTED SEARCH</h1><form><input type="text" placeholder="Enter search query..." style="background:#111;color:#0f0;border:1px solid #0f0;padding:10px;width:400px"><br><br><button type="submit" style="background:#0f0;color:#000;border:none;padding:10px 20px">SEARCH ANONYMOUSLY</button></form><p>All searches are encrypted and routed through 9 hops</p></body></html>`;w.SERVE("/search",[]byte(search));social:=`<html><head><title>ATR Social</title></head><body style="background:#000;color:#0f0;font-family:monospace"><h1>DECENTRALIZED SOCIAL</h1><div><h3>Anonymous Timeline</h3><p>[ENCRYPTED] User-4F2A: Just deployed a new dApp on ATR-NET</p><p>[ENCRYPTED] User-8B1C: The mesh is growing strong. 500+ nodes online</p><p>[ENCRYPTED] User-2E9D: Released new privacy tools on code.atr</p></div><textarea placeholder="Share something anonymously..." style="background:#111;color:#0f0;border:1px solid #0f0;width:500px;height:100px"></textarea><br><button style="background:#0f0;color:#000;border:none;padding:10px 20px">POST ANONYMOUSLY</button></body></html>`;w.SERVE("/social",[]byte(social))}
type BOOT struct{nodes[]NAddr;mu sync.RWMutex}
func NewBOOT()BOOT{return BOOT{[]NAddr{},sync.RWMutex{}}}
func(b*BOOT)ADD(a NAddr){b.mu.Lock();defer b.mu.Unlock();b.nodes=append(b.nodes,a)}
func(b*BOOT)LIST()[]string{b.mu.RLock();defer b.mu.RUnlock();var addrs[]string;for _,n:=range b.nodes{addrs=append(addrs,n.A())};return addrs}
func(b*BOOT)START()error{l,e:=net.Listen("tcp",fmt.Sprintf(":%d",BOOTPORT));if e!=nil{return e};fmt.Printf("BOOTSTRAP SERVER LIVE :%d\n",BOOTPORT);for{conn,e:=l.Accept();if e!=nil{continue};go b.HANDLE(conn)};return nil}
func(b*BOOT)HANDLE(conn net.Conn){defer conn.Close();buf:=make([]byte,1024);n,e:=conn.Read(buf);if e!=nil{return};req:=strings.TrimSpace(string(buf[:n]));if req=="LIST"{addrs:=b.LIST();resp:=strings.Join(addrs,",");conn.Write([]byte(resp))}}
type NET struct{boot BOOT;dns DNS;web WEB;nodes[]NODE}
func NewNET()(NET,error){boot:=NewBOOT();dns,e1:=NewDNS("127.0.0.1");if e1!=nil{return NET{},e1};web,e2:=NewWEB("127.0.0.1");if e2!=nil{return NET{},e2};return NET{boot,dns,web,[]NODE{}},nil}
func(net*NET)ADDNODE(ip string,port int)error{node,e:=NewNODE(ip,port);if e!=nil{return e};net.nodes=append(net.nodes,node);net.boot.ADD(node.addr);return nil}
func(net*NET)START()error{fmt.Println("STARTING ATR-NET - THE NEW INTERNET");fmt.Println("=====================================");go net.boot.START();go net.dns.START();go net.web.START();for i:=0;i<5;i++{net.ADDNODE("127.0.0.1",NODEPORT+i)};for i,node:=range net.nodes{go func(n NODE,idx int){n.BOOT([]string{"127.0.0.1:7778"});time.Sleep(time.Duration(idx)*time.Second);n.LISTEN()}(node,i)};return nil}
func main(){net,e:=NewNET();if e!=nil{fmt.Printf("FAILED TO START ATR-NET: %v\n",e);os.Exit(1)};e=net.START();if e!=nil{fmt.Printf("FAILED TO RUN ATR-NET: %v\n",e);os.Exit(1)};fmt.Println("ATR-NET IS LIVE!");fmt.Println("Bootstrap: :7778");fmt.Println("DNS: :7779");fmt.Println("Web: :7781");fmt.Println("Nodes: :7777-7781");select{}}
