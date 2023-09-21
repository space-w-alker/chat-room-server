package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Field struct {
	Name      string
	FieldType string
}

type Model struct {
	Name   string
	Fields []Field
}

func main() {
	filePath := "schema.prisma"

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fileContentStr := string(content)

	r := regexp.MustCompile("model(?P<modelName>.*?){(?P<modelFields>(.|\n)*?)}")
	matches := r.FindStringSubmatch(fileContentStr)
	parseFields(matches[2])
}

func parseFields(s string) {
	s = strings.TrimSpace(s)
	sL := strings.Split(s, "\n")
	for i := range sL {
		sL[i] = strings.TrimSpace(sL[i])
	}
	fmt.Println(strings.Join(sL, ","))
}
