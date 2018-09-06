package pixie

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// A regex to match '(X)YdZ' or 'YdZ' inputs, where X is the number of times
// to roll all dice, Y is the number of dice to roll each time, and Z is the
// type of dice to roll (ex., Z=4 rolls a 4-sided die, etc.)
const diceRegex = `(\(\d+\)){0,1}(\d+)d(\d+)`

// parseTimesInt takes the optional first match from the regex (e.g. '(2)')
// and converts it to its int value; if it wasn't passed, or can't be parsed,
// it will default to 1 (so there will be exactly one set of rolls performed)
func parseTimesInt(times string) int {
	times = strings.TrimPrefix(times, "(")
	times = strings.TrimSuffix(times, ")")
	timesInt, _ := strconv.Atoi(times)
	if timesInt == 0 {
		// if the int can't be parsed from the string (or it didn't exist), default to 1
		timesInt = 1
	}
	return timesInt
}

// Roll parses the message input via regex, rolls the dice, and returns
// the formatted output
func Roll(input []string) (string, *BotError) {
	rx := regexp.MustCompile(diceRegex)
	matches := rx.FindStringSubmatch(input[0])
	if len(matches) != 4 {
		return "", &BotError{err: fmt.Sprintf("failed to parse roll input: %s", input[0]),
			botMsg: fmt.Sprintf("Oops! To roll, ask me like this: (X)YdZ\n  X: The number of times you want to roll (optional)\n  Y: The number of dice to roll\n  Z: The type of dice to roll (doesn't have to be a real die)")}
	}
	times := parseTimesInt(matches[1])
	dieNum, _ := strconv.Atoi(matches[2])
	dieType, _ := strconv.Atoi(matches[3])
	rolls := rollDice(times, dieNum, dieType)
	output := rollToString(rolls)
	return output, nil
}

// rollDice simulates dice rolling by generating random numbers in sets and
// sequences specified by the user's parsed input
func rollDice(times int, dieNum int, dieType int) [][]int {
	rolls := make([][]int, times)
	for idx := 0; idx < times; idx++ {
		rolls[idx] = make([]int, dieNum)
		// seed a new random source for each iteration
		source := rand.NewSource(time.Now().UnixNano())
		rnd := rand.New(source)

		for jdx := 0; jdx < dieNum; jdx++ {
			roll := rnd.Intn(dieType) + 1
			rolls[idx][jdx] = roll
		}
	}
	return rolls
}

// rollToString converts an array of rolls ([][]int) to a string for output
func rollToString(rolls [][]int) string {
	output := ":game_die: Here's your rolls, as requested."
	for _, row := range rolls {
		sum := 0
		seq := ""
		for _, num := range row {
			if num == 0 {
				continue
			}
			sum += num
			seq += strconv.Itoa(num) + ", "
		}
		output += fmt.Sprintf("\nTotal: %d [%s]", sum, strings.TrimSuffix(seq, ", "))
	}
	return output
}
