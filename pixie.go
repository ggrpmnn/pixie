package pixie

import (
	"fmt"
	"strings"
)

// Command represents the contents of a command sent to the bot
type Command struct {
	Task     string
	Params   []string
	Response []byte
}

// ValidCommands contains the list of valid commands that the bot can process
var ValidCommands map[string]func([]string) (string, *BotError)

func init() {
	// initialize the list of valid commands
	ValidCommands = make(map[string]func([]string) (string, *BotError))
	ValidCommands["bg"] = nil    //TODO
	ValidCommands["class"] = nil //TODO
	ValidCommands["item"] = nil  //TODO
	ValidCommands["race"] = nil  //TODO
	ValidCommands["roll"] = Roll
	ValidCommands["rule"] = nil  //TODO
	ValidCommands["spell"] = nil //TODO
}

// ParseCommand takes a string sent to the bot and parses it for processing
func ParseCommand(cmd string) (*Command, *BotError) {
	strTokens := strings.SplitN(cmd, " ", 2)
	// validate that the command is something we can process
	val := strings.ToLower(strTokens[0])
	if _, ok := ValidCommands[val]; !ok {
		return nil, &BotError{err: fmt.Sprintf("received invalid command '%s'", val),
			botMsg: fmt.Sprintf("Sorry! I don't know what to do with '%s'!", val)}
	}

	return &Command{Task: strTokens[0], Params: strTokens[1:], Response: nil}, nil
}

// RunCommand executes the command string
func RunCommand(cmdStr string) (string, *BotError) {
	cmd, err := ParseCommand(cmdStr)
	if err != nil {
		return "", err
	}

	// get the func tied to the task and execute it to get the output
	fn := ValidCommands[cmd.Task]
	output, err := fn(cmd.Params)
	if err != nil {
		return "", err
	}

	return output, nil
}
