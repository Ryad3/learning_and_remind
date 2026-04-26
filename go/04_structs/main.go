package main

import (
	"fmt"
)

type Film struct {
	Title string
	rate int
	seen bool
}

func (f Film) recommend() bool {
	return f.rate >= 8 && f.seen
}

func main() {
	film1 := Film{Title: "Inception", rate: 8, seen: true}
	film2 := Film{Title: "The Matrix", rate: 9, seen: false}
	film3 := Film{Title: "Interstellar", rate: 4, seen: true}

	films := []Film{film1, film2, film3}

	for _, film := range films {
		if film.recommend() {
			fmt.Println("I recommend the film", film.Title, "with a rate of", film.rate)
		}
	}
}