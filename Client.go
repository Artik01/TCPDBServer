package main

import (
	"net"
	"fmt"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:10257")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	
	for {
		var path string
		fmt.Print("Path to file with action (.json/.xml): ")
		fmt.Scan(&path)
		if path == "close" {
			break
		}
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}
		defer file.Close()
		finfo, err := file.Stat()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}
		buf := make([]byte, finfo.Size())
		file.Read(buf)
		
		conn.Write(buf[:len(buf)])
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}
	}
}
