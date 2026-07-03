package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
)

func main() {
	data, err := ioutil.ReadFile("ebay.html")
	if err != nil {
		log.Fatal(err)
	}

	html := string(data)

	log.Println("Testing regex for displayPrice...")
	re := regexp.MustCompile(`"displayPrice":\s*\{"_type":"TextualDisplayValue","value":\{"value":([\d.]+)`)
	matches := re.FindAllStringSubmatch(html, 10)

	log.Printf("Found %d matches\n", len(matches))
	for i, match := range matches {
		if len(match) >= 2 {
			price, _ := strconv.ParseFloat(match[1], 64)
			log.Printf("Match %d: Price = %.2f", i+1, price)
		}
	}

	if len(matches) == 0 {
		log.Println("No matches found! Let me check the pattern...")
		// 尝试简单的模式
		re2 := regexp.MustCompile(`"value":([\d.]+),\s*"currency"`)
		matches2 := re2.FindAllStringSubmatch(html, 10)
		log.Printf("Simple pattern found %d matches", len(matches2))
		for i, match := range matches2 {
			if len(match) >= 2 {
				log.Printf("Match %d: %s", i+1, match[1])
			}
		}
	}
}