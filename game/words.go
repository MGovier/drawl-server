package game

import "fmt"

// Don't deconflict generated words per game, it could be fun to see different approaches to the same term...
func generateWord(player *Player, players []*Player) string {
	// Use a context word 10% of the time
	if len(players) >= 2 && randomInt(0, 101) < 10 {
		return generateContextWord(player, players)
	}
	// Otherwise use a random phrase
	choice := randomInt(0, len(wordList))
	return wordList[choice]
}

// Context words use other player's names.
func generateContextWord(player *Player, players []*Player) string {
	otherPlayers := make([]*Player, 0)
	for _, playa := range players {
		if player.ID != playa.ID {
			otherPlayers = append(otherPlayers, playa)
		}
	}
	chosenPlayer := otherPlayers[randomInt(0, len(otherPlayers))]
	chosenContext := contextList[randomInt(0, len(contextList))]
	return fmt.Sprintf(chosenContext, chosenPlayer.Name)
}

var contextList = []string{
	"A Portrait of %v",
	"%v's Favourite Hobby",
	"%v's Favourite Film",
	"%v's Biggest Fear",
	"%v's Favourite Animal",
	"%v's Favourite Food",
	"Something %v Would Hate",
	"%v's Idea of Hell",
	"%v in Their Happy Place",
}

var wordList = []string{
	"Hit the Nail on the Head",
	"Sausage Party",
	"One Child Policy",
	"Ring Finger",
	"Road Rage",
	"Lord of the Rings",
	"Laughing Your Head Off",
	"The Blind Leading the Blind",
	"Meat Sweats",
	"Life Giving You Lemons",
	"Freak Show",
	"Scaring Yourself Shitless",
	"1950s Sex Ed Poster",
	"Wiping the Smile Off Your Face",
	"The United Kingdom",
	"Armpit Fetish",
	"Joyless Sex",
	"Stealing Someone's Thunder",
	"Lost in IKEA",
	"Badly Drawn-On Eyebrows",
	"Crying on the Inside",
	"Brain Freeze",
	"Haunted Oven",
	"Dancing on Someone's Grave",
	"Booty Call",
	"Jazz Hands",
	"Sexually Transmitted Disease",
	"Naked Sleepwalking",
	"Butt Dialling",
	"Muffin Top",
	"Sex Bomb",
	"Tramp Stamp",
	"Crap Parkour",
	"Badly-Trained Dentist",
	"A Bizarre Gardening Accident",
	"Aztec Sacrifice",
	"Awkward Family Photo",
	"Crushing a Child's Sandcastle",
	"Giraffe Limbo Contest",
	"Walking Into a Mirror",
	"LARPing",
	"Police Line-Up",
	"Awkward Hug",
	"Staring Contest",
	"Eye-Wateringly Strong Chili",
	"Trust Exercises",
	"Frisbee Decapitation",
	"Free Willy",
	"Butt Chin",
	"Spontaneous Human Combustion",
	"Moobs",
	"Sexting",
	"Cow Tipping",
	"Cat Got Your Tongue",
	"Costing an Arm and a Leg",
	"Lost in a Ball Pit",
	"Self-Portrait (of Yourself)",
	"Face-Plant",
	"Naked Statue of Yourself",
	"Mail Order Bride",
	"Netflix and Chill",
	"Notre Dame",
	"Foam Party",
	"Up Shit Creek Without a Paddle",
	"Cereal Killer",
	"Skeleton in Your Closet",
	"At It Like Rabbits",
	"Potty Mouth",
	"Really Ugly Baby",
	"Pet Rock",
	"Unicorn Rainbow Fart",
	"Wispy Beard",
	"Stud Muffin",
	"Kicking Ducks",
	"Side Boob",
	"Rubbing Someone Up the Wrong Way",
	"Duck Face",
	"Personal Space Invasion",
	"Questionable Parenting",
	"An Unfortunate Birthmark",
	"Extreme Chafing",
	"Erotic Oven Mitts",
	"Beating Around The Bush",
	"Stealing the Covers",
	"Cheese Nightmare",
	"A Very Disappointing Christmas",
	"Drinking Problems",
	"Heavy Breathing",
	"Blaze of Glory",
	"Wheelbarrow Race",
	"Sore Loser",
	"The Pyramids",
	"Dogs Dressed as Humans",
	"Camel Toe",
	"Living in a Skip",
	"Spooning a Spoon",
	"Racist Grandparents",
	"Wet Handshake",
	"The Ugly Friend",
	"Fighting Fire With Fire",
	"Pathetic Snowman",
	"Bad Babysitter",
	"The London Eye",
	"The Australian Flag",
	"The Welsh Flag",
	"A Stereotypical Dutch Person",
	"Country Music",
	"The Beatles",
	"Stand-Up Comedian",
	"Religion",
	"The 1%",
	"Ironman",
	"The Incredible Hulk",
	"Pineapple Pizza",
	"A Haunted House",
	"A Theme Park",
	"DIY",
}
