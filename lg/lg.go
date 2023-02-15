package lg

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func Init(logfile string) error {
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY, 0644)
	// fmt.Println(err)
	if err != nil {
		// fmt.Println("HERE")
		// log.Fatalf("Failed to open log file: %v", err)
		return err
	}

	InfoLogger = log.New(f, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(f, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(f, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
	
}

