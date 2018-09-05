package pixie

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

// Command represents a command that can be sent to the bot
type Command struct {
	Name        string
	Emoji       string
	Fn          func([]string) (string, *BotError)
	Description string
}

// ValidCommands contains the list of valid commands that the bot can process
var ValidCommands map[string]Command

func init() {
	// initialize the list of valid commands
	ValidCommands = make(map[string]Command)

	ValidCommands["help"] = Command{Name: "help", Emoji: "exclamation", Fn: ListCommands, Description: "list available commands"}
	ValidCommands["roll"] = Command{Name: "roll", Emoji: "game_die", Fn: Roll, Description: "simulate dice rolls"}

	// add find commands only if data files are loaded
	if FilesLoaded() {
		ValidCommands[bg] = Command{Name: bg, Emoji: "fleur_de_lis", Fn: FindBackground, Description: "find a source book reference for a background"}
		ValidCommands[class] = Command{Name: class, Emoji: "bow_and_arrow", Fn: FindClass, Description: "find a source book reference for a class"}
		ValidCommands[item] = Command{Name: item, Emoji: "gem", Fn: FindItem, Description: "find a source book reference for an item"}
		ValidCommands[race] = Command{Name: race, Emoji: "moyai", Fn: FindRace, Description: "find a source book reference for a race"}
		ValidCommands[rule] = Command{Name: rule, Emoji: "scales", Fn: FindRule, Description: "find a source book reference for a rule"}
		ValidCommands[spell] = Command{Name: spell, Emoji: "fire", Fn: FindSpell, Description: "find a source book reference for a spell"}

		ValidCommands["generate"] = Command{Name: "generate", Emoji: "", Fn: nil, Description: "generate a random character"}
	}

	if len(ValidCommands) == 0 {
		log.Fatalf("No commands available to run; exiting")
	}
}

// ListCommands lists the available commands the user can input
func ListCommands(input []string) (string, *BotError) {
	output := ":information_source: Hey there! Here are the things I know how to do:"

	keys := ListKeys()
	for _, key := range keys {
		cmd := ValidCommands[key]
		output += "\n"
		if cmd.Emoji != "" {
			output += fmt.Sprintf(":%s: ", cmd.Emoji)
		}
		output += fmt.Sprintf("**%s**", key)
		if cmd.Description != "" {
			output += fmt.Sprintf(": %s", cmd.Description)
		}
	}
	return output, nil
}

// ListKeys returns a list of keys from the ValidCommands map, sorted by alpha
func ListKeys() []string {
	keys := make([]string, len(ValidCommands))
	idx := 0
	for key := range ValidCommands {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)
	return keys
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
