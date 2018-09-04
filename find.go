package pixie

import (
	"log"
	"strings"
)

var (
	bg    = "background"
	class = "class"
	item  = "item"
	race  = "race"
	rule  = "rule"
	spell = "spell"
)

func init() {
	// load various data files
	dataFiles := DataFiles()
	if len(dataFiles) == 0 {
		log.Printf("No data files found; ALL find functions will fail")
		return
	}
}

// Entry houses the data for a particular entry parsed from the data files
type Entry struct {
	Name   string
	Page   string
	Source string
}

// FindBackground attempts to find the specified background (source book and page number)
func FindBackground(input []string) (string, *BotError) {
	return find(bg, strings.Join(input, " "))
}

// FindClass attempts to find the specified class (source book and page number)
func FindClass(input []string) (string, *BotError) {
	return find(class, strings.Join(input, " "))
}

// FindItem attempts to find the specified item (source book and page number)
func FindItem(input []string) (string, *BotError) {
	return find(item, strings.Join(input, " "))
}

// FindRace attempts to find the specified race (source book and page number)
func FindRace(input []string) (string, *BotError) {
	return find(race, strings.Join(input, " "))
}

// FindRule attempts to find the specified rule (source book and page number)
func FindRule(input []string) (string, *BotError) {
	return find(rule, strings.Join(input, " "))
}

// FindSpell attempts to find the specified spell (source book and page number)
func FindSpell(input []string) (string, *BotError) {
	return find(spell, strings.Join(input, " "))
}

// find attempts to locate various items in the local data files
func find(entryType string, query string) (string, *BotError) {

	return "", nil
}
