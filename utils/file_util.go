package utils

import (
	"path"
	"runtime"
	"strconv"
	"strings"
)

func ParseLine(line string) (id int, xPos float64, yPos float64) {
	line = strings.TrimSpace(line)

	if len(line) == 0 {
		return -1, 0, 0
	}

	if '0' > line[0] || '9' < line[0] {
		return -1, 0, 0
	}

	words := make([]string, 0)

	for _, word := range strings.Fields(line) {
		if len(word) == 0 {
			continue
		}

		words = append(words, word)
	}

	id, _ = strconv.Atoi(words[0])
	xPos, _ = strconv.ParseFloat(words[1], 64)
	yPos, _ = strconv.ParseFloat(words[2], 64)

	return
}

func GetSourceRootPath() string {
	_, sourcePath, _, _ := runtime.Caller(0)
	sourcePath = path.Join(path.Dir(sourcePath), "..")

	return sourcePath
}
