package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

var banner = `
  ▄████  ▒█████   ██▀███   ▄▄▄       ███▄    █   ▄████ 
 ██▒ ▀█▒▒██▒  ██▒▓██ ▒ ██▒▒████▄     ██ ▀█   █  ██▒ ▀█▒
▒██░▄▄▄░▒██░  ██▒▓██ ░▄█ ▒▒██  ▀█▄  ▓██  ▀█ ██▒▒██░▄▄▄░
░▓█  ██▓▒██   ██░▒██▀▀█▄  ░██▄▄▄▄██ ▓██▒  ▐▌██▒░▓█  ██▓
░▒▓███▀▒░ ████▓▒░░██▓ ▒██▒ ▓█   ▓██▒▒██░   ▓██░░▒▓███▀▒
 ░▒   ▒ ░ ▒░▒░▒░ ░ ▒▓ ░▒▓░ ▒▒   ▓▒█░░ ▒░   ▒ ▒  ░▒   ▒ 
  ░   ░   ░ ▒ ▒░   ░▒ ░ ▒░  ▒   ▒▒ ░░ ░░   ░ ▒░  ░   ░ 
░ ░   ░ ░ ░ ░ ▒    ░░   ░   ░   ▒      ░   ░ ░ ░ ░   ░ 
      ░     ░ ░     ░           ░  ░         ░       ░
`

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func expandCidr(cidr string, result chan []string) {
	ip, ipnet, _ := net.ParseCIDR(cidr)
	var ipList []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ipList = append(ipList, ip.String())
	}
	result <- ipList
}

func CidrIp(files string) {
	result := make(chan []string)
	file, err := os.Open(files)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		go expandCidr(scanner.Text(), result)
		for _, ip := range <-result {
			fmt.Println(ip)
		}
	}
}

func main() {
	var list string
	fmt.Print(banner)
	fmt.Print("Enter Ur List : ")
	_, _ = fmt.Scan(&list)
	CidrIp(list)
}
