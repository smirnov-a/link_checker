package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

const linksFilename = "links.txt"

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	help := flag.Bool("help", false, "Show help")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	f := flag.String("file", linksFilename, "Filename with links")
	if err := checkFileExists(*f); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	initConfig()

	startTime := time.Now()
	ctx := context.Background()
	process(ctx, *f)
	endTime := time.Now()
	executionTime := endTime.Sub(startTime)
	fmt.Printf("All done! Execution time: %v second(s)", executionTime.Seconds())
}

// initConfig initialize config with viper
func initConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %v", err)
	}
	viper.AutomaticEnv()
}

func checkFileExists(fileName string) error {
	_, err := os.Stat(fileName)
	return err
}
