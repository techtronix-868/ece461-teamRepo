package main

import (
	np "app/npm"
    nd "app/output"
    log "app/lg"
    gh "app/github"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

)


func seperateLinks(links[] string) ([]*nd.NdJson){
    var re = regexp.MustCompile(`(?m)github`)

    var scores[]*nd.NdJson
    for _, url := range links {
        if re.MatchString(url){
            //log.InfoLogger.Println("Github Condition in Seperate Links , Current URL: ",url)
            //fmt.Println(url)
            scores = append(scores,gh.Score(url))
        } else if strings.Contains(url,"npm"){
            //log.InfoLogger.Println("NPM Condition in Seperate Links , Current URL: ",url)
             // urlScore:= scoreNPM(url)
		 	//fmt.Println(url)
            cn := new(np.Connect_npm);
            resp := cn.Data(url)
            if  resp != nil {
                scores = append(scores,resp)
            }
        }
    }
    return scores
}

func readInput(inputFile string)[]string{
    readfile,err := os.Open(inputFile)
    log.ErrorLogger.Println("error in opeing file: ",inputFile)

    if err != nil {
        log.ErrorLogger.Println("error in opeing file")
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
    log.Init(os.Getenv("LOG_FILE"))
    inputFile := os.Args[1]
    
    links:=readInput(inputFile)
    if links == nil {
    	return
    }
    score := seperateLinks(links)
    output:=nd.FormattedOutput(score)
    fmt.Println(output)

    os.Exit(0)
   
}
