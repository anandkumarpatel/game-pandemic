package main

import "fmt"

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
	Virus          Viruses
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
	out += fmt.Sprintf("Virus %v\n", s.Virus)
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
	DrawPAction       Action = "DrawPAction"
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
	case InfectAction:
		s.InfectAction()
	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	return nil
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

func (s *State) InfectAction() {
	card := s.Decks.VDeck.Draw()
	city := s.Cities.FindCityByName(card.Name)
	city.VirusCount[string(card.City.Region)]++
	s.Decks.VDiscard.AddCard(card)
}

type ActionMap map[Action]([]string)
type PlayersActions map[string]ActionMap

func (s *State) GetActions() PlayersActions {
	out := PlayersActions{}
	if s.Turn.ActionCount == 0 {
		out[s.Turn.CurrentPlayer.Name] = ActionMap{}
		out[s.Turn.CurrentPlayer.Name][DrawPAction] = []string{"draw"}
		return out
	}
	for _, player := range s.Players {
		if player.Name == s.Turn.CurrentPlayer.Name {
			out[player.Name] = ActionMap{}
			playerLocation := s.Cities.FindCityByName(player.Location)
			if playerLocation == nil {
				panic(fmt.Errorf("cant find player location %s", player))
			}
			out[player.Name][MoveAction] = playerLocation.Links
			out[player.Name][CureAction] = []string{playerLocation.Name}

			for _, card := range player.Hand.Cards {
				switch card.Type {
				case CityCardType:
					out[player.Name][FlyToAction] = append(out[player.Name][FlyToAction], card.City.Name)
					if card.City.Name == playerLocation.Name {
						out[player.Name][FlyAnywhereAction] = []string{playerLocation.Name}
						out[player.Name][BuildAction] = []string{playerLocation.Name}
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
