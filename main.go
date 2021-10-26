package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Check1918 - is the given IP in RFC1918 space.
func Check1918(checkIP net.IP) bool {
	localBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16"}

	for _, v := range localBlocks {
		_, vCIDR, err := net.ParseCIDR(v)
		if err != nil {
			log.Fatal(err)
		}
		if vCIDR.Contains(checkIP) {
			// fmt.Println(v, checkIP)
			return true
		}
	}

	return false
}

// DNSResolver will detect IP vs FQDN and resolve the corresponding value
func DNSResolver(resolveVal string) ([]string, error) {
	var returnSlice []string
	checkIP := net.ParseIP(resolveVal)
	if checkIP.IsLoopback() {
		return returnSlice, errors.New("IP is loopback")
	}
	if checkIP != nil {
		if checkIP.IsGlobalUnicast() && !Check1918(checkIP) {
			fmt.Printf("Resolving : %s\n\n", resolveVal)
			returnHosts, err := net.LookupAddr(resolveVal)
			if err != nil {
				return returnSlice, err
			}
			for _, v := range returnHosts {
				returnSlice = append(returnSlice, string(v))
			}
		} else {
			if Check1918(checkIP) {
				fmt.Printf("Resolving : %s\n\n", resolveVal)
				returnHosts, err := net.LookupAddr(resolveVal)
				if err != nil {
					return returnSlice, err
				}
				for _, v := range returnHosts {
					returnSlice = append(returnSlice, string(v))
				}
			}
		}
	} else {
		fmt.Printf("Resolving : %s\n\n", resolveVal)
		var err error
		returnSlice, err = net.LookupHost(resolveVal)
		if err != nil {
			return returnSlice, err
		}
	}

	return returnSlice, nil
}

func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func FileIsADirectory(file string) bool {
	if stat, err := os.Stat(file); err == nil && stat.IsDir() {
		// path is a directory
		return true
	}
	return false
}

// FileExistsAndIsNotADirectory - tests a file
func FileExistsAndIsNotADirectory(file string) bool {
	if FileExists(file) && !FileIsADirectory(file) {
		return true
	}
	return false
}

func FileReadReturnLines(fileName string) ([]string, error) {
	var lines []string
	if !FileExists(fileName) {
		return lines, errors.New("file does not exist")
	}

	file, err := os.Open(fileName)
	if err != nil {
		return lines, errors.New("could not open file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil

}

func readFile(inFile string) ([]string, error) {

	fmt.Println("reading: " + inFile)
	fileContent, err := FileReadReturnLines(inFile)
	if err != nil {
		return fileContent, errors.New("could not read config")
	}

	fmt.Printf("Read %d lines\n", (len(fileContent)))

	return fileContent, nil

}

func main() {

	var hostListFile string
	flag.StringVar(&hostListFile, "l", "", "Specify list of hosts to lookup")
	flag.Parse()

	if len(os.Args) > 2 && hostListFile == "" {
		fmt.Printf("Please lookup one thing at a time\n\n")
		os.Exit(0)
	}
	lookupVal := os.Args[1]

	if hostListFile == "" && len(os.Args) == 2 {
		resolvedThing, err := (DNSResolver(lookupVal))
		if err != nil {
			color.Red(fmt.Sprintf("\t%s\n", err))
		}

		for _, v := range resolvedThing {
			color.Green("\t%s\n", v)
		}
		fmt.Printf("\n\n\n")
	} else {
		if !FileExistsAndIsNotADirectory(hostListFile) {
			fmt.Println("This is not a file")
			os.Exit(0)
		}

		lookupList, err := readFile(hostListFile)
		if err != nil {
			color.Red("Could not read hosts file")
		}
		for _, i := range lookupList {
			if strings.HasPrefix(i, "#") || strings.HasPrefix(i, " ") {
				continue
			}
			resolvedThing, err := (DNSResolver(i))
			if err != nil {
				color.Red(fmt.Sprintf("\t%s\n", err))
			}

			for _, v := range resolvedThing {
				color.Green("\t%s\n", v)
			}
			fmt.Printf("\n\n\n")
		}
	}
}
