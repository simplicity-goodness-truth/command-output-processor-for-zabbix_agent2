package co

import (
	"bufio"
	"fmt"
	"strings"
)

// Types

type CommandOutputContent []byte

// Interfaces

type CommandOutputRecords interface {
	GetRecords() commandOutputRecords
	GetRecordsCount() int
	GetOutputСontent() CommandOutputContent
	GetHeaders()
	PrintRecords()
}

// Class definition

type CommandOutput struct {
	outputContent        CommandOutputContent
	records              commandOutputRecords
	recordsCount         int
	hasHeaderLine        bool
	headers              []string
	dataRecordsStartLine int
}

// Private types

type commandOutputRecords map[int]commandOutputRecord
type commandOutputRecord map[string]string

// Constructor

func NewCommandOutput(content CommandOutputContent, commandResultHasHeaderLine bool, dataRecordsStartLine int) *CommandOutput {

	return &CommandOutput{
		outputContent:        content,
		hasHeaderLine:        commandResultHasHeaderLine,
		dataRecordsStartLine: dataRecordsStartLine,
		records:              make(map[int]commandOutputRecord),
	}

}

// Interfaces implementations

func (c *CommandOutput) GetRecordsCount() (recordsCount int) {

	c.countRecords()

	return c.recordsCount

}

func (c *CommandOutput) GetHeaders() []string {

	return c.headers
}

func (c *CommandOutput) GetOutputСontent() CommandOutputContent {

	return c.outputContent
}

func (c *CommandOutput) GetRecords() commandOutputRecords {

	return c.records

}

func (c *CommandOutput) PrintRecords() {

	for key, value := range c.records {
		fmt.Printf("RECORD NUMBER: %d \n", key)

		for jey, walue := range value {
			fmt.Printf("HEADER: %s \n", jey )
			fmt.Printf("VALUE: %s \n", walue)
		}

	}
}

// Private methods

func (c *CommandOutput) countRecords() {

	if len(c.records) == 0 {

		c.convertOutputToRecords()

		c.recordsCount = len(c.records)

	}

}

func (c *CommandOutput) convertOutputToRecords() {

	var commandOutputRecordArray commandOutputRecord

	commandOutput := string(c.GetOutputСontent())

	if len(commandOutput) > 0 {

		// Getting headers

		if c.hasHeaderLine {

			c.setHeaders()

		}

		var lineNumber int

		scanner := bufio.NewScanner(strings.NewReader(commandOutput))

		for scanner.Scan() {

			if lineNumber < c.dataRecordsStartLine {

				lineNumber++

				continue
			}

			commandOutputRecordArray = make(map[string]string)

			commandOutputLine := scanner.Text()

			commandOutputRecord := strings.Fields(string(commandOutputLine))

			if c.hasHeaderLine {

				for i, item := range c.GetHeaders() {

					commandOutputRecordArray[item] = commandOutputRecord[i]

				}

			 } else {

				headerName := fmt.Sprintf("header%d", lineNumber)			

				commandOutputRecordArray[headerName] = string(commandOutputLine)
			 }



			c.records[lineNumber] = commandOutputRecordArray

			lineNumber++

		}

	}

}

func (c *CommandOutput) setHeaders() {

	for _, line := range strings.Split(strings.TrimRight(string(c.outputContent), "\n"), "\n") {

		if len(line) > 1 {

			c.headers = strings.Split(line, "\t")

			break
		}
	}
}
