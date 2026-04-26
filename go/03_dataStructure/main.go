package main

import (
	"fmt"
)

func main() {
	films := []string{"Inception", "The Matrix", "Interstellar", "Labyrinth", "The Lord of the Rings"}

	fmt.Println("My favorite films are:")
	
	for _, film := range films {
		fmt.Println("-", film)
	}

	notesFilms := map[string]int{
		"Inception": 8,
		"The Matrix": 9,
		"Interstellar": 4,
		"Labyrinth": 7,
		"The Lord of the Rings": 9,
	}

	for film, note := range notesFilms {
		if note >= 8 {
			fmt.Println("I really liked the film", film, "with a note of", note)
		}
	}
}