package main

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
)

const punctuation = "`~!@#$%^&*()-_=+,./<>?;':\"[]{}\\| \r\n"

func getSeed() int64 {
	b := make([]byte, 8)
	_, err := cryptoRand.Read(b)
	if err != nil {
		panic(err)
	}

	return int64(binary.LittleEndian.Uint64(b))
}

func readFileToWords(fileName string) ([]string, error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	buf := make([]byte, 1024)
	totalRead := 0
	bytesRead, err := fd.Read(buf)

	word := make([]rune, 0)
	words := make([]string, 0)

	for err == nil {
		totalRead += bytesRead

		for i := 0; i < bytesRead; i++ {
			c := rune(buf[i])
			if strings.ContainsRune(punctuation, c) {
				if len(word) > 0 {
					words = append(words, string(word))
					word = make([]rune, 0)
				}
			} else {
				word = append(word, c)
			}
		}

		bytesRead, err = fd.Read(buf)
	}

	if err != io.EOF {
		return nil, err
	}

	if len(word) > 0 {
		words = append(words, string(word))
	}

	log.Printf("processed %d bytes from %s", totalRead, fileName)
	log.Printf("read %d words", len(words))

	return words, nil
}

func main() {
	wordList := flag.String("wordlist", "/usr/share/dict/words", "file containing words for passphrase generation")
	numPhrases := flag.Int("n", 10, "number of pass phrases to generate")
	wordsPerPhrase := flag.Int("wpp", 5, "words per passphrase")
	lower := flag.Bool("lower", true, "lower case passphrase")
	seed := flag.Int64("seed", getSeed(), "64bit custom PRNG seed, unix nano is used otherwise")
	enumerate := flag.Bool("enumerate", false, "prefix each passphrase with a number")
	flag.Parse()

	words, err := readFileToWords(*wordList)
	if err != nil {
		log.Fatalf("failed to read wordlist file: %s", err.Error())
		os.Exit(-1)
	}

	log.Printf("seed: 0x%x", uint64(*seed))

	rand.Seed(*seed)

	log.Printf("passphrases:\n\n")

	for i := 0; i < *numPhrases; i++ {
		phrase := make([]string, 0)
		for j := 0; j < *wordsPerPhrase; j++ {
			index := rand.Int63n(int64(len(words)))
			phrase = append(phrase, words[index])
		}

		l := strings.Join(phrase, " ")
		if *lower {
			l = strings.ToLower(l)
		}

		if *enumerate {
			fmt.Printf("%05d %s\n", i+1, l)
		} else {
			fmt.Printf("%s\n", l)
		}
	}
}
