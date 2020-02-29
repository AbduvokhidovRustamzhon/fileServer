package main

import (
	"bufio"
	"fileServer/pkg/rpc"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

var download = flag.String("download", "default", "Download")
var upload = flag.String("upload", "default", "Upload")
var list = flag.Bool("list", false, "List")

func main() {
	file, err := os.Create("log-client.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Can't close files: %v", err)
		}
	}()
	log.SetOutput(file)
	flag.Parse()
	var cmd, fileName string
	if *download != "default" {
		fileName = *download
		cmd = rpc.Dwn
	} else if *upload != "default" {
		cmd = rpc.Upd
		fileName = *upload
	} else if *list != false {
		cmd = rpc.List
		fileName = ""
	} else {
		return
	}
	StartingOperationsLoop(cmd, fileName)
}

func StartingOperationsLoop(cmd string, fileName string) (exit bool) {
	log.Print("client connecting")
	conn, err := net.Dial(rpc.Tcp, rpc.AddrClient)
	if err != nil {
		log.Fatalf("can't connect to %s: %v", rpc.AddrClient, err)
	}
	defer func(){
		err := conn.Close()
		if err != nil {
			log.Printf("Can't close conn: %v", err)
		}
	}()
	log.Print("client connected")
	writer := bufio.NewWriter(conn)
	line := cmd + ":" + fileName
	log.Print("command sending")
	err = rpc.WriteLine(line, writer)
	if err != nil {
		log.Fatalf("can't send command %s to server: %v", line, err)
	}
	log.Print("command sent")
	switch cmd {
	case rpc.Dwn:
		downloads(conn, fileName)
	case rpc.Upd:
		uploads(conn, fileName)
	case rpc.List:
		listFile(conn)
	case rpc.Quit:
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}

func downloads(conn net.Conn, fileName string) {
	reader := bufio.NewReader(conn)
	line, err := rpc.ReadLine(reader)
	if err != nil {
		log.Printf("can't read: %v", err)
		return
	}
	if line == rpc.CheckError+rpc.Suffix {
		log.Printf("files not such: %v", err)
		fmt.Printf("Файл с названием %s не существует на сервере, пожалуйста убедитесь в правильности введенных данных\n", fileName)
		return
	}
	log.Print(line)
	bytes, err := ioutil.ReadAll(reader) // while not EOF
	if err != nil {
		if err != io.EOF {
			log.Printf("can't read data: %v", err)
		}
	}
	log.Print(len(bytes))
	err = ioutil.WriteFile(rpc.WayForClient+fileName, bytes, 0666)
	if err != nil {
		log.Printf("can't write files: %v", err)
	}
	fmt.Printf("Файл с названием %s успешно скачан\n", fileName)

}

func uploads(conn net.Conn, fileName string) {
	options := strings.TrimSuffix(fileName, rpc.Suffix)
	file, err := os.Open(rpc.WayForClient + options)
	writer := bufio.NewWriter(conn)
	if err != nil {
		log.Print("files does not exist")
		err = rpc.WriteLine(rpc.CheckError, writer)
		fmt.Printf("Файл с названием %s не существует, убедитесь в правильности введенных данных\n", fileName)
		return
	}
	err = rpc.WriteLine(rpc.CheckOk, writer)
	if err != nil {
		log.Printf("error while writing: %v", err)
		return
	}
	log.Print(fileName)
	fileByte, err := io.Copy(writer, file)
	log.Print(fileByte)
	err = writer.Flush()
	if err != nil {
		log.Printf("can't flush: %v", err)
	}
	fmt.Printf("Файл с названием %s успешно загружен на сервер\n", fileName)

}

func listFile(conn net.Conn) {
	reader := bufio.NewReader(conn)
	line, err := rpc.ReadLine(reader)
	if err != nil {
		log.Printf("can't read: %v", err)
		return
	}
	fmt.Println("Список доступных файлов в сервере")
	var list string
	for i := 0; i < len(line); i++ {
		if string(line[i]) == " " || string(line[i]) == rpc.Suffix {
			fmt.Println(list)
			list = ""
		} else {
			list = list + string(line[i])
		}
	}
	_, err = ioutil.ReadAll(reader)
	if err != nil {
		if err != io.EOF {
			log.Printf("can't read data: %v", err)
		}
	}
}
