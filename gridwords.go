// Author: Sankar <sankar.curiosity@gmail.com>
// Distributed under Creative Commons Zero License - Public Domain
// For more information see LICENSE file
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
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

	return res
}

var res map[string]struct{}
var wg sync.WaitGroup
var gridWriter *bufio.Writer

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	const usage = `Usage: gridwords <Path-to-the-words-file>

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

	gridf, err := os.Create("grid.txt")
	if err != nil {
		log.Fatal(err)
	}
	gridWriter = bufio.NewWriter(gridf)

	// Assuming os.Args[1] points to a
	// valid three.txt file having a list of words
	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	t := strings.Split(string(content), "\n")
	if t[len(t)-1] == "" {
		t = t[:len(t)-1]
	}

	res = make(map[string]struct{})
	ch := make(chan string)
	quit := make(chan bool)

	// Get matches in a go routine
	go func(ch chan string, quit chan bool) {
		var zs struct{}

		for {
			select {
			case l := <-ch:
				gridWriter.WriteString(l + "\n")
				res[l] = zs
			case <-quit:
				break
			}
		}
	}(ch, quit)

	// Send matches
	for i := 0; i < len(t); i++ {

		// First word of the grid cannot contain these letters
		// as no word can begin with these letters.
		// if strings.ContainsAny(t[i], "அக்ங்ச்ஞ்ட்ண்த்ந்ப்ம்ய்ர்ல்வ்ழ்ள்ற்ன்ங்ஙஙாஙிஙீஙுஙூஙெஙேஙைஙொஙோஙௌண்ணணணிணீணுணூணெணேணைணொணோணௌழ்ழழாழிழீழுழூழெழேழைழொழோழௌள்ளளாளிளீளுளூளெளேளைளொளோளௌற்றறாறிறீறுறூறெறேறைறொறோறௌன்னனானினீனுனூனெனேனைனொனோனௌ") {
		// continue
		// }

		firstWord := getLetters(t[i])

		for j := 0; j < len(t); j++ {

			secondWord := getLetters(t[j])
			if secondWord[0] != firstWord[1] {
				continue
			}

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
	gridWriter.Flush()
}
