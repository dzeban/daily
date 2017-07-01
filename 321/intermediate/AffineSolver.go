package main

import (
	"fmt"
	"bufio"
	"os"
	"flag"
	"log"
	"strings"
	"sync"
)

type Coeff [2]int

var wg sync.WaitGroup

func isAlpha(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
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

func affineDecode(coeff Coeff, word string) string{
	var decoded string
	for _, c := range word {
		d := int(c)
		if isAlpha(c) {
			d -= 'A'
			d = (coeff[0]*d + coeff[1]) % 26 + int('A')
		}
		decoded += string(d)
	}

	return decoded
}

func crack(word string, dict map[string]bool, coeffs chan<- Coeff) {

	// `a` coeffiecient in Affine Cipher
	// can only be coprime with 26 
	coprimes := []int{3, 5, 7, 11, 15, 17, 19, 21, 23, 25}
	for _, a := range coprimes {
		for b := 0; b < 26; b++ {
			word = strip(strings.ToUpper(word))
			decoded := affineDecode(Coeff{a, b}, word)
			_, present := dict[decoded]
			if present == true {
				// log.Println(word, decoded, d)
				coeffs <- Coeff{a, b}
			}
		}
	}

	// log.Println("Releasing crack worker")
	wg.Done()
}

func decrypt(input string, dict map[string]bool) []string {
	words := strings.Split(input, " ")

	// Create buffered channel for workers
	// to avoid deadlock because channel 
	// reading/writing is blocking like in the
	// case when we try to read from channel,
	// but nobody is writing to it.
	// Buffer size is large enough to hold the case
	// when every permutation of `a` and `b` coefficients
	// matches.
	coeffs := make(chan Coeff, len(words) * 26 * 26)
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
	// log.Println(maxCoeffs)

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
	log.Printf("Decrypting using %s dict\n", *dictFilePath)

	dict, err := loadDict(*dictFilePath)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		fmt.Println(decrypt(s, dict))
	}
}
