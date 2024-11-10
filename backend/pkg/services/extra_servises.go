package services
import (
	"strings"
	"unicode"
)

// ParseCategories parses the categories from a string based on '#' delimiter
func ParseCategories(categories string) []string {
	var catArray []string
	for _, cat := range strings.Split(categories, "#") {
		trimmedCat := strings.TrimSpace(cat)
		if trimmedCat != "" {
			catArray = append(catArray, strings.ToLower(trimmedCat))
		}
	}
	return catArray
} 

func TrimAndNormalizeSpaces(input string) string {
	// Видаляємо зайві пробіли з початку та кінця рядка
	trimmed := strings.TrimSpace(input)

	// Розділяємо рядок на слова, використовуючи функцію-предикат
	words := strings.FieldsFunc(trimmed, unicode.IsSpace)

	// Збираємо слова знову, вставляючи один пробіл між ними
	result := strings.Join(words, " ")

	return result
}
