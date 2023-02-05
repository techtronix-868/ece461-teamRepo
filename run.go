package main

import (
    "fmt"
    "os"
    "bufio"
    "regexp"
)

func seperateLinks(links[] string) ([]string,[]string){
    var re = regexp.MustCompile(`(?m)github`)

    var github[]string
    var npm[]string
    for _, line := range links {
        if re.MatchString(line){
            github = append(github,line)
        }else{
            npm = append(npm,line)
        }
    }
    return github,npm
}

func readInput(inputFile string)[]string{
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

    return fileLines
}


func main(){
    inputFile := os.Args[1]
    links:=readInput(inputFile)
    githubLinks,npmLinks := seperateLinks(links)
    fmt.Println(githubLinks)
    fmt.Println(npmLinks)
}