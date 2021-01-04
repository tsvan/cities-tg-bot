package messages

import (
	"strings"
)

func cityCheck(city string, usedCities []string) bool {
	if len(usedCities) < 1 {
		return true
	}
	lastCity := strings.ToLower(usedCities[len(usedCities)-1])
	r := []rune(strings.ToLower(city))
	firstletter := string(r[0])
	if getLastLetter(lastCity) == firstletter {
		return true
	}
	return false
}

func setCitiesList(cities []string, userCity string, randomCity string) []string {
	if len(cities) > 6 {
		cities = cities[2:]
	}
	result := append(cities, userCity, randomCity)
	return result
}

func getLastLetter(word string) string {
	replacer := strings.NewReplacer("ь", "", "ъ", "", "ы", "", "й", "")
	word = replacer.Replace(word)
	r := []rune(word)
	letter := string(r[len(r)-1:])
	return letter
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
