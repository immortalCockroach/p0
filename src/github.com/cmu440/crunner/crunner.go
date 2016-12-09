package main

import (
    "fmt"
    "net"

    "os"
    "p0/src/github.com/cmu440/p0"
)

const (
    defaultHost = "localhost"
    defaultPort = 9999
)

// To test your server implementation, you might find it helpful to implement a
// simple 'client runner' program. The program could be very simple, as long as
// it is able to connect with and send messages to your server and is able to
// read and print out the server's response to standard output. Whether or
// not you add any code to this file will not affect your grade.
func main() {
    fmt.Println("Not implemented.")

    tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
    //checkError(err)
    conn, _ := net.DialTCP("tcp", nil, tcpAddr)
    fmt.Println(conn)
    //checkError(err)


    _, _= conn.Write([]byte("get,1\n"))
    result, _ := p0.ReadLine(conn)
    fmt.Println(string(result))
    fmt.Println("get ok")

    _, _= conn.Write([]byte("put,1,123\n"))
    fmt.Println("put ok")

    _, _= conn.Write([]byte("get,1\n"))
    fmt.Println("get ok")

    result, _ = p0.ReadLine(conn)
    //checkError(err)
    fmt.Println(string(result))


    os.Exit(0)
}
