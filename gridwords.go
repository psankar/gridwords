package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"unicode"
)

func getLetters(s string) []string {
	var res []string

	p := []rune(s)
	for i := 0; i < len(p); {
		j := i + 1
		for j < len(p) && (unicode.Is(unicode.Mn, p[j]) ||
			unicode.Is(unicode.Me, p[j]) || unicode.Is(unicode.Mc, p[j])) {
			j++
		}
		res = append(res, string(p[i:j]))
		i = j
	}

	if len(res) != 9 {
		log.Fatal(s, res)
	}

	return res
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

var res map[string]struct{}
var wg sync.WaitGroup

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	const usage = "Usage: gridwords <Path-to-the-words-file>\n\n"

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, usage)
		os.Exit(1)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Fprintf(os.Stderr, usage)
		return
	}

	var t, four []string

	if true {
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
				t = append(t, l)
			case 4:
				four = append(four, l)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		log.Println("Three letter words found: ", len(t))
		f, err := os.Create("tw.txt")
		if err != nil {
			log.Fatal(err)
		}
		w := bufio.NewWriter(f)
		for _, s := range t {
			_, err := w.WriteString(s + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
		w.Flush()
	} else {
		t = []string{"cat", "bat", "mat"}
	}

	res = make(map[string]struct{})
	ch := make(chan string)
	quit := make(chan bool)

	// Get matches in a go routine
	go func(ch chan string, quit chan bool) {
		var zs struct{}

		select {
		case l := <-ch:
			log.Println(l)
			res[l] = zs
		case <-quit:
			break
		}
	}(ch, quit)

	// Send matches
	for i := 0; i < len(t); i++ {

		// First word of the grid cannot contain these letters
		// as no word can begin with these letters.
		if strings.Contains(t[i], "க்ங்ச்ஞ்ட்ண்த்ந்ப்ம்ய்ர்ல்வ்ழ்ள்ற்ன்ங்ஙஙாஙிஙீஙுஙூஙெஙேஙைஙொஙோஙௌண்ணணணிணீணுணூணெணேணைணொணோணௌழ்ழழாழிழீழுழூழெழேழைழொழோழௌள்ளளாளிளீளுளூளெளேளைளொளோளௌற்றறாறிறீறுறூறெறேறைறொறோறௌன்னனானினீனுனூனெனேனைனொனோனௌ") {
			continue
		}

		for j := 0; j < len(t); j++ {
			if i == j {
				continue
			}
			for k := 0; k < len(t); k++ {

				if i == k || j == k {
					continue
				}

				l := t[i] + t[j] + t[k]

				wg.Add(1)
				go func(l string, ch chan string, i, j, k int) {

					defer wg.Done()

					g := getLetters(l)
					// log.Println(l)
					// log.Println(g, i, j, k)

					if (g[0]+g[3]+g[6] == t[i]) &&
						(g[1]+g[4]+g[7] == t[j]) &&
						(g[2]+g[5]+g[8] == t[k]) {
						ch <- l
					}
				}(l, ch, i, j, k)
			}
		}
	}
	log.Println("All goroutines created. Now waiting...")
	wg.Wait()

	quit <- true

	log.Println(len(res))
}
