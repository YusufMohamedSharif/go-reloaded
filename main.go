package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func isFileEmpty(filename string) bool {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		// Handle error
		return false
	}
	defer file.Close()

	// Create a new empty buffer
	var buffer bytes.Buffer

	// Copy the file contents to the buffer
	_, err = buffer.ReadFrom(file)
	if err != nil {
		// Handle error
		return false
	}

	// Check if the buffer is empty
	if buffer.Len() == 0 {
		return true
	}

	return false
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Wrong number of arguments!")
		return
	}

	if isFileEmpty(os.Args[1]) {
		fmt.Println("The file is empty")
		return
	} else {
		// fmt.Println("The file is not empty")
	}

	args := os.Args[1:]

	// reading first file

	givenText, _ := os.ReadFile(args[0])

	arrayFormatting := strings.Fields(string(givenText))
	formattedString := strings.Join(arrayFormatting, " ")

	// array to push the words into
	words := strings.Split(formattedString, " ")
	for i, word := range words {
		if word == "(up)" && i != 0 {
			words[i-1] = strings.ToUpper(words[i-1])
			words = append(words[:i], words[i+1:]...)
		} else if word == "(low)" && i != 0 {
			words[i-1] = strings.ToLower(words[i-1])
			words = append(words[:i], words[i+1:]...)
		} else if word == "(cap)" && i != 0 {
			if len(words[i-1]) > 0 {
				words[i-1] = strings.ToUpper(string(words[i-1][0])) + strings.ToLower(string(words[i-1][1:]))
			}
			words = append(words[:i], words[i+1:]...)
		} else if word == "(hex)" && i != 0 {
			words[i-1] = HextoInt(words[i-1])
			words = append(words[:i], words[i+1:]...)
		} else if word == "(bin)" && i != 0 {
			words[i-1] = BintoInt(words[i-1])
			words = append(words[:i], words[i+1:]...)

			// upper with number
		} else if word == "(up," && i != 0 {
			b := strings.Trim(string(words[i+1]), words[i+1][len(words[i+1])-1:])
			number, _ := strconv.Atoi(string(b))
			for j := 1; j <= number; j++ {
				if i-j >= 0 {
					words[i-j] = strings.ToUpper(words[i-j])
				}
			}
			words = append(words[:i], words[i+2:]...)
			// lower with number
		} else if word == "(low," && i != 0 {
			b := strings.Trim(string(words[i+1]), words[i+1][len(words[i+1])-1:])
			number, _ := strconv.Atoi(string(b))
			for j := 1; j <= number; j++ {
				if i-j >= 0 {
					words[i-j] = strings.ToLower(words[i-j])
				}
			}
			words = append(words[:i], words[i+2:]...)
			// capitalize with num
		} else if word == "(cap," && i != 0 {
			b := strings.Trim(string(words[i+1]), words[i+1][len(words[i+1])-1:])
			number, _ := strconv.Atoi(string(b))
			for j := 1; j <= number; j++ {
				if i-j >= 0 {
					words[i-j] = strings.ToUpper(string(words[i-j][0])) + strings.ToLower(string(words[i-j][1:]))
				}
			}
			words = append(words[:i], words[i+2:]...)
		}
	}

	ChangeA(words)

	// join slice
	needed := strings.Join(Punctuations(words), " ")

	// write file, automatically updates manipulated file.
	man := os.WriteFile(args[1], []byte(needed), 0o644)
	if man != nil {
		panic(man)
	}
}

// conv hex to int
func HextoInt(hex string) string {
	number, _ := strconv.ParseInt(hex, 16, 64)
	return fmt.Sprint(number)
}

// conv binary to int
func BintoInt(bin string) string {
	number, _ := strconv.ParseInt(bin, 2, 64)
	return fmt.Sprint(number)
}

func ChangeA(s []string) []string {
	vowels := []string{"a", "e", "i", "o", "u", "h", "A", "E", "I", "O", "U", "H"}

	for i, word := range s {
		for _, letter := range vowels {
			if i+1 < len(s) && len(s[i+1]) > 0 && word == "a" && string(s[i+1][0]) == letter {
				s[i] = "an"
			} else if word == "A" && string(s[i+1][0]) == letter {
				s[i] = "An"
			}
		}
	}
	return s
}

func Punctuations(s []string) []string {
	puncs := []string{",", ".", "!", "?", ":", ";"}
	// punc in the middle of a string connecting to word after
	for i, word := range s {
		for _, punc := range puncs {
			if len(word) > 0 && i != 0 && string(word[0]) == punc && string(word[len(word)-1]) != punc {
				s[i-1] += punc
				s[i] = word[1:]
			}
		}
	}

	// punc at end of string
	for i, word := range s {
		for _, punc := range puncs {
			if len(s) > 1 && i != 0 && len(word) > 0 && i < len(s) && (string(word[0]) == punc) && (s[len(s)-1] == s[i]) {
				s[i-1] += word
				s = s[:len(s)-1]
			}
		}
	}

	// punc in middle of string
	for i, word := range s {
		for _, punc := range puncs {
			if len(word) > 0 && len(s) > i && i != 0 && string(word[0]) == punc && string(word[len(word)-1]) == punc && s[i] != s[len(s)-1] {
				s[i-1] += word
				s = append(s[:i], s[i+1:]...)
			}
		}
	}

	count := 0
	for i, word := range s {
		if word == "'" {
			if i+1 < len(s) && count%2 == 0 {

				s[i+1] = word + s[i+1]
				s = append(s[:i], s[i+1:]...)

			} else if i-1 < len(s) && i != 0 && count%2 == 1 {
				if i+1 < len(s) {
					s[i-1] = s[i-1] + word
					s = append(s[:i], s[i+1:]...)
				}
			}

			if i != 0 && i == len(s)-1 && s[len(s)-1] == "'" && s[len(s)-2] != " " {
				s[len(s)-2] = s[len(s)-2] + s[len(s)-1]
				s = append(s[:i], s[i+1:]...)
			}

			count += 1
		}
	}

	stuckedQuotes := regexp.MustCompile("('')")

	addedSpace := stuckedQuotes.ReplaceAllString(strings.Join(s, " "), "$1 ")

	// // for apostrophe
	// count := true
	// for i, word := range s {
	// 	if word == "'" {
	// 		if i+1 < len(s) {
	// 			if count {
	// 				s[i+1] = word + s[i+1]
	// 				s = append(s[:i], s[i+1:]...)
	// 			}
	// 			count = !count

	// 		}
	// 	}
	// }
	//  for second apostrophe
	// for i, word := range s {
	// 	if word == "'" {
	// 		if i-1 < len(s) && i != 0 {
	// 			s[i-1] = s[i-1] + word
	// 			s = append(s[:i], s[i+1:]...)
	// 		}
	// 	}
	// }
	return strings.Fields(addedSpace)
}
