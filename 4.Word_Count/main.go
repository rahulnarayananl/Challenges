package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode/utf8"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s [-lcmw] [file ...]", os.Args[0])
	}

	var flag, file string
	if len(os.Args) == 2 {
		if os.Args[1][0] == '-' {
			flag = os.Args[1]
		} else {
			file = os.Args[1]
		}
	} else {
		flag = os.Args[1]
		file = os.Args[2]
	}

	var scanner *bufio.Scanner
	var f *os.File
	var err error

	if file == "-" || file == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		f, err = os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		scanner = bufio.NewScanner(f)
	}

	if flag == "" {
		processAll(f, file)
	} else {
		processFlag(scanner, flag)
	}
}

func processAll(f *os.File, file string) {
	lineCount, wordCount, charCount := 0, 0, 0

	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lineCount++
		charCount += utf8.RuneCountInString(scanner.Text()) + 1
	}

	f.Seek(0, 0)
	scanner = bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		wordCount++
	}

	fmt.Printf("%8d %8d %8d %s\n", lineCount, wordCount, charCount, file)
}

func processFlag(scanner *bufio.Scanner, flag string) {
	if flag == "-c" {
		readChars(scanner)
	} else if flag == "-l" || flag == "-m" {
		readLines(scanner, flag)
	} else if flag == "-w" {
		readWords(scanner)
	} else {
		log.Fatalf("Invalid flag: %s. Use -c, -l, -w, or -m.", flag)
	}
}

func readChars(scanner *bufio.Scanner) {
	count := 0
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		count++
	}
	fmt.Println(count)
}

func readLines(scanner *bufio.Scanner, flag string) {
	count := 0
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if flag == "-m" {
			count += utf8.RuneCountInString(scanner.Text()) + 1
		} else {
			count++
		}
	}
	fmt.Println(count)
}

func readWords(scanner *bufio.Scanner) {
	count := 0
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		count++
	}
	fmt.Println(count)
}
