package game

// Don't deconflict generated words per game, it could be fun to see different approaches to the same term...
func generateWord() string {
	choice := randomInt(0, len(wordList)-1)
	return wordList[choice]
}

var wordList = []string{"Hit the Nail on the Head", "Sausage Party", "One Child Policy", "Road Rage", "Lord of the Rings"}
