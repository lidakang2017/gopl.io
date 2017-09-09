// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"io"
	"log"
	"net"
	"os"
	"fmt"
	//"sync"
)

//!+
func main() {
	done := make(chan struct{},2)

	for server := range Genservers() {
		go func(server string) {
			conn, err := net.Dial("tcp", server)
			fmt.Printf("connecting to server: %d \n", server)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Time from server %d \n", server)
			io.Copy(os.Stdout, conn) // NOTE: ignoring errors
			log.Println("done.")
			done <- struct{}{} // signal the main goroutine
		}(server)
	}

	select {
	case <-done:
		close(done)
	}
}

func Genservers() <-chan string {
	ch := make(chan string)
	go func() {
		for _, url := range []string{
			"localhost:8000",
			"localhost:8001",
		} {
			ch <- url
		}
		close(ch)
	}()
	return ch
}
