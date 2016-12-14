package main

import (
    "fmt"
    "net"


    "time"
    "bufio"
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

    messageReader := bufio.NewReader(conn)

    _, _= conn.Write([]byte("put,1,123\n"))

    _, _= conn.Write([]byte("get,1\n"))
    result, _ := ReadLine(messageReader)

    _, _= conn.Write([]byte("put,1,123213\n"))

    fmt.Println(string(result))
    fmt.Println("get ok")

    fmt.Println("put ok")

    _, _= conn.Write([]byte("get,1\n"))
    result, _ = ReadLine(messageReader)
    fmt.Println("get ok")


    //checkError(err)
    fmt.Println(string(result))

    time.Sleep(10 * time.Second)
}

func ReadLine(messageReader *bufio.Reader) ([]byte, error) {


    res, err := messageReader.ReadBytes(byte('\n'))

    if err != nil {
        return nil, err
    }

    return res[:len(res) - 1], nil
}
