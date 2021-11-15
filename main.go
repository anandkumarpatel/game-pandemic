package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
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

func genCities() Cities {
	type cityDef struct {
		Name  string
		Links []string
	}

	var cityMap = []cityDef{
		{"a", []string{"b", "z"}},
		{"b", []string{"a", "c"}},
		{"c", []string{"b", "d"}},
		{"d", []string{"c", "e"}},
		{"e", []string{"d", "f"}},
		{"f", []string{"e", "g"}},
		{"g", []string{"f", "h"}},
		{"h", []string{"g", "i"}},
		{"i", []string{"h", "j"}},
		{"j", []string{"i", "k"}},
		{"k", []string{"j", "l"}},
		{"l", []string{"k", "m"}},
		{"m", []string{"l", "n"}},
		{"n", []string{"m", "o"}},
		{"o", []string{"n", "p"}},
		{"p", []string{"o", "q"}},
		{"q", []string{"p", "r"}},
		{"r", []string{"q", "s"}},
		{"s", []string{"r", "t"}},
		{"t", []string{"s", "u"}},
		{"u", []string{"t", "v"}},
		{"v", []string{"u", "w"}},
		{"w", []string{"v", "x"}},
		{"x", []string{"w", "y"}},
		{"y", []string{"x", "z"}},
		{"z", []string{"y", "a"}},
	}

	cities := Cities{}
	for _, cityDef := range cityMap {

		cities = append(cities, &City{
			Name: cityDef.Name,
			VirusCounts: VirusCounts{
				Blue:   0,
				Red:    0,
				Yellow: 0,
				Black:  0,
			},
			Buildings: Buildings{
				ResearchBuilding: false,
			},
			VirusType: Black,
		})
	}

	for _, cityDef := range cityMap {
		city := cities.FindCityByName(cityDef.Name)
		city.Links = append(city.Links, cityDef.Links...)
	}

	cities[0].Buildings[ResearchBuilding] = true
	// cities[0].VirusCounts[Black] = 4
	// cities[1].VirusCounts[Black] = 3
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
			Type:      CityCardType,
			Name:      city.Name,
			City:      city,
			VirusType: Black,
		}

		decks.VDeck.AddCard(card)
		decks.PDeck.AddCard(card)
	}
	decks.VDeck.Shuffle()
	decks.PDeck.Shuffle()
	// TODO: add event cards
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

	viruses := Viruses{
		Black:  NoneVirusStatus,
		Blue:   NoneVirusStatus,
		Yellow: NoneVirusStatus,
		Red:    NoneVirusStatus,
	}

	state := State{
		Players:        players,
		Cities:         cities,
		Viruses:        viruses,
		Decks:          decks,
		InfectionLevel: 3,
		OutbreakCount:  0,
		Turn: &Turn{
			CurrentPlayer: players[firstPlayer],
			ActionCount:   4,
			Step:          StartStep,
			DrawCount:     0,
		},
	}

	state.SetupVirus()

	for {
		if state.HasWon() {
			panic("HAS WON")
		}
		fmt.Printf("start step %s\n", state.Turn.Step)
		switch state.Turn.Step {
		case StartStep:
			state.Turn.ActionCount = 4
			state.Turn.Step = ActionStep
			continue
		case ActionStep:
			if state.Turn.ActionCount < 1 {
				state.Turn.DrawCount = 2
				state.Turn.Step = DrawStep
				continue
			}
		case DrawStep:
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
			if state.Turn.InfectCount < 1 {
				state.Turn.Step = NextPlayerStep
				continue
			}
		case NextPlayerStep:
			firstPlayer = (firstPlayer + 1) % playerCount
			state.Turn.CurrentPlayer = players[firstPlayer]
			state.Turn.Step = StartStep
			continue
		default:
			panic(fmt.Errorf("invalid action: %s", state.Turn.Step))
		}

		fmt.Println(state)
		fmt.Println()

		doInput(&state)
	}
}

type DoMap [][]string

func (s DoMap) String() string {
	out := ""
	for i, do := range s {
		out += fmt.Sprintf("%d\tplayer (%s) action %s target %s\n", i, do[2], do[0], do[1])
	}
	return out
}

func (s *DoMap) Len() int {
	return len(*s)
}
func (s *DoMap) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}
func (s *DoMap) Less(i, j int) bool {
	return (*s)[i][0] < (*s)[j][0]
}

func doInput(state *State) {
	actions := state.GetActions()
	do := DoMap{}
	for player, actions := range actions {
		for action, targets := range actions {
			for _, target := range targets {
				do = append(do, []string{string(action), target, player})
			}
		}
	}

	sort.Sort(&do)
	fmt.Println(do)
	fmt.Printf("Enter Option: \n\n")
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		num, err := strconv.Atoi(strings.Trim(text, "\n"))
		if err != nil || num >= do.Len() {
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
