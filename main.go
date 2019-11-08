package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	print("Input File Name: ")
	fileName, _ := reader.ReadString('\n')
	fileName = fileName[:len(fileName)-1]

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	frameArrTemp := strings.Split(string(b), " ")
	frameArr := []string{}
	for _, seg := range frameArrTemp {
		frameArr = append(frameArr, seg[:2])
		frameArr = append(frameArr, seg[2:])
	}

	println(strings.Join(frameArr, " "), "\n\n")

	processEthernetHeader(frameArr[:14])
	processIpHeader(frameArr[14 : 14+int(frameArr[14][1]-'0')*4])
}

func processEthernetHeader(header []string) {
	println("Ethernet Header:", strings.Join(header, " "))
	println("  Destination Address:", strings.Join(header[:6], " "))
	println("  Source Address:", strings.Join(header[6:12], " "))
	payloadType := strings.Join(header[12:], " ")
	if payloadType == "08 00" {
		println("  Payload Type:", payloadType, "- IP")
	} else if payloadType == "08 06" {
		println("  Payload Type:", payloadType, "- ARP")
	}
}

func processIpHeader(header []string) {
	println("IP Header:", strings.Join(header, " "))
}
