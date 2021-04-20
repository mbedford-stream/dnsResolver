package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
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
				returnHosts, err := net.LookupHost(resolveVal)
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

func main() {

	if len(os.Args) > 2 {
		fmt.Printf("Please lookup one thing at a time\n\n")
		os.Exit(0)
	}
	lookupVal := os.Args[1]

	resolvedThing, err := (DNSResolver(lookupVal))
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range resolvedThing {
		fmt.Printf("\t%s\n", v)
	}
	fmt.Printf("\n\n\n")
}
