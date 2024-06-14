package config

import (
	"fmt"
	"commandOutputProcessorForZabbix/file"

	"gopkg.in/yaml.v3"
)

// Interfaces

type Configuration interface {
	Print()
	Get() Config
}

// Classes definition

type AppConfig struct {
	config Config
}

type CommandArguments []struct {
	Arg string `yaml:"arg"`
} 

type Config struct {
	CommandLine      string `yaml:"commandLine"`
	CommandArguments CommandArguments `yaml:"commandArguments"`
	CommandResultWaitTimeSeconds int `yaml:"commandResultWaitTimeSeconds"`
	CommandResultHasHeaderLine bool `yaml:"commandResultHasHeaderLine"`
	DataRecordsStartLine int `yaml:"dataRecordsStartLine"`

}

// Constructor

func (ac *AppConfig) NewConfig(filePath string) {

	configFile := file.NewFile(filePath)

	var config Config

	err := yaml.Unmarshal(configFile.FileContent, &config)

	if err != nil {

		fmt.Printf("Error in configuration file parsing")

	}

	ac.setConfig(config)

}

// Interfaces implementation

func (ac AppConfig) Print() {

	fmt.Print(ac.Get())
}

func (ac *AppConfig) Get() Config {

	return ac.config

}

// Private methods

func (ac *AppConfig) setConfig(config Config) {

	ac.config = config

}
