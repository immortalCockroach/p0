p0
==

This repository contains the starter code that you will use as the basis of your key-value database server
implementation. It also contains the tests that we will use to test your implementation,
and an example 'server runner' binary that you might find useful for your own testing purposes.

If at any point you have any trouble with building, installing, or testing your code, the article
titled [How to Write Go Code](http://golang.org/doc/code.html) is a great resource for understanding
how Go workspaces are built and organized. You might also find the documentation for the
[`go` command](http://golang.org/cmd/go/) to be helpful. As always, feel free to post your questions
on Piazza as well.

## Running the official tests

To test your submission, we will execute the following command from inside the
`src/github.com/cmu440/p0` directory:

```sh
$ go test
```

We will also check your code for race conditions using Go's race detector by executing
the following command:

```sh
$ go test -race
```

To execute a single unit test, you can use the `-test.run` flag and specify a regular expression
identifying the name of the test to run. For example,

```sh
$ go test -race -test.run TestBasic1
```

**Our reference solution was tested on the AFS clusters, so we strongly recommend that you test your solution on AFS before submitting.**
## Submission

Submit the `server_impl.go` file on Autolab. **Do not change the name of the file as this will cause the tests to fail. Please submit your code with all print statements removed.**

## Testing your implementation using `srunner`

To make testing your server a bit easier (especially during the early stages of your implementation
when your server is largely incomplete), we have given you a simple `srunner` (server runner)
program that you can use to create and start an instance of your `KeyValueServer`. The program
simply creates an instance of your server, starts it on a default port, and blocks forever,
running your server in the background.

To compile and build the `srunner` program into a binary that you can run, execute the three
commands below (these directions assume you have cloned this repo to `$HOME/p0`):

```bash
$ export GOPATH=$HOME/p0
$ go install github.com/cmu440/srunner
$ $GOPATH/bin/srunner
```

The `srunner` program won't be of much use to you without any clients. It might be a good exercise
to implement your own `crunner` (client runner) program that you can use to connect with and send
messages to your server. We have provided you with an unimplemented `crunner` program that you may
use for this purpose if you wish. Whether or not you decide to implement a `crunner` program will not
affect your grade for this project.

You could also test your server using Netcat as you saw shortly in lecture (i.e. run the `srunner`
binary in the background, execute `nc localhost 9999`, type the message you wish to send, and then
click enter).

## Using Go on AFS

For those students who wish to write their Go code on AFS (either in a cluster or remotely), you will
need to set the `GOROOT` environment variable as follows (this is required because Go is installed
in a custom location on AFS machines):

```bash
$ export GOROOT=/usr/local/depot/go
```

## 实现

Server的结构在`server_impl.go`的顶部有描述.
如果需要debug运行，将server的`debugMode`修改为`true`(结构体的第二个参数)。
然后在需要调试的地方加入下面的语句。

```Golang
if debugMode {
    fmt.Println("your debug and test code")
}
```

## 和[官方实现](http://www.cs.cmu.edu/~srini/15-440/lectures/code_p0_solution.go)的区别

官方的实现中，当client主动退出的时候，client对应的`readRoutine`和`writeRoutine`
仍然处于运行中。

这里的实现是在`readRoutine`读取到`io.EOF`的时候，server的`mainRoutine`向client的
2个channel写入退出信息，2个routine读到之后return

有个坑在于`readRoutine`读取到`io.EOF`之后一直循环读到`io.EOF`,并向`exitClientChannel`循环
发送client本身。如果client再向`readRoutine`写入退出信息的话可能会死锁。

例如下面的代码：第一个是client的`readRoutine`，第二个是server的`mainRoutine`

当client发送第二个`io.EOF`的时候，server可能刚刚向`client.exitReadChannel`中写入0
而此时如果server想往下执行，那么必须在client端执行到`<-client.exitReadChannel`之后，
但是client端被`exitClientsChannel`阻塞。
而client端想往下执行，那么必须在server端再次从`exitClientsChannel`中读到才可以，
而server也因为client端无法读取到退出信息而阻塞。

最终导致了死锁

```Golang
 case <-client.exitReadChannel:
			return
	default:
				res, err := messageReader.ReadBytes(byte('\n'))

				if err == io.EOF {
					if kvs.debugMode {
						fmt.Println("a channel exit initiative")
					}
					kvs.exitClientsChannel <- client
```

```Golang
 case exitClient := <-kvs.exitClientsChannel:
    for i, client := range kvs.clients {
        if client == exitClient {
            // client端主动关闭之后，将对应的2个routine结束
            client.exitReadChannel <- 0
            client.exitWriteChannel <- 0
            kvs.clients = append(kvs.clients[:i], kvs.clients[i+1:]...)
            break
        }
    }
```



故此处设置了第一次读到`io.EOF`之后修改标志位`exited`，`default`块不再执行。

```Golang
if !exited { // 防止无限读取到EOF的死锁
				res, err := messageReader.ReadBytes(byte('\n'))

				if err == io.EOF {
					if kvs.debugMode {
						fmt.Println("a channel exit initiative")
					}
					exited = true
					kvs.exitClientsChannel <- client
```