package chapter04

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func Runner() {
	fmt.Println("*** Chapter 04 ***")
	ReverseArray()
	DeleteDuplications()
	RemoveDuplicatedSpaces()
	//CountElements()
	//CountWords()
	SearchWebComic()
}

func ReverseArray() {
	fmt.Println("*** Running ReverseArray ***")

	//!+array
	a := [...]int{0, 1, 2, 3, 4, 5}
	reverse(&a)
	fmt.Println(a) // "[5 4 3 2 1 0]"
	//!-array
}

// !+rev
// reverse reverses a slice of ints in place.
func reverse(s *[6]int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

//!-rev

func DeleteDuplications() {
	fmt.Println("*** Running DeleteDuplications ***")

	//!+array
	a := []string{"test", "test", "this is a test"}
	removeDuplications(a)
	//!-array
}

// removeDuplications remove adjacent duplicated strings.
func removeDuplications(s []string) {
	for i := 0; i < len(s)-1; i++ {
		nextValue := s[i+1]
		if s[i] == nextValue {
			s = removeElementFromSlice(s, i)
			// once the elements was changed, must return to the previous
			i--
		}
	}
	fmt.Println("result:", strings.Join(s, ", "))
}

// removeElementFromSlice removes the i element from a slice.
func removeElementFromSlice(slice []string, i int) []string {
	// copy will move the elements to the left (rewriting the i element)
	copy(slice[i:], slice[i+1:])
	// removes the last element (is duplicated after move elements)
	return slice[:len(slice)-1]
}

func RemoveDuplicatedSpaces() {
	fmt.Println("*** Running RemoveDuplicatedSpaces ***")

	b := []byte("test  test this is a test")
	result := b[:0]
	for i := 0; i < len(b); i++ {
		if unicode.IsSpace(rune(b[i])) && unicode.IsSpace(rune(b[i+1])) {
			continue
		} else {
			result = append(result, b[i])
		}
	}

	fmt.Println("result:", string(result))
}

func CountElements() {
	fmt.Println("*** Running CountElements ***")

	counts := make(map[rune]int) // counts of Unicode characters
	countLetters := 0
	countDigits := 0
	countSpaces := 0
	var utflen [utf8.UTFMax + 1]int // count of lengths of UTF-8 encodings
	invalid := 0                    // count of invalid UTF-8 characters

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		switch {
		case unicode.IsLetter(r):
			countLetters++
		case unicode.IsDigit(r):
			countDigits++
		case unicode.IsSpace(r):
			countSpaces++
		}

		counts[r]++
		utflen[n]++
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}

func CountWords() {
	fmt.Println("*** Running CountWords ***")

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Errorf("error: %s", err.Error())
		return
	}

	result := readFile(file)
	fmt.Println("result:", result)
}

func readFile(file *os.File) map[string]int {
	wordFreq := make(map[string]int)
	input := bufio.NewScanner(file)

	input.Split(bufio.ScanWords)
	for input.Scan() {
		wordFreq[input.Text()]++
	}
	err := input.Err()
	if err != nil {
		fmt.Errorf("error: %s", err.Error())
	}
	return wordFreq
}

func SearchWebComic() {
	fmt.Println("*** Running SearchWebComic ***")

	fmt.Println("Choose a comic index...")
	var input string
	fmt.Scan(&input)

	url := fmt.Sprintf("https://xkcd.com/%s/info.0.json", input)
	result, err := xkcd(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for key, value := range result {
		fmt.Println(key, value)
	}
}

func xkcd(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("got error: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got unexpected status: %s", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("got error while read body: %s", err.Error())
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("got error while unmarshal body: %s", err.Error())
	}

	return result, nil
}
