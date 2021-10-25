package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// start server on port 8080 and get connection
	fmt.Println("Launching server on port:", 35703)
	ln, _ := net.Listen("tcp", ":35703")
	defer ln.Close()

	for {
		connection, connErr := ln.Accept()
		checkError(connErr)
		defer connection.Close()
		go handleConnection(connection)
	}

}


func handleConnection(connection net.Conn){
	// read the file size from the client and print

	reader := bufio.NewReader(connection)
	originalFileSize, readErr := reader.ReadString('\n')
	checkError(readErr)
	fmt.Printf("uploaded file size: %s", originalFileSize)

	fileSizeAsNumber, converrErr := strconv.Atoi(strings.Split(originalFileSize, "\n")[0])
	checkError(converrErr)
	totalSizeUploaded := 0
	clientUpload := []string{}
	for {
		line, err := reader.ReadString('\n')
		totalSizeUploaded += len(line)
		clientUpload = append(clientUpload, line) // add the line to upload
		if totalSizeUploaded >= fileSizeAsNumber {
			break
		}
		if len(line) == 0 && err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
	}

	fmt.Println("Client uploaded content:")
	fmt.Println(strings.Join(clientUpload, "")) // print the servers response

	newFile := getFileWithLineNumbers(clientUpload, "output.txt")
	newFileSize := getFileSize(newFile)

	time.Sleep(5 * time.Second)
	connectionWriter := bufio.NewWriter(connection)
	message := fmt.Sprintf("original file size: %s. new file size: %s \n", strings.Split(originalFileSize, "\n")[0], newFileSize)
	connectionWriter.WriteString(message)
	connectionWriter.Flush()
}

// function that appends line numbers to the contents of a given file
// then writes the new file to the given file path
func getFileWithLineNumbers(fileLines []string, outputFileName string) *os.File {
	currentLineNumber := 1                             // line number of text
	lineNumberFile, error := os.Create(outputFileName) // create file and check error
	checkError(error)

	// for each line write to the new file the line number + the original text
	writer := bufio.NewWriter(lineNumberFile)
	for _, line := range fileLines {
		str := fmt.Sprintf("%d %s", currentLineNumber, line)
		writer.WriteString(str)
		currentLineNumber += 1
	}

	writer.Flush()
	return lineNumberFile
}

// utility function to get the file size
func getFileSize(file *os.File) string {
	// get the file info struct and then get file size
	fileInfo, err := file.Stat()
	checkError(err)
	sizeAsString := strconv.Itoa(int(fileInfo.Size()))
	return sizeAsString
}
