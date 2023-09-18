# BIT-TORRENT
BitTorrent is a decentralized peer-to-peer (P2P) communication protocol and technology designed for distributing and sharing large files over the internet. In this project, I utilized the GoLang programming language to implement a range of functionalities, such as setting up both UDP and TCP connections, as well as managing data sharing and reception.

## How to run

1. Cloning The Project :
  ```
   git clone https://github.com/Ghanshyam144/Bit-Torrent-Client.git
  ```
2. Installing required golang packages : 
  ```
    go get .
  ```
3. Download the required torrent file in 'torrent' folder
4. Run program
  ```
 go run . "torrent file name"
 ```
   Eg :-  
    ``` go run . Sample.torrent ```

### Till Now
1. The bit-torrent client currently supports the download of a single file using the bit-torrent protocol.
2. The client currently only supports the use of UDP trackers for communication and coordination between peers.
3. The client uses a brute force algorithm for mapping pieces to connections between peers.
4. It supports download of multiple files in  a torrent simultaneously.
5. It also supports re-handshaking mechanism to automatically re-establish TCP connections that have failed after a successful handshake, thereby ensuring the client can continue to participate    in the swarm.
   
### Todo
 1. Enhance the client's functionality by incorporating additional tracker protocols such as HTTP, TCP, and WebSockets, in order to broaden the range of supported tracker types.
 2. Develop a robust upload feature that allows users to share data with other clients in the swarm, and prioritize uploads to maximize overall download speed.
 3. Periodically scan for and connect to new peers to maintain a healthy swarm with high availability and fast download speeds.
 4. Support for Distributed Hash-Table (DHT) and Network Address Translation (NAT) Traversal.
    
   


  
