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
	ValidCommands["bg"] = Command{}    //TODO
	ValidCommands["class"] = Command{} //TODO
	ValidCommands["help"] = Command{Name: "help", Fn: ListCommands, Description: "list available commands"}
	ValidCommands["item"] = Command{} //TODO
	ValidCommands["race"] = Command{} //TODO
	ValidCommands["roll"] = Command{Name: "roll", Fn: Roll, Description: "simulate dice rolls"}
	ValidCommands["rule"] = Command{}  //TODO
	ValidCommands["spell"] = Command{} //TODO

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
			botMsg: fmt.Sprintf("Sorry, I don't know what to do with '%s'!", val)}
	}

	return strTokens[0], strTokens[1:], nil
}

// Run attempts to process and execute the user's input, if possible
func Run(input string) (string, *BotError) {
	name, params, err := ParseUserInput(input)
	if err != nil {
		return "", err
	}

	// get the func tied to the task and execute it to get the output
	cmd := ValidCommands[name]

	//TODO: remove this block when all commands are implemented
	if cmd.Fn == nil {
		return "Sorry, but this command is coming soon! :sob:", nil
	}

	output, err := cmd.Fn(params)
	if err != nil {
		return "", err
	}

	return output, nil
}

// ListCommands lists the available commands the user can input
func ListCommands(input []string) (string, *BotError) {
	output := "Hey there. Here are the things I know how to do: "
	for key, cmd := range ValidCommands {
		if cmd.Description != "" {
			output += fmt.Sprintf("%s (%s), ", key, cmd.Description)
		} else {
			output += key + ", "
		}
	}
	return strings.TrimSuffix(output, ", "), nil
}
