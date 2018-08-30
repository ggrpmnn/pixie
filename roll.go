package pixie

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const diceRegex = `(\(\d+\)){0,1}(\d+)d(\d+)`

// Roll rolls some dice
func Roll(input []string) (string, *BotError) {
	rx := regexp.MustCompile(diceRegex)
	matches := rx.FindStringSubmatch(input[0])
	if len(matches) != 4 {
		return "", &BotError{err: fmt.Sprintf("failed to parse roll input: %s", input[0]),
			botMsg: fmt.Sprintf("Oops! To roll, ask me like this: (X)YdZ\n  X: The number of times you want to roll (optional)\n  Y: The number of dice to roll\n  Z: The type of dice to roll (doesn't have to be a real die!)")}
	}
	// fmt.Printf("%v\n", matches)
	times := parseTimesInt(matches[1])
	dieNum, _ := strconv.Atoi(matches[2])
	dieType, _ := strconv.Atoi(matches[3])
	// fmt.Printf("%d, %d, %d\n", times, dieNum, dieType)
	return rollDice(times, dieNum, dieType), nil
}

// rollDice rolls the dice and returns the formatted output
func rollDice(times int, dieNum int, dieType int) string {
	output := "Here's your rolls, as requested!"
	for idx := 0; idx < times; idx++ {
		// seed a new random source for each iteration
		source := rand.NewSource(time.Now().UnixNano())
		rnd := rand.New(source)

		sum := 0
		rollList := ""
		for jdx := 0; jdx < dieNum; jdx++ {
			roll := rnd.Intn(dieType) + 1
			sum += roll
			if jdx != 0 {
				rollList += ", "
			}
			rollList += strconv.Itoa(roll)
		}
		output += fmt.Sprintf("\nTotal: %d [%s]", sum, rollList)
	}
	return output
}

// parseTimesInt takes the optional first match from the regex (e.g. '(2)')
// and converts it to its int value
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
