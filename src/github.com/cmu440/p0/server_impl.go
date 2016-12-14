// Implementation of a KeyValueServer. Students should write their code in this file.

package p0

import (
	"errors"
	"fmt"

	"net"

	"bufio"
	"bytes"
	"io"
	"strconv"
)

var PUT string = "put"
var GET string = "get"
var MAX_BUFFER_SIZE int = 500

type ClientProxy struct {
	conn             net.Conn    // 连接真正的client的conn
	messageChannel   chan []byte // 和server通信的channel
	exitReadChannel  chan int
	exitWriteChannel chan int
}

type keyValueServer struct {
	closed                bool
	debugMode             bool
	serverListener        *net.TCPListener
	clients               []*ClientProxy
	newConnectionsChannel chan net.Conn
	queryChannel          chan *Query
	exitClientsChannel    chan *ClientProxy
	countChannel          chan int
	exitMainChannel       chan int
	exitAcceptChannel     chan int
}

type Query struct {
	isGet bool
	key   string
	value []byte
}

// New creates and returns (but does not start) a new KeyValueServer.
func New() KeyValueServer {

	return &keyValueServer{false, false, nil, nil, make(chan net.Conn), make(chan *Query), make(chan *ClientProxy), make(chan int), make(chan int), make(chan int)}
}

func (kvs *keyValueServer) Start(port int) error {
	if kvs.closed == true {
		return errors.New("the server has been closed")
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))


	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		return err
	}
	init_db()
	kvs.serverListener = listener
	go kvs.acceptRequest()
	go kvs.serverMainRoutine()

	return nil
}

func (kvs *keyValueServer) Close() {
	if kvs.debugMode {
		fmt.Println("server closing")
	}
	kvs.closed = true
	kvs.serverListener.Close()
	kvs.exitMainChannel <- 0
	kvs.exitAcceptChannel <- 0
}

func (kvs *keyValueServer) Count() int {
	if kvs.debugMode {
		fmt.Println("get connect clients count")
	}
	if kvs.serverListener == nil || kvs.closed{
		return -1
	}
	kvs.countChannel <- 0

	count := <- kvs.countChannel;

	if kvs.debugMode {
		fmt.Println("connect clients count:", count)
	}
	return count
}

func (kvs *keyValueServer) serverMainRoutine() {
	for {
		select {
		case newConnection := <-kvs.newConnectionsChannel:
			newClient := &ClientProxy{newConnection, make(chan []byte, MAX_BUFFER_SIZE), make(chan int), make(chan int)}
			kvs.clients = append(kvs.clients, newClient)
			go clientReadRoutine(kvs, newClient)
			go clientWriteRoutine(newClient)
		case query := <-kvs.queryChannel:
			if query.isGet {
				keyPrefix := append([]byte(query.key), byte(','))
				response := append(keyPrefix, get(query.key)...)
				//response = append(response, byte('\n'))
				if kvs.debugMode {
					fmt.Println("get result:", response)
				}
				for _, client := range kvs.clients {
					if (len(client.messageChannel) == MAX_BUFFER_SIZE) {
						continue
					}
					client.messageChannel <- response
				}
			} else {
				put(query.key, query.value)
			}
		case exitClient := <-kvs.exitClientsChannel:
			for i, client := range kvs.clients {
				if client == exitClient {
					kvs.clients = append(kvs.clients[:i], kvs.clients[i+1:]...)
					break
				}
			}
		case <-kvs.countChannel:
			kvs.countChannel <- len(kvs.clients)
		case <-kvs.exitMainChannel:
			if kvs.debugMode {
				fmt.Println("server main routine exit")
			}
			for _, client := range kvs.clients {
				client.conn.Close()
				client.exitReadChannel <- 0
				client.exitWriteChannel <- 0
			}
			return
		}
	}
}

// Go routine used to accpet multiple connections
func (kvs *keyValueServer) acceptRequest() {

	for {
		select {
		case <-kvs.exitAcceptChannel:
			if kvs.debugMode {
				fmt.Println("accept routine exit")
			}
			return
		default:
			conn, err := kvs.serverListener.Accept()

			if err == nil {
				if kvs.debugMode {
					fmt.Println("connection accepted")
				}
				kvs.newConnectionsChannel <- conn

			} else {
				if kvs.debugMode {
					fmt.Println("listener failed to accept new connection", err)
				}
			}


		}

	}
}

func clientReadRoutine(kvs *keyValueServer, client *ClientProxy) {
	messageReader := bufio.NewReader(client.conn)
	for {
		select {
		case <-client.exitReadChannel:
			return
		default:
			res, err := messageReader.ReadBytes(byte('\n'))

			if err == io.EOF {
				if kvs.debugMode {
					fmt.Println("a channel exit initiative")
				}
				kvs.exitClientsChannel <- client
			}

			if err != nil {
				return
			}

			queryParam := bytes.Split(res, []byte(","))
			if kvs.debugMode {
				if string(queryParam[0]) == PUT {
					fmt.Println("%s:%v",string(queryParam[1]),queryParam[2])
				} else {
					fmt.Println("%s",string(queryParam[1]))
				}
			}
			switch string(queryParam[0]) {

			case PUT:
				kvs.queryChannel <- &Query{false, string(queryParam[1]), queryParam[2]}
			case GET:
				kvs.queryChannel <- &Query{isGet: true, key: string(queryParam[1][:len(queryParam[1]) - 1])}
			}
		}

	}
}

func clientWriteRoutine(client *ClientProxy) {
	for {
		select {
		case <-client.exitWriteChannel:
			return
		case message := <-client.messageChannel:
			client.conn.Write(message)
		}
	}
}
