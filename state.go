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
	Players        Players
	Cities         Cities
	Viruses        Viruses
	Decks          Decks
	InfectionLevel int
	Turn           *Turn
	OutbreakCount  int
	OutbreakCities Cities
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
	OutbreakAction    Action = "outbreak"
	EpidemicAction    Action = "epidemic"
	GiveCardAction    Action = "GiveCard"
	GetCardAction     Action = "GetCard"
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
		s.InfectAction(1)
	case CureAction:
		s.CureAction(target)
	case BuildAction:
		s.BuildAction(target)
	case ResearchAction:
		s.ResearchAction(target)
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

	return nil
}

func (s *State) SetupVirus() {
	for i := 3; i > 0; i-- {
		for j := 0; j < 3; j++ {
			s.InfectAction(i)
		}
	}
}

func (s *State) ResearchAction(target string) {
	// First split is always virus
	split := strings.Split(target, ":")
	virusName := split[0]
	for i := 1; i < len(split); i++ {
		card := s.Turn.CurrentPlayer.Hand.RemoveCard(split[i])
		s.Decks.PDiscard.AddCard(card)
	}
	v := VirusType(virusName)
	s.Viruses[v] = CuredVirusStatus

	if s.Cities.IsEradicated(v) {
		s.Viruses[v] = EradicatedVirusStatus
	}
	s.Turn.ActionCount--
}

func (s *State) MoveAction(cityName string) {
	s.Turn.CurrentPlayer.Move(cityName)
	s.Turn.ActionCount--
}

func (s *State) GiveCardAction(target string) {
	split := strings.Split(target, ":")
	cardName, playerName := split[0], split[1]
	card := s.Turn.CurrentPlayer.Hand.RemoveCard(cardName)
	toPlayer := s.Players.FindByName(playerName)
	toPlayer.Hand.AddCard(card)
	s.Turn.ActionCount--
}

func (s *State) GetCardAction(target string) {
	split := strings.Split(target, ":")
	cardName, playerName := split[0], split[1]
	fromPlayer := s.Players.FindByName(playerName)
	card := fromPlayer.Hand.RemoveCard(cardName)
	s.Turn.CurrentPlayer.Hand.AddCard(card)
	s.Turn.ActionCount--
}

func (s *State) EpidemicAction() {
	s.InfectionLevel++
	s.Turn.CurrentPlayer.Hand.RemoveCard("epidemic")

	card := s.Decks.VDeck.BackDraw()
	if s.Viruses[card.City.VirusType] != EradicatedVirusStatus {
		card.City.VirusCounts[card.City.VirusType] += 3
	}
	s.Decks.VDiscard.AddCard(card)
	s.Decks.VDiscard.Shuffle()

	s.Decks.VDeck.AddDeck(s.Decks.VDiscard)
	s.Decks.VDiscard.Clear()
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
	v := VirusType(virusName)

	city := s.Cities.FindByName(cityName)
	city.VirusCounts[v]--
	if s.Cities.IsEradicated(v) {
		s.Viruses[v] = EradicatedVirusStatus
	}
	s.Turn.ActionCount--
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

func (s State) HasWon() bool {
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

func (s *State) BuildAction(target string) {
	split := strings.Split(target, ":")
	cityName, buildName := split[0], split[1]
	card := s.Turn.CurrentPlayer.Hand.RemoveCard(cityName)
	s.Decks.PDiscard.AddCard(card)
	s.Cities.FindByName(cityName).Buildings[Building(buildName)] = true
	s.Turn.ActionCount--
}

func (s *State) DrawPAction() {
	card := s.Decks.PDeck.Draw()
	s.Turn.CurrentPlayer.Hand.AddCard(card)
	s.Turn.DrawCount--
}

func (s *State) InfectAction(count int) {
	card := s.Decks.VDeck.Draw()
	if s.Viruses[card.City.VirusType] != EradicatedVirusStatus {
		card.City.VirusCounts[card.City.VirusType] += count
	}
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
	hasOutbreak := s.Cities.HasOutbreak()
	if hasOutbreak {
		oCity := s.Cities.FindOutbreakCity()
		out[s.Turn.CurrentPlayer.Name][OutbreakAction] = []string{oCity.Name}
	} else {
		s.OutbreakCities = Cities{}
	}

	hasEpidemic := s.Turn.CurrentPlayer.Hand.HasEpidemic()
	if hasEpidemic {
		out[s.Turn.CurrentPlayer.Name][EpidemicAction] = []string{"epidemic"}
	}

	if !hasEpidemic && !hasOutbreak {
		if s.Turn.Step == InfectionStep {
			out[s.Turn.CurrentPlayer.Name][InfectAction] = []string{"draw"}
		}

		if s.Turn.Step == DrawStep {
			out[s.Turn.CurrentPlayer.Name][DrawPAction] = []string{"draw"}
		}

		if s.Turn.Step == DiscardStep {
			out[s.Turn.CurrentPlayer.Name][DiscardAction] = s.Turn.CurrentPlayer.Hand.CardNames()
		}
	}

	for _, player := range s.Players {
		if player.Name == s.Turn.CurrentPlayer.Name && s.Turn.Step == ActionStep && !hasOutbreak && !hasEpidemic {
			playerCity := s.Cities.FindByName(player.Location)
			if playerCity.Buildings[ResearchBuilding] {
				groups := player.Hand.HasN(5)
				for virus, deck := range groups {
					if s.Viruses[virus] != NoneVirusStatus {
						continue
					}
					r := append([]string{string(virus)}, deck.CardNames()...)
					out[player.Name][ResearchAction] = append(out[player.Name][ResearchAction], strings.Join(r, ":"))
				}
			}
			out[player.Name][MoveAction] = playerCity.Links
			if playerCity.HasVirus() {
				cureing := []string{}
				for virusName := range playerCity.GetActive() {
					cureing = append(cureing, player.Location+":"+string(virusName))
				}
				out[player.Name][CureAction] = cureing
			}

			for _, oPlayer := range s.Players {
				if oPlayer.Name == s.Turn.CurrentPlayer.Name {
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
					out[player.Name][FlyToAction] = append(out[player.Name][FlyToAction], card.City.Name)
					if card.City.Name == player.Location {
						out[player.Name][FlyAnywhereAction] = []string{player.Location}

						if !playerCity.Buildings[ResearchBuilding] {
							out[player.Name][BuildAction] = []string{player.Location + ":" + string(ResearchBuilding)}
						}

						for _, oPlayer := range s.Players {
							if oPlayer.Name == s.Turn.CurrentPlayer.Name {
								continue
							}
							if oPlayer.Location == player.Location {
								out[player.Name][GiveCardAction] = []string{card.City.Name + ":" + oPlayer.Name}
								if oPlayer.Hand.Contains(player.Location) {
									out[player.Name][GetCardAction] = []string{card.City.Name + ":" + oPlayer.Name}
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
