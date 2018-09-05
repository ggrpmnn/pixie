package pixie

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	bg    = "background"
	class = "class"
	item  = "item"
	race  = "race"
	rule  = "rule"
	spell = "spell"

	lookup Lookup

	filesLoaded = false
)

// Lookup provides an interface for looking up data on different entities
type Lookup map[string]Entry

// Entry houses the data for a particular entry parsed from the data files
type Entry struct {
	Name   string
	Type   string
	Page   string
	Source string
}

func init() {
	lookup = make(Lookup, 0)

	files, err := ioutil.ReadDir("./data")
	if err != nil || len(files) == 0 {
		log.Printf("Failed to load data files; disabling reference functionality")
		return
	}

	for _, file := range files {
		fileName := file.Name()
		fileHandle, err := os.Open("./data/" + fileName)
		if err != nil {
			log.Printf("Failed to load file '%s'", fileName)
			continue
		}
		reader := csv.NewReader(fileHandle)
		data, err := reader.ReadAll()
		if err != nil {
			log.Printf("Failed to read file '%s'", fileName)
			continue
		}
		lookup.LoadFromCSV(data)
	}

	filesLoaded = true
}

// LoadFromCSV adds data to a Lookup object
func (l Lookup) LoadFromCSV(data [][]string) {
	// trim header row
	data = data[1:][:]
	for _, row := range data {
		name := row[1]
		if val, ok := l[name]; ok {
			log.Printf("Found colliding values for '%s'; original is type '%s' and new is type '%s'; skipping", name, val.Type, "")
			continue
		}
		l[name] = Entry{Name: name, Type: row[0], Page: row[2], Source: row[3]}
	}
}

// FilesLoaded returns whether the files were loaded correctly or not; set in init()
func FilesLoaded() bool {
	return filesLoaded
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
	query = strings.ToLower(query)
	if val, ok := lookup[query]; ok && entryType == val.Type {
		return fmt.Sprintf("Got it! %s is on p. %s of %s :smile:", val.Name, val.Page, val.Source), nil
	}
	return "", &BotError{err: fmt.Sprintf("Failed to find %s-type entry for '%s'", entryType, query),
		botMsg: fmt.Sprintf("Sorry! I couldn't find a %s called '%s'... :cry:", entryType, query)}
}
