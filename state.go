package main

import (
	"fmt"
	"strings"
)

type Turn struct {
	CurrentPlayer *Player
	ActionCount   int
	Step          Step
	DrawCount     int
	InfectCount   int
}

type State struct {
	Players        []*Player
	Cities         Cities
	Viruses        Viruses
	Decks          Decks
	InfectionLevel int
	Turn           *Turn
}

func (s State) String() string {
	out := ""
	for i, p := range s.Players {
		out += fmt.Sprintf("Player(%d) %s\n", i, p)
	}
	out += fmt.Sprintf("Cities %s\n", s.Cities)
	out += fmt.Sprintf("Viruses %v\n", s.Viruses)
	out += fmt.Sprintf("PDeck %s\n", s.Decks.PDeck)
	out += fmt.Sprintf("VDeck %s\n", s.Decks.VDeck)
	out += fmt.Sprintf("PDiscard %s\n", s.Decks.PDiscard)
	out += fmt.Sprintf("VDiscard %s\n", s.Decks.VDiscard)
	out += fmt.Sprintf("Current Player %+v\n", s.Turn)

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
	NextAction        Action = "next"
	DrawPAction       Action = "drawP"
	ResearchAction    Action = "reseach"
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
	case DiscardAction:
		s.DiscardAction(target)
	case DrawPAction:
		s.DrawPAction()
	case InfectAction:
		s.InfectAction()
	case CureAction:
		s.CureAction(target)
	case BuildAction:
		s.BuildAction(target)
	case ResearchAction:
		s.ResearchAction(target)
	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	return nil
}

func (s *State) ResearchAction(target string) {
	// First split is always virus
	split := strings.Split(target, ":")
	virusName := split[0]
	for i := 1; i < len(split); i++ {
		card := s.Turn.CurrentPlayer.Hand.RemoveCard(split[i])
		s.Decks.PDiscard.AddCard(card)
	}
	s.Viruses[VirusType(virusName)] = CuredVirusStatus
	s.Turn.ActionCount--
}

func (s *State) MoveAction(cityName string) {
	s.Turn.CurrentPlayer.Move(cityName)
	s.Turn.ActionCount--
}

func (s *State) FlyToAction(cityName string) {
	card := s.Turn.CurrentPlayer.Hand.RemoveCard(cityName)
	s.Decks.PDiscard.AddCard(card)
	s.Turn.CurrentPlayer.Move(cityName)
	s.Turn.ActionCount--
}

func (s *State) DiscardAction(cityName string) {
	card := s.Turn.CurrentPlayer.Hand.RemoveCard(cityName)
	s.Decks.PDiscard.AddCard(card)
}

func (s *State) CureAction(target string) {
	split := strings.Split(target, ":")
	cityName, virusName := split[0], split[1]
	s.Cities.FindCityByName(cityName).VirusCounts[VirusType(virusName)]--
	s.Turn.ActionCount--
}

func (s *State) BuildAction(target string) {
	split := strings.Split(target, ":")
	cityName, buildName := split[0], split[1]
	card := s.Turn.CurrentPlayer.Hand.RemoveCard(cityName)
	s.Decks.PDiscard.AddCard(card)
	s.Cities.FindCityByName(cityName).Buildings[Building(buildName)] = true
	s.Turn.ActionCount--
}

func (s *State) DrawPAction() {
	card := s.Decks.PDeck.Draw()
	s.Turn.CurrentPlayer.Hand.AddCard(card)
	s.Turn.DrawCount--
}

func (s *State) InfectAction() {
	card := s.Decks.VDeck.Draw()
	city := s.Cities.FindCityByName(card.Name)
	city.VirusCounts[card.City.VirusType]++
	s.Decks.VDiscard.AddCard(card)
	s.Turn.InfectCount--
}

type ActionMap map[Action]([]string)
type PlayersActions map[string]ActionMap

func (s *State) GetActions() PlayersActions {
	out := PlayersActions{}
	for _, player := range s.Players {
		out[player.Name] = ActionMap{}
	}

	if s.Turn.Step == InfectionStep {
		out[s.Turn.CurrentPlayer.Name][InfectAction] = []string{"draw"}
	}

	if s.Turn.Step == DrawStep {
		out[s.Turn.CurrentPlayer.Name][DrawPAction] = []string{"draw"}
	}

	if s.Turn.Step == DiscardStep {
		out[s.Turn.CurrentPlayer.Name][DiscardAction] = s.Turn.CurrentPlayer.Hand.CardNames()
	}

	for _, player := range s.Players {
		if player.Name == s.Turn.CurrentPlayer.Name && s.Turn.Step == ActionStep {
			playerCity := s.Cities.FindCityByName(player.Location)
			if playerCity == nil {
				panic(fmt.Errorf("cant find player location %s", player))
			}
			out[player.Name][MoveAction] = playerCity.Links
			if playerCity.HasVirus() {
				cureing := []string{}
				for virusName, _ := range playerCity.GetActive() {
					cureing = append(cureing, player.Location+":"+string(virusName))
				}
				out[player.Name][CureAction] = cureing
			}

			for _, card := range player.Hand.Cards {
				switch card.Type {
				case CityCardType:
					out[player.Name][FlyToAction] = append(out[player.Name][FlyToAction], card.City.Name)
					if card.City.Name == player.Location {
						out[player.Name][FlyAnywhereAction] = []string{player.Location}
						out[player.Name][BuildAction] = []string{player.Location + ":" + string(ResearchBuilding)}
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
