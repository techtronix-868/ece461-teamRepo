package main

import (
	np "app/npm"
    nd "app/output"
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func seperateLinks(links[] string) ([]*nd.NdJson){
    var re = regexp.MustCompile(`(?m)github`)

    var scores[]*nd.NdJson
    for _, url := range links {
        if re.MatchString(url){
            // urlScore:= scoreGithub(url)
            //fmt.Println(url)
            //scores = append(scores,url)
        }else if strings.Contains(url,"npm"){
            // urlScore:= scoreNPM(url)
			//fmt.Println(url)
            cn := new(np.Connect_npm)
            scores = append(scores,cn.Data(url))
        }
    }
    return scores
}

func readInput(inputFile string)[]string{
    readfile,err := os.Open(inputFile)

    if err != nil {
        log.Println("error in oprning file")
        return nil
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
	logfile := os.Getenv("LOG_FILE")
	// if logfile == "" {
	// 	logfile = "./app.log"
	// }

	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	defer f.Close()

	log.SetOutput(f)

    inputFile := os.Args[1]
    links:=readInput(inputFile)
    if links == nil {
    	return
    }
    score := seperateLinks(links)
    output:=nd.FormattedOutput(score)
    fmt.Println(output)
   
}
