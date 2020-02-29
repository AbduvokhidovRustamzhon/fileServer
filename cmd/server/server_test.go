package main

import (
	"bufio"
	"bytes"
	"fileServer/pkg/rpc"
	"io/ioutil"
	"log"
	"net"
	"testing"
	"time"

)

func Test_DownloadFileInServer(t *testing.T) {
	const addr = "localhost:9995"
	go func() {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatalf("can't listen on %s: %v", addr, err)
		}
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Fatalf("can't accept %v\n", err)
			}
			go handleConn(conn)
		}
	}()
	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't connect to %s: %v", addr, err)
	}
	writer := bufio.NewWriter(conn)
	options := "sms.txt"
	cmd := "download"
	err = rpc.WriteLine(cmd, writer)
	if err != nil {
		t.Fatalf("can't write command %v\n", err)
	}
	reader := bufio.NewReader(conn)
	downloadBytes,  err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("can't reader file error: %v\n", err)
	}
	err = conn.Close()
	if err != nil {
		t.Fatalf("can't conn close %v\n", err)
	}
	err = ioutil.WriteFile("./files/"+options, downloadBytes, 0666)
	if err != nil {
		t.Fatalf("can't write file: %v\n", err)
	}
	downloadFile, err := ioutil.ReadFile("./files/"+options)
	if err != nil {
		t.Fatalf("can't read file download error:  %v\n",err)
	}
	if !bytes.Equal(downloadBytes,downloadFile) {
		t.Fatalf("files are not equal: %v", err)
	}
}

func Test_UploadFileInServer(t *testing.T) {
	const addr = "localhost:9994"
	go func() {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatalf("can't listen on %s: %v", addr, err)
		}
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Fatalf("can't accept %v\n", err)
			}
			go handleConn(conn)
		}
	}()

	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	writer := bufio.NewWriter(conn)
	options := "sms.txt"
	line := "upload" + ":" + options
	err = rpc.WriteLine(line, writer)
	if err != nil {
		t.Fatalf("can't send command %s to server: %v", line, err)
	}
	src, err := ioutil.ReadFile("files/sms.txt")
	if err != nil {
		log.Fatalf("Can't read file: %v",err)
	}
	_, err = writer.Write(src)
	if err != nil {
		log.Fatalf("Can't write: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatalf("Can't flush: %v", err)
	}
	err = conn.Close()
	if err != nil {
		log.Fatalf("Can't close conn: %v", err)
	}
	dst, err := ioutil.ReadFile("files/" + options)
	if err != nil {
		log.Fatalf("can't Read file: %v",err)
	}
	if !bytes.Equal(src, dst) {
		t.Fatalf("files are not equal: %v", err)
	}
}