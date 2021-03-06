package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type State struct {
	Players        Players
	Cities         Cities
	Viruses        Viruses
	Decks          Decks
	InfectionLevel int
	OutbreakCount  int
	OutbreakCities Cities
	CurrentPlayerN int
	ActionCount    int
	Step           Step
	DrawCount      int
	InfectCount    int
}

func (s State) String() string {
	out := ""
	out += fmt.Sprintf("Cities %s\n", s.Cities)
	out += fmt.Sprintf("Viruses %v\n", s.Viruses)
	out += fmt.Sprintf("PDeck %s\n", s.Decks.PDeck)
	out += fmt.Sprintf("VDeck %s\n", s.Decks.VDeck)
	out += fmt.Sprintf("PDiscard %s\n", s.Decks.PDiscard)
	out += fmt.Sprintf("VDiscard %s\n", s.Decks.VDiscard)
	out += fmt.Sprintf("InfectCount %d ", s.InfectCount)
	out += fmt.Sprintf("InfectionLevel %d ", s.InfectionLevel)
	out += fmt.Sprintf("ActionCount %d ", s.ActionCount)
	out += fmt.Sprintf("DrawCount %d ", s.DrawCount)
	out += fmt.Sprintf("OutbreakCount %d\n", s.OutbreakCount)
	for i, p := range s.Players {
		out += fmt.Sprintf("Player(%d) %s\n", i, p)
	}

	return out
}

type Action string

const (
	MoveAction        Action = "move"
	FlyAnywhereAction Action = "flyAnywhere"
	FlyToAction       Action = "flyTo"
	CureAction        Action = "cure"
	BuildAction       Action = "build"
	DiscardAction     Action = "discard"
	InfectAction      Action = "infect"
	DrawAction        Action = "draw"
	ResearchAction    Action = "research"
	OutbreakAction    Action = "outbreak"
	EpidemicAction    Action = "epidemic"
	GiveCardAction    Action = "give"
	GetCardAction     Action = "get"
)

func (a Action) String() string {
	return string(a)
}

func (s *State) Do(action string, target string) error {
	switch Action(action) {
	case MoveAction:
		s.MoveAction(target)
	case FlyToAction:
		s.FlyToAction(target)
	case FlyAnywhereAction:
		s.FlyAnywhereAction(target)
	case DiscardAction:
		s.DiscardAction(target)
	case DrawAction:
		s.DrawAction()
	case InfectAction:
		s.InfectAction(1)
	case CureAction:
		s.CureAction(target)
	case BuildAction:
		s.BuildAction(target)
	case ResearchAction:
		return s.ResearchAction(target)
	case OutbreakAction:
		s.OutbreakAction(target)
	case EpidemicAction:
		s.EpidemicAction()
	case GiveCardAction:
		s.GiveCardAction(target)
	case GetCardAction:
		s.GetCardAction(target)
	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	s.Update()

	return nil
}

func (s *State) SetupVirus() {
	for i := 3; i > 0; i-- {
		for j := 0; j < 3; j++ {
			s.InfectAction(i)
		}
	}
}

func (s State) CurrentPlayer() *Player {
	return s.Players[s.CurrentPlayerN]
}

func (s *State) ResearchAction(target string) error {
	// First split is always virus
	split := strings.Split(target, ":")
	virusName := split[0]
	cureCount := s.CurrentPlayer().CureCount
	fmt.Printf("Cure counts %d, %t, %d, %d", cureCount, (len(split) < cureCount+1), len(split), cureCount+1)
	if len(split) < cureCount+1 {
		return fmt.Errorf("require %d cards for research", cureCount)
	}
	for i := 1; i < cureCount; i++ {
		card := s.CurrentPlayer().Hand.RemoveCard(split[i])
		s.Decks.PDiscard.AddCard(card)
	}
	v := VirusType(virusName)
	s.Viruses[v] = CuredVirusStatus

	if s.Cities.IsEradicated(v) {
		s.Viruses[v] = EradicatedVirusStatus
	}
	s.ActionCount--
	return nil
}

func (s *State) MoveAction(cityName string) {
	s.CurrentPlayer().Move(cityName)
	s.ActionCount--
}

func (s *State) GiveCardAction(target string) {
	split := strings.Split(target, ":")
	cardName, playerName := split[0], split[1]
	card := s.CurrentPlayer().Hand.RemoveCard(cardName)
	toPlayer := s.Players.FindByName(playerName)
	toPlayer.Hand.AddCard(card)
	s.ActionCount--
}

func (s *State) GetCardAction(target string) {
	split := strings.Split(target, ":")
	cardName, playerName := split[0], split[1]
	fromPlayer := s.Players.FindByName(playerName)
	card := fromPlayer.Hand.RemoveCard(cardName)
	s.CurrentPlayer().Hand.AddCard(card)
	s.ActionCount--
}

func (s *State) EpidemicAction() {
	s.InfectionLevel++
	s.CurrentPlayer().Hand.RemoveCard("epidemic")

	card := s.Decks.VDeck.BackDraw()
	city := s.Cities.FindByName(card.Name)
	if s.Viruses[city.VirusType] != EradicatedVirusStatus {
		city.VirusCounts[city.VirusType] += 3
	}
	s.Decks.VDiscard.AddCard(card)
	s.Decks.VDiscard.Shuffle()

	s.Decks.VDeck.AddDeck(s.Decks.VDiscard)
	s.Decks.VDiscard.Clear()
}

func (s *State) FlyAnywhereAction(cityName string) {
	card := s.CurrentPlayer().Hand.RemoveCard(s.CurrentPlayer().Location)
	s.Decks.PDiscard.AddCard(card)
	s.CurrentPlayer().Move(cityName)
	s.ActionCount--
}

func (s *State) FlyToAction(cityName string) {
	card := s.CurrentPlayer().Hand.RemoveCard(cityName)
	s.Decks.PDiscard.AddCard(card)
	s.CurrentPlayer().Move(cityName)
	s.ActionCount--
}

func (s *State) DiscardAction(cityName string) {
	card := s.CurrentPlayer().Hand.RemoveCard(cityName)
	s.Decks.PDiscard.AddCard(card)
}

func (s *State) CureAction(target string) {
	split := strings.Split(target, ":")
	cityName, virusName := split[0], split[1]
	v := VirusType(virusName)

	city := s.Cities.FindByName(cityName)
	city.VirusCounts[v]--
	if s.Viruses[v] == CuredVirusStatus {
		city.VirusCounts[v] = 0
	}

	if s.Cities.IsEradicated(v) {
		s.Viruses[v] = EradicatedVirusStatus
	}
	s.ActionCount--
}

func (s *State) OutbreakAction(target string) {
	s.OutbreakCount++
	city := s.Cities.FindByName(target)
	virus := city.GetOutbreak()
	city.VirusCounts[virus] = 3
	for _, name := range city.Links {
		if s.OutbreakCities.Contains(name) {
			continue
		}
		nCity := s.Cities.FindByName(name)
		nCity.VirusCounts[virus]++
	}
	s.OutbreakCities = append(s.OutbreakCities, city)
}

func (s State) HasWon() {
	// TODO check # virus
	if s.Viruses.AllCured() {
		panic("HAS WON")
	}

	if s.Decks.VDeck.Count() < 1 {
		panic("LOST: No more Virus Cards")
	}
	if s.Decks.PDeck.Count() < 1 {
		panic("LOST: No more Player Cards")
	}
	if s.OutbreakCount > 8 {
		panic("LOST: To many outbreak")
	}
}

func (s *State) Save() {
	b, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("./here.json", b, 0666)
}

func (s *State) Load() error {
	b, err := os.ReadFile("./here.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, s)
	if err != nil {
		return err
	}

	return nil
}

func (s *State) BuildAction(target string) {
	split := strings.Split(target, ":")
	cityName, buildName := split[0], split[1]
	card := s.CurrentPlayer().Hand.RemoveCard(cityName)
	s.Decks.PDiscard.AddCard(card)
	s.Cities.FindByName(cityName).Buildings[Building(buildName)] = true
	s.ActionCount--
}

func (s *State) DrawAction() {
	card := s.Decks.PDeck.Draw()
	s.CurrentPlayer().Hand.AddCard(card)
	s.DrawCount--
}

func (s *State) InfectAction(count int) {
	card := s.Decks.VDeck.Draw()
	city := s.Cities.FindByName(card.Name)
	if s.Viruses[city.VirusType] != EradicatedVirusStatus {
		city.VirusCounts[city.VirusType] += count
	}
	s.Decks.VDiscard.AddCard(card)
	s.InfectCount--
}

type ActionMap map[Action]([]string)
type PlayersActions map[string]ActionMap

func (s *State) GetActions() PlayersActions {
	out := PlayersActions{}
	for _, player := range s.Players {
		out[player.Name] = ActionMap{}
	}
	hasOutbreak := s.Cities.HasOutbreak()
	if hasOutbreak {
		oCity := s.Cities.FindOutbreakCity()
		out[s.CurrentPlayer().Name][OutbreakAction] = []string{oCity.Name}
	} else {
		s.OutbreakCities = Cities{}
	}

	hasEpidemic := s.CurrentPlayer().Hand.HasEpidemic()
	if hasEpidemic {
		out[s.CurrentPlayer().Name][EpidemicAction] = []string{"epidemic"}
	}

	if !hasEpidemic && !hasOutbreak {
		if s.Step == InfectionStep {
			out[s.CurrentPlayer().Name][InfectAction] = []string{"draw"}
		}

		if s.Step == DrawStep {
			out[s.CurrentPlayer().Name][DrawAction] = []string{"draw"}
		}

		if s.Step == DiscardStep {
			out[s.CurrentPlayer().Name][DiscardAction] = s.CurrentPlayer().Hand.CardNames()
		}
	}

	for _, player := range s.Players {
		if player.Name == s.CurrentPlayer().Name && s.Step == ActionStep && !hasOutbreak && !hasEpidemic {
			playerCity := s.Cities.FindByName(player.Location)
			out[player.Name][MoveAction] = playerCity.Links

			if playerCity.Buildings[ResearchBuilding] {
				citiesWithResearch := s.Cities.GetAllWithBuilding(ResearchBuilding).FilterOne(player.Location)
				if len(citiesWithResearch) > 0 {
					out[player.Name][MoveAction] = append(out[player.Name][MoveAction], citiesWithResearch.CityNames()...)
				}

				groups := player.Hand.HasN(5)
				for virus, deck := range groups {
					if s.Viruses[virus] != NoneVirusStatus {
						continue
					}
					r := append([]string{string(virus)}, deck.CardNames()...)
					out[player.Name][ResearchAction] = append(out[player.Name][ResearchAction], strings.Join(r, ":"))
				}
			}
			if playerCity.HasVirus() {
				cureing := []string{}
				for virusName := range playerCity.GetActive() {
					cureing = append(cureing, player.Location+":"+string(virusName))
				}
				out[player.Name][CureAction] = cureing
			}

			for _, oPlayer := range s.Players {
				if oPlayer.Name == s.CurrentPlayer().Name {
					continue
				}
				if oPlayer.Location == player.Location {
					if oPlayer.Hand.Contains(player.Location) {
						out[player.Name][GetCardAction] = []string{player.Location + ":" + oPlayer.Name}
					}
				}
			}

			for _, card := range player.Hand.Cards {
				switch card.Type {
				case CityCardType:
					out[player.Name][FlyToAction] = append(out[player.Name][FlyToAction], card.Name)
					if card.Name == player.Location {
						out[player.Name][FlyAnywhereAction] = []string{"anywhere"}

						if !playerCity.Buildings[ResearchBuilding] {
							out[player.Name][BuildAction] = []string{player.Location + ":" + string(ResearchBuilding)}
						}

						for _, oPlayer := range s.Players {
							if oPlayer.Name == s.CurrentPlayer().Name {
								continue
							}
							if oPlayer.Location == player.Location {
								out[player.Name][GiveCardAction] = []string{card.Name + ":" + oPlayer.Name}
								if oPlayer.Hand.Contains(player.Location) {
									out[player.Name][GetCardAction] = []string{card.Name + ":" + oPlayer.Name}
								}
							}

						}
					}

				}
			}
		} else {
			for _, card := range player.Hand.Cards {
				switch card.Type {
				case EventCardType:
					out[player.Name][Action(card.Name)] = []string{player.Name}
				}
			}
		}
	}
	return out
}

func (s *State) Update() {
	if s.Step == StartStep {
		s.ActionCount = 4
		s.Step = ActionStep
	}
	if s.Step == ActionStep {
		if s.ActionCount < 1 {
			s.DrawCount = 2
			s.Step = DrawStep
		}
	}
	if s.Step == DrawStep {
		if s.DrawCount < 1 {
			s.Step = DiscardStep
		}
	}
	if s.Step == DiscardStep {
		if s.CurrentPlayer().Hand.Count() < 8 {
			s.InfectCount = s.InfectionLevel
			s.Step = InfectionStep
		}
	}
	if s.Step == InfectionStep {
		if s.InfectCount < 1 {
			s.Step = NextPlayerStep
		}
	}
	if s.Step == NextPlayerStep {
		s.CurrentPlayerN = (s.CurrentPlayerN + 1) % playerCount
		s.Step = StartStep
		s.Update()
	}

	s.HasWon()
	s.Save()
}
