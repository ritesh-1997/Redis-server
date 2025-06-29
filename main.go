package main

import (
    "fmt"
	"io"
	"net"
	"os"
	"bufio"
	"strings"
	"strconv"

)

func parseInputData(buffer []byte, input int) {
	fmt.Println("Received data:", string(buffer[:input]))
	reader := bufio.NewReader(strings.NewReader(string(buffer[:input])))
	// Reading the first byte to check if it is a valid request
	// Redis RESP protocol starts with '*'
	b, _ := reader.ReadByte()
	if b != '*'{
		fmt.Println("Invalid request, expected '*' at the start of the request")
		os.Exit(1)
	}
	// "*3\r\n$3\r\nset\r\n$4\r\nkey1\r\n$6\r\nritesh\r\n"
	size, _ := reader.ReadByte()
	strSize, _ := strconv.ParseInt(string(size), 10, 64)
	fmt.Println("Size of the request:", strSize)
	reader.ReadByte()
	reader.ReadByte()

	b, _ = reader.ReadByte()
	if b != '$' {
		fmt.Println("Invalid request, expected '$' before the command")
		os.Exit(1)
	}	

	// Read the command size
	commandSize, _ := reader.ReadByte()
	commandSizeInt, _ := strconv.ParseInt(string(commandSize), 10, 64)
	fmt.Println("Command size:", commandSizeInt)
	// name := make([]byte, strSize)
	// reader.Read(name)
	// fmt.Println(string(name))
}

func handleConnection(connection net.Conn) {
    defer connection.Close()

    for {
        resp := NewResp(connection)
        value, err := resp.Read()
        if err != nil {
            if err == io.EOF {
                fmt.Println("Client disconnected, when closing connection")
                break
            }
            fmt.Println("Error reading from client:", err.Error())
            return
        }

        if value.typ != "array" {
            fmt.Println("Invalid request, expected array")
            continue
        }

        if len(value.array) == 0 {
            fmt.Println("Invalid request, expected non-empty array")
            continue
        }

        command := strings.ToUpper(value.array[0].bulk)
        args := value.array[1:]

        writer := NewWriter(connection)

        handler, ok := Handlers[command]
        if !ok {
            fmt.Println("Unknown command:", command)
            writer.Write(Value{typ: "string", str: ""})
            continue
        }
        result := handler(args)
        writer.Write(result)
    }
}


func main() {

	// Create a TCP listener on port 6379
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error starting server on PORT 6379:", err)
		os.Exit(1)
	}

	fmt.Println("Listening on port 6379 ...")

	defer listener.Close()

	// Accept incoming connections
	
	for {
        connection, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(connection)
    }	
}
