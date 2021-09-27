package util

import (
	"math/rand"
	"strconv"
	"time"
)

func FloatFromJson(key string, json map[string]interface{}) (float64, bool) {
	value, ok := json[key]
	if !ok {
		return 0.0, false
	}
	return value.(float64), true
}

func IntFromJson(key string, json map[string]interface{}) (int, bool) {
	value, ok := json[key]
	if !ok {
		return 0, false
	}
	return int(value.(float64)), true
}

func StringFromJson(key string, json map[string]interface{}) (string, bool) {
	value, ok := json[key]
	if !ok {
		return "", false
	}
	return value.(string), true
}

func FloatToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func ParseFloat(str string) (float64, error) {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	return value, err
}

func ParseInt(str string) (int64, error) {
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func RandomConnectionName() string {
	rand.Seed(time.Now().Unix())
	names := []string{
		"0xDECAFBAD", "0xBADDCAFE", "0xCAFED00D",
		"Boobytrap", "Sobeit Void", "Transpotato",
		"A Butt Tuba", "Evil Olive", "Flee to Me, Remote Elf",
		"Sit on a Potato Pan, Otis", "Taco Cat", "UFO Tofu",
	}
	return names[rand.Intn(len(names))]
}
