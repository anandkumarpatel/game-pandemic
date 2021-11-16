package main

import (
	"bufio"
	"encoding/json"
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

	var cityMap = []cityDef{}
	b, err := os.ReadFile("./cities.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &cityMap)
	if err != nil {
		panic(err)
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
			VirusType: city.VirusType,
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	state := &State{}
	err := state.Load()
	if err != nil {
		state = createState()
	}

	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/state", func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{
			"State":   state,
			"Actions": state.GetActions(),
		})
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
			fmt.Printf("ERROR %s\n", err)
			c.JSON(500, err.Error())
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
