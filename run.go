package main

import (
    "fmt"
    "os"
    "bufio"
    "regexp"
)

func readInput(inputFile string){
    readfile,err := os.Open(inputFile)

    if err != nil {
        fmt.Println("error in oprning file")
    }

    fileScanner := bufio.NewScanner(readfile)
    fileScanner.Split(bufio.ScanLines)

    var fileLines []string

    //The following read the file and adds to an array
    for fileScanner.Scan() {
        fileLines = append(fileLines, fileScanner.Text())
    }
    readfile.Close()

    var re = regexp.MustCompile(`(?m)github`)

    var github[]string
    var npm[]string
    for _, line := range fileLines {
        // Testing if we are able to print correct output from input file
        // fmt.Println(line)

        // Seperating github and npm in to their respecting arrays
        if re.MatchString(line){
            github = append(github,line)
        }else{
            npm = append(npm,line)
        }
    }


}

func main(){
    inputFile := os.Args[1]
    fmt.Println(inputFile)
    readInput(inputFile)
}