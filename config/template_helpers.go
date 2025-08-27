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
        "stars": func(rating interface{}) template.HTML {
            var r float64
            switch v := rating.(type) {
            case int:
                r = float64(v)
            case float64:
                r = v
            default:
                r = 0
            }
            fullStars := int(r)
            if fullStars > 5 {
                fullStars = 5
            }
            if fullStars < 0 {
                fullStars = 0
            }
            stars := strings.Repeat("★", fullStars)
            emptyStars := strings.Repeat("☆", 5-fullStars)
            return template.HTML(stars + emptyStars)
        },
    }
}
