package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/youssefmaher99/equilib/internal/server"
)

type ServersJSON struct {
	Servers []string `json="servers"`
}

func loadServers(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file_byte_value, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var servers ServersJSON
	err = json.Unmarshal(file_byte_value, &servers)
	if err != nil {
		log.Fatal(err)
	}

	err = parseIpAndPort(servers.Servers)
	if err != nil {
		log.Fatal(err)
	}

	// add http:// to all addresses
	return servers.Servers
}

func parseIpAndPort(addresses []string) error {

	type pair struct {
		ip   string
		port string
	}

	pairs := make([]pair, 0, len(addresses))
	for _, address := range addresses {
		vals := strings.Split(address, ":")
		pairs = append(pairs, pair{ip: vals[0], port: vals[1]})
	}

	for _, pair := range pairs {
		valid := net.ParseIP(pair.ip)
		if valid == nil {
			return fmt.Errorf("invalid ip address : %s", pair.ip)
		}
		val, err := strconv.Atoi(pair.port)
		if err != nil {
			return fmt.Errorf("invalid port number : %s", pair.port)
		}
		if val > 65535 || val < 1 {
			return fmt.Errorf("invalid port number : %s", pair.port)

		}
	}
	return nil
}

func main() {
	servers_list := loadServers("servers.json")
	s := server.New("127.0.0.1:8080", len(servers_list), servers_list, 1)
	log.Fatal(s.Start())
}
