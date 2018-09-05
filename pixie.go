package pixie

import (
	"fmt"
	"log"
	"strings"
)

// Command represents a command that can be sent to the bot
type Command struct {
	Name        string
	Fn          func([]string) (string, *BotError)
	Description string
}

// ValidCommands contains the list of valid commands that the bot can process
var ValidCommands map[string]Command

func init() {
	// initialize the list of valid commands
	ValidCommands = make(map[string]Command)

	ValidCommands["help"] = Command{Name: "help", Fn: ListCommands, Description: "list available commands"}
	ValidCommands["roll"] = Command{Name: "roll", Fn: Roll, Description: "simulate dice rolls"}

	// add find commands only if data files are loaded
	if FilesLoaded() {
		ValidCommands[bg] = Command{Name: bg, Fn: FindBackground, Description: "find a source book reference for a background"}
		ValidCommands[class] = Command{Name: class, Fn: FindClass, Description: "find a source book reference for a class"}
		ValidCommands[item] = Command{Name: item, Fn: FindItem, Description: "find a source book reference for an item"}
		ValidCommands[race] = Command{Name: race, Fn: FindRace, Description: "find a source book reference for a race"}
		ValidCommands[rule] = Command{Name: rule, Fn: FindRule, Description: "find a source book reference for a rule"}
		ValidCommands[spell] = Command{Name: spell, Fn: FindSpell, Description: "find a source book reference for a spell"}
	}

	if len(ValidCommands) == 0 {
		log.Fatalf("No commands available to run; exiting")
	}
}

// ParseUserInput takes a string sent to the bot and parses it for processing
func ParseUserInput(input string) (string, []string, *BotError) {
	strTokens := strings.SplitN(input, " ", 2)
	// validate that the command is something we can process
	val := strings.ToLower(strTokens[0])
	if _, ok := ValidCommands[val]; !ok {
		return "", nil, &BotError{err: fmt.Sprintf("received invalid command '%s'", val),
			botMsg: fmt.Sprintf(":x: Sorry, I don't know what to do with '%s'!", val)}
	}

	return val, strTokens[1:], nil
}

// Run attempts to process and execute the user's input, if possible
func Run(input string) (string, *BotError) {
	name, params, err := ParseUserInput(input)
	if err != nil {
		return "", err
	}

	// get the func tied to the task and execute it to get the output
	cmd := ValidCommands[name]
	output, err := cmd.Fn(params)
	if err != nil {
		return "", err
	}

	return output, nil
}

// ListCommands lists the available commands the user can input
func ListCommands(input []string) (string, *BotError) {
	output := "Hey there. Here are the things I know how to do:"
	for key, cmd := range ValidCommands {
		if cmd.Description != "" {
			output += fmt.Sprintf("\n:arrow_forward: %s (%s)", key, cmd.Description)
		} else {
			output += fmt.Sprintf("\n:arrow_forward: %s", key)
		}
	}
	return strings.TrimSuffix(output, ", "), nil
}
