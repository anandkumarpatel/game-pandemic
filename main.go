package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Decks struct {
	VDeck    *Deck
	VDiscard *Deck
	PDeck    *Deck
	PDiscard *Deck
}

func setupPlayerDeck(deck *Deck, epidemicCount int) {
	newDecks := make([]Deck, epidemicCount)
	deck.Shuffle()
	for i, c := range deck.Cards {
		newDecks[i%epidemicCount].AddCard(c)
	}
	deck.Clear()
	for i := 0; i < epidemicCount; i++ {
		newDecks[i].AddCard(&Card{
			Name: "epidemic",
			Type: PandemicCardType,
		})
		newDecks[i].Shuffle()
		deck.AddDeck(&newDecks[i])
	}

}

func setupVirusDeck(deck *Deck) {
	deck.Shuffle()
}

func genCities() Cities {
	type cityDef struct {
		Name  string
		Links []string
	}

	var cityMap = []cityDef{
		{"a", []string{"b"}},
		{"b", []string{"a"}},
		{"c", []string{}},
		{"d", []string{}},
		{"e", []string{}},
		{"f", []string{}},
		{"g", []string{}},
		{"h", []string{}},
		{"i", []string{}},
		{"j", []string{}},
		{"k", []string{}},
		{"l", []string{}},
		{"m", []string{}},
		{"n", []string{}},
		{"o", []string{}},
		{"p", []string{}},
		{"q", []string{}},
		{"r", []string{}},
		{"s", []string{}},
		{"t", []string{}},
		{"u", []string{}},
		{"v", []string{}},
		{"w", []string{}},
		{"x", []string{}},
		{"y", []string{}},
		{"z", []string{}},
	}

	cities := Cities{}
	for _, cityDef := range cityMap {
		cities = append(cities, &City{
			Name: cityDef.Name,
		})
	}

	for _, cityDef := range cityMap {
		city := cities.FindCityByName(cityDef.Name)
		city.Links = append(city.Links, cityDef.Links...)
	}
	return cities
}

func genDecks(cities Cities, epidemicCount int) Decks {
	decks := Decks{
		VDeck:    &Deck{},
		VDiscard: &Deck{},
		PDeck:    &Deck{},
		PDiscard: &Deck{},
	}

	for _, city := range cities {
		card := &Card{
			Type: CityCardType,
			Name: city.Name,
			City: city,
		}

		decks.VDeck.AddCard(card)
		decks.PDeck.AddCard(card)
	}
	// TODO: add event cards
	setupVirusDeck(decks.PDeck)
	return decks
}

func genPlayers(playerCount int, startLocation string, pDeck *Deck) []*Player {
	players := []*Player{}

	for i := 0; i < playerCount; i++ {
		player := &Player{
			Name:     fmt.Sprint(i),
			Location: startLocation,
		}
		players = append(players, player)
		for i := 0; i < 6-playerCount; i++ {
			player.Hand.AddCard(pDeck.Draw())
		}
	}
	return players
}

type Step string

const (
	StartStep      Step = "StartStep"
	ActionStep     Step = "ActionStep"
	DrawStep       Step = "DrawStep"
	DiscardStep    Step = "DiscardStep"
	InfectionStep  Step = "InfectionStep"
	NextPlayerStep Step = "NextPlayerStep"
)

func main() {
	const epidemicCount = 5
	const playerCount = 2
	const startLocation = "a"

	cities := genCities()
	decks := genDecks(cities, epidemicCount)

	players := genPlayers(playerCount, startLocation, decks.PDeck)

	setupPlayerDeck(decks.PDeck, epidemicCount)
	firstPlayer := rand.Int() % playerCount

	viruses := Viruses{}

	state := State{
		Players:        players,
		Cities:         cities,
		Virus:          viruses,
		Decks:          decks,
		InfectionLevel: 2,
		Turn: &Turn{
			CurrentPlayer: players[firstPlayer],
			ActionCount:   4,
			Step:          StartStep,
			DrawCount:     0,
		},
	}

	for {
		fmt.Printf("start step %s\n", state.Turn.Step)
		switch state.Turn.Step {
		case StartStep:
			state.Turn.ActionCount = 4
			state.Turn.Step = ActionStep
			continue
		case ActionStep:
			fmt.Printf("XX Count %d, True? %b", state.Turn.ActionCount, state.Turn.ActionCount < 1)
			if state.Turn.ActionCount < 1 {
				state.Turn.DrawCount = 2
				state.Turn.Step = DrawStep
				fmt.Println("XX here")
				continue
			}
		case DrawStep:
			card := decks.PDeck.Draw()
			state.Turn.CurrentPlayer.Hand.AddCard(card)
			state.Turn.DrawCount--
			// TODO pandemic card
			if state.Turn.DrawCount < 1 {
				state.Turn.Step = DiscardStep
				continue
			}
		case DiscardStep:
			if state.Turn.CurrentPlayer.Hand.Count() < 8 {
				state.Turn.InfectCount = state.InfectionLevel
				state.Turn.Step = InfectionStep
				continue
			}
		case InfectionStep:
			state.Turn.InfectCount--
			if state.Turn.InfectCount < 1 {
				state.Turn.Step = NextPlayerStep
				continue
			}
			if err := state.Do("infect", ""); err != nil {
				panic(err)
			}
		case NextPlayerStep:
			firstPlayer = (firstPlayer + 1) % playerCount
			state.Turn.CurrentPlayer = players[firstPlayer]
			state.Turn.Step = StartStep
		default:
			panic(fmt.Errorf("invalid action: %d", state.Turn.Step))
		}

		fmt.Println(state)
		fmt.Println()

		doInput(&state)
	}
}

func doInput(state *State) {
	actions := state.GetActions()
	count := 0
	do := [][]string{}
	for player, actions := range actions {
		for action, targets := range actions {
			for _, target := range targets {
				fmt.Printf("%d\tplayer (%s) action %s target %s\n", count, player, action, target)
				do = append(do, []string{string(action), target})
				count++
			}
		}
	}
	fmt.Printf("Enter Option: \n\n")
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		num, err := strconv.Atoi(strings.Trim(text, "\n"))
		if err != nil {
			fmt.Printf("invalid input  (%s,%d) : try again \n", text, num)
			continue
		}
		command := do[num]
		fmt.Printf("\n Running command %s with %s\n\n", command[0], command[1])
		if err := state.Do(command[0], command[1]); err != nil {
			panic(err)
		}
		break
	}
}
