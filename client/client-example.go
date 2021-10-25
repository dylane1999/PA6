package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// create a connection to the server runnin on port 11111 on localhost
	connection, err := net.Dial("tcp", "localhost:35703")
	checkError(err)          // check for errors in connection
	defer connection.Close() // defer to system for when to close connection

	inputFilePath := getInputFileName() // get the input file path
	inputFile := getFile(inputFilePath) // get the input file
	defer inputFile.Close()
	fileSize := getFileSize(inputFile) // get input file size

	fmt.Println("write file size to server")
	connectionWriter := bufio.NewWriter(connection) // open connection writer
	connectionWriter.WriteString(fileSize)          // send the file size to the connection
	connectionWriter.Flush()

	fmt.Println("write file to server")
	fileContent := getFileAsString(inputFile) // get the file as a string
	for _, element := range fileContent {
		_, writeError := connectionWriter.WriteString(fmt.Sprintf("%s\n", element)) // write the file content to the connection
		checkError(writeError)
		connectionWriter.Flush()
	}
	connectionWriter.WriteString("done \n")          // signal finish
	connectionWriter.Flush()


	reader := bufio.NewReader(connection)
	fileSizeMessage, readErr := reader.ReadString('\n')
	checkError(readErr)
	fmt.Println("server responded with: ", fileSizeMessage) // print the servers response
}

//function to take all the text in a file and concat it to one string.
func getFileAsString(file *os.File) []string {
	var fileLines = []string{} // init slice for file lines
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text()) // for each line add the text to the file lines slice
	}
	return fileLines //strings.Join(fileLines, "\n") // return the file array as a string joined w/ newlines
}

// function used to get the input file from user input
func getInputFileName() string {
	fmt.Println("What is the upload filename?")
	inputFilePath := ""
	fmt.Scanf("%s", &inputFilePath) // get the input file from user input
	return inputFilePath
}

// utility function that takes a file path and returns an *os.File struct
func getFile(filePath string) *os.File {
	file, err := os.Open(filePath) // open file and check error
	checkError(err)
	return file
}

// utility function to get the file size
func getFileSize(file *os.File) string {
	// get the file info struct and then get file size
	fileInfo, err := file.Stat()
	checkError(err)
	sizeAsString := fmt.Sprintf("%s\n", strconv.Itoa(int(fileInfo.Size())))
	return sizeAsString
}

//1. connects to the server implemented and run on the
//workstation at port 11111 DONE
//2. prompts the user for the upload filename DONE
//3. sends first the file size (just the number in a single line) DONE
//4. sends next the file content (the entire file) DONE
//5. receives a message back from the server DONE
//6. prints what the server says DONE
//7. closes the connection and terminates the program DONE
