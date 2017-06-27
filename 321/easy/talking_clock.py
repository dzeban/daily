#!/usr/bin/env python3

digits = [
    "twelve", "one", "two",
    "three", "four", "five",
    "six", "seven", "eight",
    "nine", "ten", "eleven",
]

minutes_scores = [
    "oh", "", "twenty", "thirty", "forty", "fifty"
]

minutes_tenth = [
    "ten", "eleven", "twelve", "thirteen", "fourteen",
    "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"
]


def translate_hour(hour):
    return digits[hour % 12]


def translate_minutes(minutes):
    if minutes == 0:
        return ''

    ten, one = (minutes // 10), (minutes % 10)
    if ten == 1:
        return ' ' + minutes_tenth[one]

    return ' {}{}'.format(
                   minutes_scores[ten],
                   (' ' + digits[one]) if (one != 0) else '')


def translate(s):
    hour, minutes = map(int, s.split(':'))
    return "It's {}{}{}".format(
            translate_hour(hour),
            translate_minutes(minutes),
            ' am' if (hour // 12 == 0) else ' pm')


def main():
    while True:
        try:
            s = input()
            print(translate(s))
        except EOFError:
            return

if __name__ == '__main__':
    main()
