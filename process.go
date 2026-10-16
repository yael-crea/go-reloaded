package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	markerRE          = regexp.MustCompile(`\((hex|bin|up|low|cap)(?:,\s*(\d+))?\)`)
	markerTokenRE     = regexp.MustCompile(`^\((hex|bin|up|low|cap)(?:,(\d+))?\)$`)
	spaceBeforePunct  = regexp.MustCompile(`\s+([.,!?:;]+)`)
	missingAfterPunct = regexp.MustCompile(`([.,!?:;]+)([^\s.,!?:;])`)
)

func Process(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = processLine(line)
	}
	return strings.Join(lines, "\n")
}

func processLine(line string) string {
	line = applyMarkers(line)
	line = formatPunctuation(line)
	line = formatQuotes(line)
	line = formatArticles(line)
	return line
}

func applyMarkers(line string) string {
	normalized := markerRE.ReplaceAllStringFunc(line, func(marker string) string {
		parts := markerRE.FindStringSubmatch(marker)
		normalizedMarker := "(" + parts[1]
		if parts[2] != "" {
			normalizedMarker += "," + parts[2]
		}
		return " " + normalizedMarker + ") "
	})

	words := strings.Fields(normalized)
	out := make([]string, 0, len(words))
	for _, word := range words {
		match := markerTokenRE.FindStringSubmatch(word)
		if match == nil {
			out = append(out, word)
			continue
		}

		count := 1
		if match[2] != "" {
			if n, err := strconv.Atoi(match[2]); err == nil {
				count = n
			}
		}

		for i := len(out) - 1; i >= 0 && count > 0; i-- {
			if !containsLetterOrDigit(out[i]) {
				continue
			}
			out[i] = transformWord(out[i], match[1])
			count--
		}
	}

	return strings.Join(out, " ")
}

func transformWord(word, command string) string {
	switch command {
	case "hex":
		return toDecimal(word, 16)
	case "bin":
		return toDecimal(word, 2)
	case "up":
		return strings.ToUpper(word)
	case "low":
		return strings.ToLower(word)
	case "cap":
		return capitalize(word)
	default:
		return word
	}
}

func toDecimal(word string, base int) string {
	n, err := strconv.ParseInt(word, base, 64)
	if err != nil {
		return word
	}
	return strconv.FormatInt(n, 10)
}

func capitalize(word string) string {
	runes := []rune(strings.ToLower(word))
	for i, r := range runes {
		if unicode.IsLetter(r) {
			runes[i] = unicode.ToUpper(r)
			break
		}
	}
	return string(runes)
}

func containsLetterOrDigit(word string) bool {
	for _, r := range word {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func formatPunctuation(line string) string {
	line = spaceBeforePunct.ReplaceAllString(line, "$1")
	return missingAfterPunct.ReplaceAllString(line, "$1 $2")
}

func formatQuotes(line string) string {
	tokens := strings.Fields(line)
	var result strings.Builder
	quoteOpen := false

	for _, token := range tokens {
		if token == "'" {
			if quoteOpen {
				trimTrailingSpace(&result)
				result.WriteByte('\'')
			} else {
				writeSpaceIfNeeded(&result)
				result.WriteByte('\'')
			}
			quoteOpen = !quoteOpen
			continue
		}

		if result.Len() > 0 && !(quoteOpen && strings.HasSuffix(result.String(), "'")) {
			result.WriteByte(' ')
		}
		result.WriteString(token)
	}

	return result.String()
}

func formatArticles(line string) string {
	tokens := strings.Fields(line)
	for i := 0; i+1 < len(tokens); i++ {
		if tokens[i] != "a" && tokens[i] != "A" {
			continue
		}

		next := strings.TrimLeft(tokens[i+1], "'")
		first, ok := firstLetter(next)
		if !ok || !strings.ContainsRune("aeiouh", unicode.ToLower(first)) {
			continue
		}

		if tokens[i] == "A" {
			tokens[i] = "An"
		} else {
			tokens[i] = "an"
		}
	}
	return strings.Join(tokens, " ")
}

func firstLetter(word string) (rune, bool) {
	for _, r := range word {
		if unicode.IsLetter(r) {
			return r, true
		}
	}
	return 0, false
}

func trimTrailingSpace(builder *strings.Builder) {
	value := strings.TrimRight(builder.String(), " ")
	builder.Reset()
	builder.WriteString(value)
}

func writeSpaceIfNeeded(builder *strings.Builder) {
	if builder.Len() > 0 && !strings.HasSuffix(builder.String(), " ") {
		builder.WriteByte(' ')
	}
}
