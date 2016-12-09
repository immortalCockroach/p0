// Implementation of a KeyValueServer. Students should write their code in this file.

package p0

import (
	"errors"
	"fmt"

	"net"

	"strconv"
	"bufio"
	"bytes"
)

var PUT string = "put"
var GET string = "get"

type keyValueServer struct {
	Closed bool
}

// New creates and returns (but does not start) a new KeyValueServer.
func New() KeyValueServer {

	return &keyValueServer{false}
}

func (kvs *keyValueServer) Start(port int) error {
	if kvs.Closed == true {
		return errors.New("the server has been closed")
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	init_db()
	kvstore["1"] = []byte("213123")
	go kvs.acceptRequest(listener)

	return nil
}

func (kvs *keyValueServer) Close() {
	// TODO: implement this!
}

func (kvs *keyValueServer) Count() int {
	// TODO: implement this!
	return -1
}

// Go routine used to accpet multiple connections
func (kvs *keyValueServer) acceptRequest(listener *net.TCPListener) {

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("listener failed to accept conn", err)
			continue
		}
		fmt.Println("connection accepted")

		go kvs.processRequest(conn)
	}
}

func (kvs *keyValueServer) processRequest(conn net.Conn) {
	messageReader := bufio.NewReader(conn)
	for {
		res, err := messageReader.ReadBytes(byte('\n'))

		if err != nil {
			return
		}

		bytesParamArray := bytes.Split(res[:len(res) - 1], []byte(","))
		fmt.Println(string(bytesParamArray[0]))

		switch string(bytesParamArray[0]) {

		case PUT:
			kvstore[string(bytesParamArray[1])] = bytesParamArray[2]
		case GET:
			conn.Write(append(kvstore[string(bytesParamArray[1])], byte('\n')))
		}
	}
}

func ReadLine(conn net.Conn) ([]byte, error) {
	messageReader := bufio.NewReader(conn)

	res, err := messageReader.ReadBytes(byte('\n'))

	if err != nil {
		return nil, err
	}

	return res[:len(res) - 1], nil
}
