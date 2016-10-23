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


var res map[string]struct{}
var wg sync.WaitGroup

 else {
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
