package talking_clock

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var numbers = map[int]string{
	0: "oh",

	1: "one", 2: "two", 3: "three",
	4: "four", 5: "five", 6: "six",
	7: "seven", 8: "eight", 9: "nine",

	10: "ten", 11: "eleven",
	12: "twelve", 13: "thirteen",
	14: "fourteen", 15: "fifteen",
	16: "sixteen", 17: "seventeen",
	18: "eighteen", 19: "nineteen",

	20: "twenty", 30: "thirty",
	40: "fourty", 50: "fifty",
}

var am_pm = map[int]string{
	0: "am", 1: "pm",
}

func Translate(s string) string {
	parts := strings.Split(s, ":")
	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return "Invalid hour"
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute > 59 {
		return "Invalid minute"
	}

	hour_index := hour % 12
	if hour_index == 0 {
		hour_index = 12
	}

	hour_str := numbers[hour_index]

	minute1, minute2 := minute/10, minute%10

	var minute_str string
	if minute == 0 {
		minute_str = ""
	} else if minute1 == 1 {
		minute_str = numbers[minute] + " "
	} else if minute2 == 0 {
		minute_str = numbers[minute1*10] + " "
	} else {
		minute_str = numbers[minute1*10] + " " + numbers[minute2] + " "
	}

	am_pm_str := am_pm[hour/12]

	return "It's " + hour_str + " " + minute_str + am_pm_str
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		fmt.Println(Translate(s))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "stdin reading error:", err)
	}
}
