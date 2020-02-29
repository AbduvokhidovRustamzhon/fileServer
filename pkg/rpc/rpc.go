package rpc

import (
	"bufio"
	"io/ioutil"
	"log"
)

func ReadLine(reader *bufio.Reader) (line string, err error) {
	return reader.ReadString('\n')
}

func WriteLine(line string, writer *bufio.Writer) (err error) {
	_, err = writer.WriteString(line + "\n")
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return
}

func ReadDir(line string) (fileList string) {
	files, err := ioutil.ReadDir(line)
	if err != nil {
		log.Printf("Can't read dir: %v", err)
	}
	for _, file := range files {
		if fileList == "" {
			fileList = fileList + file.Name()
		} else {
			fileList = fileList + " " + file.Name()
		}
	}
	fileList = fileList + "\n"
	return fileList
}


const Dwn = "Download"
const Upd = "Upload"
const List = "List"
const Quit = "quit"
const CheckError = "result: error"
const CheckOk = "result: ok"
const Tcp = "tcp"
const Suffix = "\n"
const AddrClient = "localhost:9999"

const WayForServer = "files/"
const WayInServer = "files"
const WayForClient = "files/"

