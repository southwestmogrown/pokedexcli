package cli

import "strings"

func cleanInput(text string) []string {
    words := strings.Split(text, " ")
    cleanWords := make([]string, 0, len(words))

    for _, word := range words {
        if word == "" {
            continue
        }
        trimmed := strings.TrimSpace(word)
        lower := strings.ToLower(trimmed)
        cleanWords = append(cleanWords, lower)
    }

    return cleanWords
}
