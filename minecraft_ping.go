package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/whatupdave/mcping"
)

// PingResult info about the result
type PingResult struct {
	isOnline    bool
	serverIP    string
	playerCount int
}

// PingServer pings the server for information
func PingServer(ip string, result chan PingResult) {
	resp, err := mcping.Ping(ip + ":25565")

	tmp := PingResult{isOnline: true, serverIP: AdjustStringLength(ip, 15), playerCount: resp.Online}

	if err != nil {
		tmp.isOnline = false
	}

	result <- tmp
}

// LoadFileLines loads the given filename
// and returns the lines as an array
func LoadFileLines(filename string) []string {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return []string{}
	}

	contents := string(bytes)

	return strings.Split(contents, "\n")
}

// AdjustStringLength takes in a string, and adjusts it to
// the given length. May shorten the string, or add in whitespace.
func AdjustStringLength(input string, length int) string {
	tmp := ""

	for i := 0; i < length; i++ {
		if i < len(input) {
			tmp += string(input[i])
		} else {
			tmp += " "
		}
	}

	return tmp
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
	ips := LoadFileLines(os.Args[1])

	// Start pinging the ips
	fmt.Println("> Pinging ips")
	result := make(chan PingResult)
	for _, ip := range ips {
		go PingServer(ip, result)
	}

	fmt.Println("> Waiting...")
	fmt.Println("\n[Server IP]\t\t[Count]")

	// Wait for the results
	finished := 0
	total := len(ips)
	for finished < total {
		server := <-result

		// Only print if people are online
		if server.isOnline && server.playerCount > 0 {
			fmt.Println(server.serverIP, "\t", server.playerCount)
		}

		finished++
	}

	fmt.Println("\nDone.")
}
