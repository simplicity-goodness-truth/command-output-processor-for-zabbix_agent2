package main

import (
	"commandOutputProcessorForZabbix/co"
	"commandOutputProcessorForZabbix/config"
	"flag"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"git.zabbix.com/ap/plugin-support/plugin"
	"git.zabbix.com/ap/plugin-support/plugin/container"
)

type CommandOutput struct {
	output []byte
	error  error
}

// Zabbix Plugin part

type Plugin struct {
	plugin.Base
}

var impl Plugin

var pluginConfig config.Config

func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (result interface{}, err error) {

	var metricValue int

	switch key {

	case "commandoutputrecordscount":

		var recordsCount int

		// Command line mode

		if len(pluginConfig.CommandLine) > 0 {

			// Getting delay setup

			waitTimeSeconds := pluginConfig.CommandResultWaitTimeSeconds

			c := make(chan CommandOutput)

			// This will work in background

			go runOSCommand(c, pluginConfig.CommandLine, pluginConfig.CommandArguments, waitTimeSeconds)

			// When everything is done, you can check your background process result
			res := <-c

			if res.error != nil {

				fmt.Println("Failed to execute command: ", res.error)

			} else {

				// You will be here, runOSCommand has finish successfuly

				co := co.NewCommandOutput(res.output, pluginConfig.CommandResultHasHeaderLine, pluginConfig.DataRecordsStartLine)

				recordsCount = co.GetRecordsCount()

				metricValue = recordsCount

			}

		}

	}

	return metricValue, nil

}

// Mandatory functions

func init() {

	plugin.RegisterMetrics(&impl, "CommandOutputProcessor", "commandoutputrecordscount", "Returns an amount of records from command output.")

}

func main() {

	// When running in Zabbix Agent 2 mode there is no change to execute a plugin with command line parameters
	// In addition it is not clear which is a current directory for Zabbix Agent 2
	// That is why for Windows configuration file to be put by default in C:\Program Files\Zabbix Agent 2\config.yml

	configFilePath := flag.String("configfile", getDefaultConfigPath(), "Full path to configuration file")

	flag.Parse()

	// Getting configuration

	var appConfig config.AppConfig

	appConfig.NewConfig(*configFilePath)

	pluginConfig = appConfig.Get()

	// Creating handler

	h, err := container.NewHandler(impl.Name())

	if err != nil {
		panic(fmt.Sprintf("failed to create plugin handler %s", err.Error()))
	}

	impl.Logger = &h

	err = h.Execute()

	if err != nil {
		panic(fmt.Sprintf("failed to execute plugin handler %s", err.Error()))
	}

}

func getDefaultConfigPath() string {

	var defaultConfigFilePath string

	if runtime.GOOS == "windows" {

		defaultConfigFilePath = "C:\\Program Files\\Zabbix Agent 2\\config.yml"

	} else {

		defaultConfigFilePath = "/usr/local/sbin/config.yml"

	}

	return defaultConfigFilePath

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
