package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func main() {
	data, err := ioutil.ReadFile("ebay.html")
	if err != nil {
		panic(err)
	}

	html := string(data)

	re := regexp.MustCompile(`"displayPrice":\s*\{[^}]+"value":\s*\{[^}]+"value":\s*([\d.]+),\s*"currency":\s*"([^"]+)"`)
	matches := re.FindAllStringSubmatch(html, 10)

	fmt.Printf("Found %d matches\n", len(matches))
	for i, match := range matches {
		fmt.Printf("Match %d: %s\n", i+1, match[0][:min(100, len(match[0]))])
		if len(match) >= 3 {
			fmt.Printf("  Price: %s, Currency: %s\n", match[1], match[2])
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}