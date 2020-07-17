package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"encoding/json"

	covidData "../lib"
)

var covidDataset = covidData.Load("../covid_final_data.csv")

// This program implements a Covid lookup service
// over TCP or Unix Data Socket. It loads CSV Covid Dataset
// information using package lib (see above) and uses a simple
// text-based protocol to interact with the client and send
// the data.
//
// Clients send covid dataset search requests as a textual command in the form:
//
// > nc localhost 4040
// > {"query": {"region": "Sindh"}}
// > {"query": {"date": "4/4/2020"}}

//
// When the server receives the request, it is parsed and is then used
// to search the list of covid Dataset. The search result is then printed
// JSON Format dataset to the client.
//
// Focus:
// This version of the server uses TCP sockets (or UDS) to implement a simple
// text-based application-level protocol. There are no streaming strategy
// employed for the read/write operations. Buffers are read in one shot
// creating opportunities for missing data during read.
//
// Testing:
// Netcat or telnet can be used to test this server by connecting and
// sending command using the format described above.
//
// Usage: server [options]
// options:
//   -e host endpoint, default ":4040"
//   -n network protocol [tcp,unix], default "tcp"

func main() {
	var addr string
	var network string
	flag.StringVar(&addr, "e", ":4040", "service endpoint [ip address or socket path]")
	flag.StringVar(&network, "n", "tcp", "network protocol [tcp,unix]")
	flag.Parse()
	
	// validate supported network protocols
	switch network {
	case "tcp", "tcp4", "tcp6", "unix":
	default:
		log.Fatalln("unsupported network protocol: ", network)
	}

	// create a listener for provided network and host address
	ln, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal("failed to create listener", err)
	}
	defer ln.Close()
	log.Println("*** COVID Dataset TCP Server ***")
	log.Printf("Service Started: (%s) %s \n", network, addr)
	
	// connection-loop - handle incoming requests
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			if err := conn.Close(); err != nil {
				log.Println("failed to close listner: ", err)
			}
			continue
		}
		log.Println("Connected to", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
	}()

	if _, err := conn.Write([]byte("Covid TCP Server Connected\n\n***\n")); err != nil {
		log.Println("error writing:", err)
		return
	}

	// loop to stay connected with client until client breaks connection
	for {
		// buffer for client command
		cmdLine := make([]byte, (1024 * 4))
		n, err := conn.Read(cmdLine)
		if n == 0 || err != nil {
			log.Println("Connection read error:", err)
			return
		}
		cmd, param := parseCommand(string(cmdLine[0:n]))
		if cmd == "" {
			if _, err := conn.Write([]byte("Invalid command\n")); err != nil {
				log.Println("failed to write:", err)
				return
			}
			continue
		}

		// execute command
		switch strings.ToUpper(cmd) {
		case "GET":
			result := covidData.Find(covidDataset, param)
			if len(result) == 0 {
				if _, err := conn.Write([]byte("Nothing found\n")); err != nil {
					log.Println("failed to write:", err)
				}
				continue
			}
			//send each covid info in form of JSON 
			
			var res []byte
			res, err := json.MarshalIndent(result, "", "    ")
			if err != nil {
				log.Println(err)
			}
			conn.Write([]byte(
				fmt.Sprintf("{\"response:\" %s}\n***\n",string(res)),
			))
			
		default:
			if _, err := conn.Write([]byte("Invalid Command\n")); err != nil {
				log.Println("Failed to write:", err)
				return
			}
		}
	}
}

func parseCommand(cmdLine string) (cmd, param string) {
	parts := strings.Split(cmdLine, " ")
	if len(parts) != 3 {
		return "", ""
	}
	i := 0
	for len(parts) != i {
		fmt.Println(strings.TrimSpace(parts[i]),"\t", i,"\n")
		i++
	}
	cmd1 := strings.TrimSpace(parts[0])
	cmd2 := strings.TrimSpace(parts[1])
	cmd3 := strings.TrimSpace(parts[2])

	str1 := "{\"query\":"
	str2 := "{\"date\":"
	str3 := "{\"region\":"
	str4 := "\"}"
	flag := true
	result1 := cmd1 == str1 
    	result2 := cmd2 == str2 
    	result3 := cmd2 == str3
	
	fmt.Println("\nResult 1: ", result1) 
    	fmt.Println("Result 2: ", result2) 
    	fmt.Println("Result 3: ", result3)
	//j := 0
	length := len(cmd3)-1
	if (result1 && (result2 || result3)){
		fmt.Println("Initial checks complete")
		i := 0
		for len(cmd3) != i {
			fmt.Println(cmd3[i],"\t", i,"\n")
			if(cmd3[0]!=str4[0]){
				flag = false	
			}
			if(cmd3[length]!=str4[1] || cmd3[length-1]!=str4[1]){
				flag = false	
			}
			//if(cmd3[i]!=str4[0] || cmd3[i]!=str4[1] || cmd3[i]!=str4[1]){
			//	para[j]=cmd3[i]
			//	j++ 
			//}else {
			//	flag = false
			//}
			i++
		}
	}else{
		flag = false
	}
	if(flag == false){
		return ""," "
	}
	
	cmd = "Get"
	param = cmd3[2:length-2]
	return
}
