package main

import (
	"bufio"
	"fmt"
	"github.com/cheggaaa/pb/v3"
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
	bar := pb.StartNew(len(result))
	for scanner.Scan() {
		go expandCidr(scanner.Text(), result)
		for _, ip := range <-result {
			w, _ := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			_, _ = w.WriteString(ip + "\n")
			bar.Increment()
		}
	}
	bar.Finish()
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	close(result)
}

func main() {
	var list string
	fmt.Print(banner)
	fmt.Print("Enter Ur List : ")
	_, _ = fmt.Scan(&list)
	CidrIp(list)
}
