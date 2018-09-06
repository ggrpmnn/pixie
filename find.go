package pixie

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	bg    = "background"
	class = "class"
	item  = "item"
	race  = "race"
	rule  = "rule"
	spell = "spell"

	dataDir = "./data/"

	filesLoaded = false

	// Lookup is a lookup table for all the data loaded from CSV
	Lookup []SourceEntry
)

// SourceEntry houses the data for a particular entry parsed from the data files
type SourceEntry struct {
	Name   string
	Type   string
	Page   string
	Source string
}

func init() {
	files, err := ioutil.ReadDir(dataDir)
	if err != nil || len(files) == 0 {
		log.Printf("Failed to load data files; disabling reference functionality")
		return
	}

	Lookup = make([]SourceEntry, 0)
	// load the data from each CSV file
	for _, file := range files {
		log.Printf("Loading data file '%s'", file.Name())
		Lookup = append(Lookup, LoadFromCSV(file)...)
	}

	filesLoaded = true
}

// FilesLoaded returns whether the files were loaded correctly or not; set in init()
func FilesLoaded() bool {
	return filesLoaded
}

// filter filters the list based on the specified function's criteria
func filter(entries []SourceEntry, criteria string, fn func(SourceEntry, string) bool) *[]SourceEntry {
	results := make([]SourceEntry, 0)
	for _, entry := range entries {
		if fn(entry, criteria) {
			results = append(results, entry)
		}
	}
	return &results
}

// filterType is a filter function for trimming entries based on type (exact match)
func filterType(entry SourceEntry, kind string) bool {
	if entry.Type == kind {
		return true
	}
	return false
}

// filterName is a filter function for trimming entries based on name (via regex)
func filterName(entry SourceEntry, nameQuery string) bool {
	match, _ := regexp.MatchString(nameQuery, entry.Name)
	return match
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
	// lookup the type, then the name
	results := filter(Lookup, entryType, filterType)
	if len(*results) > 0 {
		results = filter(*results, query, filterName)
	}

	// process and return filtered results
	numResults := len(*results)
	entries := *results
	if numResults > 0 {
		if numResults == 1 {
			return fmt.Sprintf(":book: Found it! %s is on p. %s of %s.", strings.Title(entries[0].Name), entries[0].Page, entries[0].Source), nil
		}
		resp := ":book: I found multiple entries for you!"
		for _, entry := range entries {
			resp += fmt.Sprintf("\n%s is on p. %s of %s", strings.Title(entry.Name), entry.Page, entry.Source)
		}
		return resp, nil
	}

	return "", &BotError{err: fmt.Sprintf("Failed to find %s-type entry for '%s'", entryType, query),
		botMsg: fmt.Sprintf(":x: Sorry! I couldn't find a %s called '%s'...", entryType, strings.Title(query))}
}

// LoadFromCSV adds data to a Lookup object
func LoadFromCSV(file os.FileInfo) []SourceEntry {
	fileName := file.Name()
	fileHandle, err := os.Open(dataDir + fileName)
	if err != nil {
		log.Printf("Failed to open data file '%s'", fileName)
		return nil
	}
	reader := csv.NewReader(fileHandle)
	data, err := reader.ReadAll()
	if err != nil {
		log.Printf("Failed to read data file '%s'", fileName)
		return nil
	}

	// trim header row
	data = data[1:][:]
	entries := make([]SourceEntry, len(data))
	for idx, row := range data {
		entries[idx] = SourceEntry{
			Name:   row[0],
			Type:   strings.TrimSuffix(fileName, ".csv"),
			Page:   row[1],
			Source: row[2],
		}
	}

	return entries
}

// returns a random entry of the specified type
func randomEntry(kind string) *SourceEntry {
	results := filter(Lookup, kind, filterType)
	source := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(source)
	return &(*results)[rnd.Intn(len(*results))]
}
