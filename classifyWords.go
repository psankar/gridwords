// Author: Sankar <sankar.curiosity@gmail.com>
// Distributed under Creative Commons Zero License - Public Domain
// For more information see LICENSE file
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	const usage = `Usage: classifyWords <Path-to-the-words-file>
    
The program will generate two files in $PWD named three.txt and four.txt

`

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, usage)
		os.Exit(1)
	}

	if strings.HasPrefix(os.Args[1], "-") {
		fmt.Fprintf(os.Stderr, usage)
		return
	}

	var three, four []string

	// Assuming os.Args[1] points to a
	// valid file having a list of words
	fd, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		l := scanner.Text()

		if strings.ContainsAny(l, "ஃ௦௪௫௰௯௩௬௨௲௵௴௷௱+") {
			continue
		}

		// Skip words with grandham
		grandhed := false
		grandhams := []rune{'ஜ', 'ஷ', 'ஸ', 'ஹ'}
		for _, g := range grandhams {
			if strings.ContainsRune(l, g) {
				grandhed = true
			}
		}
		if grandhed {
			continue
		}

		// Skip non-tamil words (if any)
		runes := []rune(l)
		tamil := true

		for _, r := range runes {
			if !unicode.Is(unicode.Tamil, r) {
				tamil = false
				break
			}
		}

		if !tamil {
			continue
		}

		// Skip words ending in some letters
		suffixes := []string{"க்", "ங்", "ச்", "ஞ்", "ட்", "த்", "ந்", "ப்", "வ்", "ற்"}
		hasSuffix := false
		for _, suffix := range suffixes {
			if strings.HasSuffix(l, suffix) {
				hasSuffix = true
				break
			}
		}
		if hasSuffix {
			continue
		}

		// For now, we will worry about only 3x3 and 4x4 grids
		switch strlen(l) {
		case 3:
			three = append(three, l)
		case 4:
			four = append(four, l)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	log.Println("Three letter words found: ", len(three))
	f, err := os.Create("three.txt")
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	for _, s := range three {
		_, err := w.WriteString(s + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	w.Flush()

	log.Println("Four letter words found: ", len(four))
	f, err = os.Create("four.txt")
	if err != nil {
		log.Fatal(err)
	}
	w = bufio.NewWriter(f)
	for _, s := range four {
		_, err := w.WriteString(s + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	w.Flush()
}

func strlen(s string) int {
	// fmt.Println(s)
	p := []rune(s)
	c := 0 // Count of letters (excluding diacritics)
	i := 0 // Runes iterator
	for {
		//fmt.Println(i)
		for i < len(p) && (unicode.Is(unicode.Mn, p[i]) ||
			unicode.Is(unicode.Me, p[i]) ||
			unicode.Is(unicode.Mc, p[i])) {
			i++
		}

		if i >= len(p) {
			return c
		}
		c++
		i++
	}
}
