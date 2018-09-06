package pixie

import (
	"fmt"
	"sort"
)

// Generate gathers selections for each character category (race, background,
// and class) and helps create a character
func Generate(input []string) (string, *BotError) {
	raceVal := randomEntry(race)
	bgVal := randomEntry(bg)
	classVal := randomEntry(class)

	str, dex, con, intl, wis, cha := genStats()

	raceArticle := ""
	if startsWithVowel(raceVal.Name) {
		raceArticle = "an"
	} else {
		raceArticle = "a"
	}
	bgArticle := ""
	if startsWithVowel(bgVal.Name) {
		bgArticle = "an"
	} else {
		bgArticle = "a"
	}

	return fmt.Sprintf(":tada: You are %s %s (p. %s, %s) %s (p. %s, %s), who is (or was) %s %s (p. %s, %s).\nYour vitals (not including racial bonuses/penalties):\n%s%s%s%s%s%s",
		raceArticle,
		raceVal.Name,
		raceVal.Page,
		raceVal.Source,
		classVal.Name,
		classVal.Page,
		classVal.Source,
		bgArticle,
		bgVal.Name,
		bgVal.Page,
		bgVal.Source,
		str,
		dex,
		con,
		intl,
		wis,
		cha), nil
}

// generate the character stats in one bunch
func genStats() (str, dex, con, intl, wis, cha string) {
	rolls := rollDice(6, 4, 6)
	str = "Strength: " + removeLowestInRow(rolls[0])
	dex = "Dexterity: " + removeLowestInRow(rolls[1])
	con = "Constitution: " + removeLowestInRow(rolls[2])
	intl = "Intelligence: " + removeLowestInRow(rolls[3])
	wis = "Wisdom: " + removeLowestInRow(rolls[4])
	cha = "Charisma: " + removeLowestInRow(rolls[5])
	return
}

// returns the output string indicating the dropped number
func removeLowestInRow(row []int) string {
	output := ""
	sum := 0
	sort.Ints(row)
	reverse(row)
	for idx, num := range row {
		if idx+1 != len(row) {
			sum += num
			output += fmt.Sprintf("%d, ", num)
		} else {
			output += fmt.Sprintf("~~%d~~", num)
		}
	}
	return fmt.Sprintf("%d: [%s]\n", sum, output)
}

// reverses a row's contents
func reverse(row []int) {
	for idx, jdx := 0, len(row)-1; idx < jdx; idx, jdx = idx+1, jdx-1 {
		row[idx], row[jdx] = row[jdx], row[idx]
	}
}

func startsWithVowel(str string) bool {
	ch := str[0]
	return ch == 'a' || ch == 'e' || ch == 'i' || ch == 'o' || ch == 'u'
}
