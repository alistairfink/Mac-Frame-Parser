package main

import (
	// "bufio"
	"io/ioutil"
	"log"
	// "os"
	"strconv"
	"strings"
)

func main() {
	// reader := bufio.NewReader(os.Stdin)
	// print("Input File Name: ")
	// fileName, _ := reader.ReadString('\n')
	// fileName = fileName[:len(fileName)-1]
	b, err := ioutil.ReadFile("test.txt") //fileName)
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
	ipHeaderEnd := 14 + int(frameArr[14][1]-'0')*4
	processIpHeader(frameArr[14:ipHeaderEnd])
}

func processEthernetHeader(header []string) {
	println("Ethernet Header:", strings.Join(header, " "))
	println("\tDestination MAC Address:\t", strings.Join(header[:6], ":"))
	println("\tSource MAC Address:\t\t", strings.Join(header[6:12], ":"))
	payloadType := "0x" + strings.Join(header[12:], "")
	if payloadType == "0x0800" {
		println("\tPayload Type:\t\t\t", payloadType, "- IP")
	} else if payloadType == "0x0806" {
		println("\tPayload Type:\t\t\t", payloadType, "- ARP")
	}

	println()
}

func processIpHeader(header []string) {
	println("IP Header:", strings.Join(header, " "))
	println("\tVersion:\t\t\t IPv" + string(header[0][0]))
	println("\tHeader Length:\t\t\t", len(header), "Bytes")

	// Type of Service
	serviceType := hexToBin(header[1])
	println("\tType of Service:\t\t", "0b"+serviceType)
	if serviceType[:3] == "000" {
		println("\t\tRoutine Precedence")
	}

	if serviceType[3] == '0' {
		println("\t\tNormal Delay")
	} else {
		println("\t\tDelay:", serviceType[3])
	}

	if serviceType[4] == '0' {
		println("\t\tNormal Throughput")
	} else {
		println("\t\tThroughput:", serviceType[4])
	}

	if serviceType[5] == '0' {
		println("\t\tNormal Reliability")
	} else {
		println("\t\tReliability:", serviceType[5])
	}

	println("\tTotal Length:\t\t\t", "0x"+strings.Join(header[2:4], ""), "-", hexToDec(strings.Join(header[2:4], "")), "Bytes")
	println("\tIdentification:\t\t\t", "0x"+strings.Join(header[4:6], ""))

	// Flags and Offset
	flagsAndOffset := hexToBin(strings.Join(header[6:8], ""))
	println("\tFlags and Offset:\t\t", "0b"+flagsAndOffset)
	if flagsAndOffset[1] == '1' {
		println("\t\tDo Not Fragment")
	} else {
		println("\t\tOk to Fragment")
	}

	if flagsAndOffset[2] == '1' {
		println("\t\tMore Fragments to Come")
	} else {
		println("\t\tNo More Fragments")
	}

	println("\t\tFragement Offset:", int(flagsAndOffset[3]-'0'))

	println("\tTime to Live:\t\t\t", hexToDec(header[8]), "more hops")
	// Protocol
	protocol := "0x" + header[9]
	if protocol == "0x06" {
		println("\tProtocol:\t\t\t TCP")
	} else if protocol == "0x01" {
		println("\tProtocol:\t\t\t ICMP")
	} else {
		println("\tProtocol:\t\t\t UDP")
	}

	println("\tHeader Checksum:\t\t", "0x"+strings.Join(header[10:12], ""))
	println("\tSource Address:\t\t\t", hexToIp(header[12:16]))
	println("\tDestination Address:\t\t", hexToIp(header[16:20]))
}

var hexToBinConversion map[byte]string = map[byte]string{
	'0': "0000",
	'1': "0001",
	'2': "0010",
	'3': "0011",
	'4': "0100",
	'5': "0101",
	'6': "0110",
	'7': "0111",
	'8': "1000",
	'9': "1001",
	'A': "1010",
	'B': "1011",
	'C': "1100",
	'D': "1101",
	'E': "1110",
	'F': "1111",
	'a': "1010",
	'b': "1011",
	'c': "1100",
	'd': "1101",
	'e': "1110",
	'f': "1111",
}

func hexToBin(hex string) string {
	result := ""
	for i := 0; i < len(hex); i++ {
		result += hexToBinConversion[hex[i]]
	}

	return result
}

func hexToDec(hex string) int {
	bin := hexToBin(hex)
	result := 0
	bitPos := 1
	for i := len(bin) - 1; i >= 0; i-- {
		if bin[i] == '1' {
			result += bitPos
		}

		if bitPos == 0 {
			bitPos = 1
		} else {
			bitPos = bitPos << 1
		}
	}

	return result
}

func hexToIp(hex []string) string {
	if len(hex) != 4 {
		return ""
	}

	for i := 0; i < 4; i++ {
		hex[i] = strconv.Itoa(hexToDec(hex[i]))
	}

	return strings.Join(hex, ".")
}
