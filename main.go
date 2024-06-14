package main

import (
	"flag"
	"fmt"
	"linuxProcessesCounterStandalone/co"
	"linuxProcessesCounterStandalone/config"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type CommandOutput struct {
	output []byte
	error  error
}

// Mandatory functions

func main() {

	// Getting current directory as a configuration file default location
	// Alternative is to use command line flags, but it doesn't work for Zabbix plugin

	currentDirectory, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
	}

	// Getting flags

	configFilePath := flag.String("configfile", currentDirectory + getOSPathDelimiter() + "config.yml", "Full path to configuration file")

	flag.Parse()

	// Getting configuration

	var appConfig config.AppConfig

	appConfig.NewConfig(*configFilePath)

	// appConfig.Print()

	// Command line mode

	if len(appConfig.Get().CommandLine) > 0 {

		// Getting delay setup

		waitTimeSeconds := appConfig.Get().CommandResultWaitTimeSeconds

		c := make(chan CommandOutput)

		// This will work in background

		go runOSCommand(c, appConfig.Get().CommandLine, appConfig.Get().CommandArguments, waitTimeSeconds)

		// When everything is done, you can check your background process result
		res := <-c

		if res.error != nil {

			fmt.Println("Failed to execute command: ", res.error)

		} else {

			// You will be here, runOSCommand has finish successfuly

			co := co.NewCommandOutput(res.output, appConfig.Get().CommandResultHasHeaderLine, appConfig.Get().DataRecordsStartLine)

			recordsCount := co.GetRecordsCount()

			fmt.Printf("%d \n", recordsCount)

		}

	}

}

func getOSPathDelimiter() string {

	var delimiter string

	delimiter = "/"
	
	if runtime.GOOS == "windows" {
	
		delimiter = "\\"

	}

	return delimiter
}

func runOSCommand(ch chan<- CommandOutput, command string, arguments config.CommandArguments, waitTimeSeconds int) {

	// Preparing command arguments

	var args []string
	for _, item := range arguments {

		args = append(args, item.Arg)

	}

	cmd := exec.Command(command, args...)

	data, err := cmd.CombinedOutput()

	if waitTimeSeconds > 0 {

		time.Sleep(time.Duration(waitTimeSeconds) * time.Second)

	}

	ch <- CommandOutput{
		error:  err,
		output: data,
	}
}
