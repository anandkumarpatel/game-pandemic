package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
		Virus VirusType
		Links []string
	}

	var cityMap = []cityDef{
		{"San_Francisco", Blue, []string{"Tokyo", "Manila", "Los_Angeles", "Chicago"}},
		{"Chicago", Blue, []string{"San_Francisco", "Toronto", "Atlanta", "Mexico_City"}},
		{"Atlanta", Blue, []string{"Chicago", "Washington", "Miami"}},
		{"Toronto", Blue, []string{"Chicago", "Washington", "New_York"}},
		{"Washington", Blue, []string{"Atlanta", "Miami", "Toronto", "New_York"}},
		{"New_York", Blue, []string{"Washington", "Madrid", "Toronto", "London"}},
		{"London", Blue, []string{"New_York", "Paris", "Madrid"}},
		{"Paris", Blue, []string{"London", "Madrid", "Essen", "Milan", "Algiers"}},
		{"Madrid", Blue, []string{"New_York", "San_Paulo", "Paris", "Algiers"}},
		{"Essen", Blue, []string{"London", "Paris", "Milan", "St_Petersburg"}},
		{"St_Petersburg", Blue, []string{"Essen", "Istanbul", "Moscow"}},
		{"Milan", Blue, []string{"Essen", "Paris", "Istanbul"}},

		{"Los_Angeles", Yellow, []string{"San_Francisco", "Mexico_City", "Sydney"}},
		{"Mexico_City", Yellow, []string{"Los_Angeles", "Chicago", "Miami", "Bogota", "Lima"}},
		{"Miami", Yellow, []string{"Mexico_City", "Atlanta", "Washington", "Bogota"}},
		{"Lima", Yellow, []string{"Mexico_City", "Santiago", "Bogota"}},
		{"Santiago", Yellow, []string{"Lima"}},
		{"Bogota", Yellow, []string{"Lima", "Miami", "Mexico_City", "San_Paulo", "Buenos_Aires"}},
		{"San_Paulo", Yellow, []string{"Bogota", "Buenos_Aires", "Lagos", "Madrid"}},
		{"Buenos_Aires", Yellow, []string{"Bogota", "San_Paulo"}},
		{"Lagos", Yellow, []string{"San_Paulo", "Kinshasa", "Khartoum"}},
		{"Kinshasa", Yellow, []string{"Johannesburg", "Lagos", "Khartoum"}},
		{"Khartoum", Yellow, []string{"Johannesburg", "Lagos", "Kinshasa", "Cairo"}},
		{"Johannesburg", Yellow, []string{"Khartoum", "Kinshasa"}},

		{"Cairo", Black, []string{"Khartoum", "Algiers", "Istanbul", "Baghdad", "Riyadh"}},
		{"Algiers", Black, []string{"Cairo", "Istanbul", "Paris", "Madrid"}},
		{"Istanbul", Black, []string{"Algiers", "Milan", "St_Petersburg", "Moscow", "Baghdad", "Cairo"}},
		{"Moscow", Black, []string{"Istanbul", "St_Petersburg", "Tehran"}},
		{"Baghdad", Black, []string{"Cairo", "Istanbul", "Tehran", "Riyadh", "Karachi"}},
		{"Riyadh", Black, []string{"Baghdad", "Cairo", "Karachi"}},
		{"Tehran", Black, []string{"Moscow", "Baghdad", "Karachi", "Delhi"}},
		{"Karachi", Black, []string{"Baghdad", "Riyadh", "Tehran", "Delhi", "Mumbai"}},
		{"Delhi", Black, []string{"Karachi", "Tehran", "Kolkata", "Chennai", "Mumbai"}},
		{"Mumbai", Black, []string{"Karachi", "Delhi", "Chennai"}},
		{"Kolkata", Black, []string{"Delhi", "Bangkok", "Hong_Kong", "Chennai"}},
		{"Chennai", Black, []string{"Delhi", "Bangkok", "Mumbai", "Kolkata", "Jakarta"}},

		{"Bangkok", Red, []string{"Chennai", "Kolkata", "Hong_Kong", "Ho_Chi_Minh_City", "Jakarta"}},
		{"Hong_Kong", Red, []string{"Bangkok", "Kolkata", "Shanghai", "Ho_Chi_Minh_City", "Manila", "Taipei"}},
		{"Shanghai", Red, []string{"Hong_Kong", "Taipei", "Tokyo", "Seoul", "Beijing"}},
		{"Seoul", Red, []string{"Shanghai", "Beijing", "Tokyo"}},
		{"Beijing", Red, []string{"Seoul", "Shanghai"}},
		{"Tokyo", Red, []string{"Seoul", "Shanghai", "Osaka", "San_Francisco"}},
		{"Osaka", Red, []string{"Tokyo", "Taipei"}},
		{"Taipei", Red, []string{"Osaka", "Shanghai", "Hong_Kong", "Manila"}},
		{"Manila", Red, []string{"Taipei", "Hong_Kong", "Ho_Chi_Minh_City", "Sydney", "San_Francisco"}},
		{"Ho_Chi_Minh_City", Red, []string{"Manila", "Hong_Kong", "Bangkok", "Jakarta"}},
		{"Jakarta", Red, []string{"Ho_Chi_Minh_City", "Chennai", "Bangkok", "Sydney"}},
		{"Sydney", Red, []string{"Jakarta", "Manila", "Los_Angeles"}},
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
			VirusType: cityDef.Virus,
		})
	}

	for _, cityDef := range cityMap {
		city := cities.FindByName(cityDef.Name)
		city.Links = append(city.Links, cityDef.Links...)
	}

	cities[0].Buildings[ResearchBuilding] = true
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

func genPlayers(playerCount int, startLocation string, pDeck *Deck) Players {
	players := Players{}

	for i := 0; i < playerCount; i++ {
		player := &Player{
			Name:     fmt.Sprint("P", i),
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

const epidemicCount = 1
const playerCount = 2
const startLocation = "Atlanta"

func createState() *State {
	cities := genCities()
	decks := genDecks(cities, epidemicCount)

	players := genPlayers(playerCount, startLocation, decks.PDeck)

	setupPlayerDeck(decks.PDeck, epidemicCount)

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
		InfectionLevel: 2,
		OutbreakCount:  0,
		CurrentPlayerN: rand.Int() % playerCount,
		ActionCount:    4,
		Step:           StartStep,
		DrawCount:      0,
	}
	state.SetupVirus()

	return &state
}
func main() {
	state := &State{}
	err := state.Load()
	if err != nil {
		state = createState()
	}

	r := gin.Default()
	r.GET("/state", func(c *gin.Context) {
		c.JSON(200, state)
	})
	r.GET("/actions", func(c *gin.Context) {
		c.JSON(200, state.GetActions())
	})
	r.POST("/action/:action", func(c *gin.Context) {
		action := c.Param("action")
		if action == "" {
			c.JSON(403, "missing action")
			return
		}
		player := c.Query("player")
		if player == "" {
			c.JSON(403, "missing player")
			return
		}
		target := c.Query("target")
		if target == "" {
			c.JSON(403, "missing target")
			return
		}
		if err := state.Do(action, target); err != nil {
			c.JSON(500, err)
			return
		}
		c.JSON(200, "OK")
	})
	state.Update()
	go r.Run()
	for {
		doInput(state)
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
	fmt.Println(state)
	fmt.Println()

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
	fmt.Printf("Enter Option: ")
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
