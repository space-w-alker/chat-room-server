package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
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
	generateDatabase()
	generateGenerics()

	filePath := "schema.prisma"

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fileContentStr := string(content)
	models := parseModels(fileContentStr)
	for _, model := range models {
		generateModel(model)
	}
}

func parseModels(s string) (models []Model) {
	r := regexp.MustCompile("model(?P<modelName>.*?){(?P<modelFields>(.|\n)*?)}")
	matches := r.FindAllStringSubmatch(s, -1)
	for i := 0; i < len(matches); i++ {
		match := matches[i]
		models = append(models, Model{Name: strings.TrimSpace(match[1]), Fields: parseFields(match[2])})
	}
	return models
}

func parseFields(s string) (fields []Field) {
	s = strings.TrimSpace(s)
	sL := strings.Split(s, "\n")
	for i := range sL {
		r := regexp.MustCompile(`(\w+?)\s+?(\w+)(\??)\s*(@\S+\s+@\S+\s+)*(\/\/\S+)*`)
		sL[i] = strings.TrimSpace(sL[i])
		match := r.FindStringSubmatch(sL[i])
		ft := match[2]
		if ft == "String" || ft == "Int" || ft == "DateTime" || ft == "Boolean" {
			fields = append(fields, Field{Name: match[1], FieldType: match[2]})
		}
	}
	return fields
}

func firstCharToLower(input string) string {
	if len(input) == 0 {
		return input // Handle empty string gracefully
	}
	return strings.ToLower(input[:1]) + input[1:]
}

func firstCharToUpper(input string) string {
	if len(input) == 0 {
		return input // Handle empty string gracefully
	}
	return strings.ToUpper(input[:1]) + input[1:]
}

func camelToSnake(s string) string {
	var result bytes.Buffer

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func generateModel(model Model) {
	name := strings.ToLower(camelToSnake(firstCharToLower(model.Name)))
	filePath := fmt.Sprintf("model/%v/%v.model.go", name, name)
	controllerPath := fmt.Sprintf("model/%v/%v.controller.go", name, name)

	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		fmt.Printf("Error creating directories: %v\n", err)
		return
	}

	content, err := os.ReadFile("generator/model.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "__lowerModelName__", name)
	contentStr = strings.ReplaceAll(contentStr, "__upperModelName__", model.Name)
	contentStr = strings.ReplaceAll(contentStr, "__fields__", renderFields(model))

	if exists := checkExists(filePath); !exists {
		os.WriteFile(filePath, []byte(contentStr), os.ModePerm)
	}

	controllerContent, err := os.ReadFile("generator/controller.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	controllerContentStr := string(controllerContent)
	controllerContentStr = strings.ReplaceAll(controllerContentStr, "__lowerModelName__", name)
	controllerContentStr = strings.ReplaceAll(controllerContentStr, "__upperModelName__", model.Name)

	if exists := checkExists(controllerPath); !exists {
		os.WriteFile(controllerPath, []byte(controllerContentStr), os.ModePerm)
	}
}

func generateGenerics() {
	content, err := os.ReadFile("generator/generic.model.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	filePath := "model/generic/generic.model.go"
	exists := checkExists(filePath)
	if !exists {
		dirPath := filepath.Dir(filePath)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			fmt.Printf("Error creating directories: %v\n", err)
			return
		}
		os.WriteFile(filePath, content, os.ModePerm)
	}
}

func generateDatabase() {
	content, err := os.ReadFile("generator/database.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	filePath := "database/database.go"
	exists := checkExists(filePath)
	if !exists {
		dirPath := filepath.Dir(filePath)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			fmt.Printf("Error creating directories: %v\n", err)
			return
		}
		os.WriteFile(filePath, content, os.ModePerm)
	}
}

func renderFields(model Model) (fieldsStr string) {
	acc := []string{}
	for _, field := range model.Fields {
		acc = append(acc, renderField(field))
	}
	return strings.Join(acc, "\n")
}

func renderField(field Field) string {
	name := field.Name
	template := fmt.Sprintf(`%v %v %vjson:"%v" form:"%v" db:"%v"%v`,
		firstCharToUpper(field.Name),
		renderPrismaType(field.FieldType),
		"`",
		name,
		name,
		name,
		"`",
	)
	return template
}

func renderPrismaType(t string) string {
	switch t {
	case "Int":
		return "*int"
	case "String":
		return "*string"
	case "DateTime":
		return "*time.Time"
	case "Boolean":
		return "*bool"
	}
	panic("Invalid type")
}

func checkExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}
