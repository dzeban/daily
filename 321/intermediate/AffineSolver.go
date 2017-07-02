package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Coeff [2]int

var wg sync.WaitGroup

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func isAlpha(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func isLowercase(r rune) bool {
	return (r >= 'a' && r <= 'z')
}

func strip(s string) string {
	var stripped string

	for _, r := range s {
		if isAlpha(r) {
			stripped += string(r)
		}
	}
	return stripped
}

func affineDecode(coeff Coeff, word string) string {
	var decoded string
	for _, c := range word {
		d := int(c)
		if isAlpha(c) {
			var baseChar rune
			if isLowercase(c) {
				baseChar = 'a'
			} else {
				baseChar = 'A'
			}

			d -= int(baseChar)
			d = (coeff[0]*d+coeff[1])%26 + int(baseChar)
		}
		decoded += string(d)
	}

	return decoded
}

func crack(word string, dict map[string]bool, coeffs chan<- Coeff) {
	// Find Affine Cipher coefficients by brute forcing `a` and `b`
	// coefficients with dict lookup

	// `a` coeffiecient in Affine Cipher
	// can only be coprime with 26
	coprimes := []int{3, 5, 7, 11, 15, 17, 19, 21, 23, 25}
	for _, a := range coprimes {
		// Because of mod 26 operation,
		// `b` can only be 0..25
		for b := 0; b < 26; b++ {
			// We crack only stripped uppercase words
			word = strip(strings.ToUpper(word))

			// Try to decode current combination
			// and match with dict
			decoded := affineDecode(Coeff{a, b}, word)
			_, present := dict[decoded]
			if present == true {
				coeffs <- Coeff{a, b}
			}
		}
	}

	wg.Done()
}

func decrypt(input string, dict map[string]bool) []string {
	defer timeTrack(time.Now(), "decrypt")
	words := strings.Split(input, " ")

	// Create buffered channel for workers
	// to avoid deadlock because channel
	// reading/writing is blocking like in the
	// case when we try to read from channel,
	// but nobody is writing to it.
	// Buffer size is large enough to hold the case
	// when every permutation of `a` and `b` coefficients
	// is matching.
	coeffs := make(chan Coeff, len(words)*26*26)
	for _, word := range words {
		wg.Add(1)
		go crack(word, dict, coeffs)
	}
	// Wait for child goroutines to finish
	wg.Wait()

	// Close the channel to read it
	// with `range` and avoid blocking
	close(coeffs)

	table := make(map[Coeff]int)
	for v := range coeffs {
		table[v] += 1
	}

	// We have brute forced coefficients for every word.
	// Find the most frequent matching coefficient
	maxCoeffs := make([]Coeff, 1)
	maxCount := 0
	for k, v := range table {
		if v == maxCount && len(maxCoeffs) >= 1 {
			maxCoeffs = append(maxCoeffs, k)
		}

		if v > maxCount {
			maxCoeffs = make([]Coeff, 1)
			maxCoeffs[0] = k
			maxCount = v
		}
	}

	// Now we have affine cipher coefficients, so
	// decrypt the original message
	output := make([]string, 0)
	for _, coeff := range maxCoeffs {
		output = append(output, affineDecode(coeff, input))
	}

	return output
}

func loadDict(path string) (map[string]bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dict := make(map[string]bool)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		dict[strings.ToUpper(t)] = true
	}

	return dict, nil
}

func main() {
	dictFilePath := flag.String("dict", "/usr/share/dict/words", "Path to dict file")
	flag.Parse()

	dict, err := loadDict(*dictFilePath)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		for _, d := range decrypt(s, dict) {
			fmt.Println(s)
			fmt.Println(d)
		}
	}
}
