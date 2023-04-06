package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var characterNamePattern = regexp.MustCompile(`(?s)<a class="charName.*?">\s*(?:<span class="c"></span>\s*)?(?:.*?)<span>(?P<characterName>.*?)</span>.*?<a class="profileLink".*?>\s*<span>(?P<streamerName>.*?)</span>`)

func Handler(w http.ResponseWriter, r *http.Request) {
	var input string
	if response, err := http.Get(`https://nopixel.hasroot.com/characters.php`); err != nil {
		log.Fatalln(err)
	} else if body, err := io.ReadAll(response.Body); err != nil {
		log.Fatalln(err)
	} else {
		input = string(body)
	}

	records := make([]Record, 0)

	for _, match := range characterNamePattern.FindAllStringSubmatch(input, -1) {
		characterName := match[characterNamePattern.SubexpIndex("characterName")]
		streamerName := strings.ToLower(match[characterNamePattern.SubexpIndex("streamerName")])
		records = append(records, Record{characterName, streamerName})
	}

	result := search(records, r.URL.Query().Get("character"))

	if result.Streamer == "" {
		fmt.Fprintf(w, "%s does not match any entry in the database", r.URL.Query().Get("character"))
	} else {
		fmt.Fprintf(w, "%s is played by twitch.tv/%s", result.Character, result.Streamer)
	}
}

type Record struct {
	Character string
	Streamer  string
}

func search(items []Record, queryStr string) Record {
	var bestMatch Record
	bestMatchScore := 0

	queryWords := strings.Split(strings.ToLower(queryStr), " ")

	for _, item := range items {
		score := 0
		charTokens := strings.Split(strings.ToLower(item.Character), " ")
		for _, queryWord := range queryWords {
			for _, charToken := range charTokens {
				if strings.HasPrefix(charToken, queryWord) {
					score += len(queryWord)
					break
				}
			}
		}
		if score > bestMatchScore {
			bestMatchScore = score
			bestMatch = item
		}
	}

	return bestMatch
}
