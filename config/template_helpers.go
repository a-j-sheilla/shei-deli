package config

import (
    "html/template"
    "strings"
)

// GetTemplateFunctions returns template helper functions
func GetTemplateFunctions() template.FuncMap {
    return template.FuncMap{
        "add": func(a, b int) int {
            return a + b
        },
        "sub": func(a, b int) int {
            return a - b
        },
        "stars": func(rating float64) template.HTML {
            fullStars := int(rating)
            stars := strings.Repeat("★", fullStars)
            emptyStars := strings.Repeat("☆", 5-fullStars)
            return template.HTML(stars + emptyStars)
        },
    }
}
