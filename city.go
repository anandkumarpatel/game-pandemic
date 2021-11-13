package main

import "fmt"

type Cities []*City

type City struct {
	Name       string
	VirusCount map[VirusType]int
	Buildings  map[string]bool
	Links      []string
	VirusType  VirusType
}

func (s City) String() string {
	out := fmt.Sprintf("|%s%s", s.Name, s.Links)
	for virus, count := range s.VirusCount {
		if count > 0 {
			out += fmt.Sprintf("%s{%d}", string(virus[0]), count)
		}
	}

	out += "|"
	return out
}

func (s Cities) FindCityByName(name string) *City {
	for _, v := range s {
		if v.Name == name {
			return v
		}
	}

	return nil
}

func (s Cities) Strings() []string {
	cities := []string{}
	for _, city := range s {
		cities = append(cities, city.Name)
	}
	return cities
}
