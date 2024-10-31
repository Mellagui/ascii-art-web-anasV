package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadBanner(filename string) (map[rune][]string, error) {
	if err := checkBannerLineCount(filename, 854); err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open banner file: %w", err)
	}
	defer file.Close()

	bannerMap := make(map[rune][]string)
	scanner := bufio.NewScanner(file)

	var bannerLines []string
	for scanner.Scan() {
		bannerLines = append(bannerLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading banner file: %w", err)
	}

	const (
		charHeight = 8
		startChar  = 32
	)

	for i := 0; i+charHeight <= len(bannerLines); i += charHeight + 1 {
		characterLines := bannerLines[i : i+charHeight]
		bannerMap[rune(startChar+i/(charHeight+1))] = characterLines
	}

	return bannerMap, nil
}

func checkBannerLineCount(filename string, expectedLineCount int) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open banner file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading banner file: %w", err)
	}

	if lineCount != expectedLineCount {
		return fmt.Errorf("banner file has %d lines; expected %d lines", lineCount, expectedLineCount)
	}

	return nil
}

func PrintAsciiArt(input string, bannerMap map[rune][]string, color string, sb *strings.Builder) {
	inputLines := strings.Split(input, "\n")

	for i, line := range inputLines {
		if line == "" {
			fmt.Fprintln(sb)
			continue
		}

		asciiLines := make([]string, 8)

		for _, char := range line {
			asciiArt, exists := bannerMap[char]
			if !exists {
				asciiArt = bannerMap[' ']
			}

			for j := 0; j < 8; j++ {
				asciiLines[j] += asciiArt[j]
			}
		}

		for _, line := range asciiLines {
			if line = strings.TrimRight(line, " "); line != "" {
				fmt.Fprintln(sb, line)
			}
		}

		if i < len(inputLines)-1 {
			fmt.Fprintln(sb)
		}
	}

	fmt.Fprintln(sb)
	fmt.Fprintln(sb)
}