package main

import (
    "fmt"
    "os"
    "bufio"
)

func readInput(inputFile string){
    readfile,err := os.Open(inputFile)

    if err != nil {
        fmt.Println("error in oprning file")
    }

    fileScanner := bufio.NewScanner(readfile)
    fileScanner.Split(bufio.ScanLines)

    var fileLines []string

    for fileScanner.Scan() {
        fileLines = append(fileLines, fileScanner.Text())
    }

    readfile.Close()

    for _, line := range fileLines {
        fmt.Println(line)
    }

    fmt.Println(fileLines)

}

func main(){
    inputFile := os.Args[1]
    fmt.Println(inputFile)
    readInput(inputFile)
}