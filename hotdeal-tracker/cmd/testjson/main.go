package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

func main() {
	data, err := ioutil.ReadFile("ebay.html")
	if err != nil {
		panic(err)
	}

	html := string(data)

	// 查找所有包含 displayPrice 的行
	lines := strings.Split(html, "\n")
	for i, line := range lines {
		if strings.Contains(line, "displayPrice") && strings.Contains(line, "value") {
			fmt.Printf("Line %d: %s\n", i, line[:min(200, len(line))])
		}
	}

	// 尝试提取价格模式
	re := regexp.MustCompile(`"displayPrice":\s*\{[^}]+"value":\s*\{[^}]+"value":\s*([\d.]+)`)
	matches := re.FindAllStringSubmatch(html, 10)
	
	fmt.Printf("\nFound %d price matches\n", len(matches))
	for i, match := range matches {
		if len(match) >= 2 {
			fmt.Printf("Price %d: %s\n", i+1, match[1])
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}