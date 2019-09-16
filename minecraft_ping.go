package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/whatupdave/mcping"
)

// PingResult holds info the result
type PingResult struct {
	isOnline    bool
	serverIP    string
	playerCount int
}

// PingServer pings the server for information
func PingServer(ip string, result chan PingResult) {
	resp, err := mcping.Ping(ip + ":25565")

	tmp := PingResult{isOnline: true, serverIP: ip, playerCount: resp.Online}

	if err != nil {
		tmp.isOnline = false
	}

	result <- tmp
}

func loadFileLines(filename string) []string {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return []string{}
	}

	contents := string(bytes)

	return strings.Split(contents, "\n")
}

func main() {
	fmt.Print("MINECRAFT PINGER\n")

	// Show the usage
	if len(os.Args) < 2 {
		fmt.Println("Usage: minecraft_ping [file]")

		fmt.Println("\nThe file must contain a list of")
		fmt.Println("minecraft server IP's which are")
		fmt.Println("seperated by a newline.")

		os.Exit(0)
	}

	// Load the ip file
	fmt.Println("> Loading file")
	ips := loadFileLines(os.Args[1])

	// Start pinging the ips
	fmt.Println("> Pinging ips")
	result := make(chan PingResult)
	for _, ip := range ips {
		go PingServer(ip, result)
	}

	fmt.Println("> Waiting...")
	fmt.Println("\n[Server IP]\t[Count]")

	// Wait for the results
	finished := 0
	total := len(ips)
	for finished < total {
		server := <-result

		if server.isOnline && server.playerCount > 0 {
			fmt.Println(server.serverIP+":\t", server.playerCount)
		}

		finished++
	}

	fmt.Println("\nDone.")
}
