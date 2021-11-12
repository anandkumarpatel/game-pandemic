package main

import "fmt"

type Cities []*City

type Region string

const (
	Blue   Region = "Blue"
	Red    Region = "Red"
	Yellow Region = "Yellow"
	Black  Region = "Black"
)

type City struct {
	Name       string
	VirusCount map[string]int
	Buildings  map[string]bool
	Links      []string
	Region     Region
}

func (s City) String() string {
	return fmt.Sprintf("%s(%s)", s.Name, s.Links)
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
