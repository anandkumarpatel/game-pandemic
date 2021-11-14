package main

import "fmt"

type Cities []*City
type VirusCounts map[VirusType]int

type Building string
type Buildings map[Building]bool

const (
	ResearchBuilding Building = "ResearchBuilding"
)

type City struct {
	Name        string
	VirusCounts VirusCounts
	Buildings   Buildings
	Links       []string
	VirusType   VirusType
}

func (s City) String() string {
	out := fmt.Sprintf("|%s%s", s.Name, s.Links)
	for virus, count := range s.VirusCounts {
		if count > 0 {
			out += fmt.Sprintf(",%s%d", string(virus[0]), count)
		}
	}

	for building, has := range s.Buildings {
		if has {
			out += fmt.Sprintf("^%s", string(building[0]))
		}
	}

	out += "|"
	return out
}

func (s City) HasVirus() bool {
	for _, count := range s.VirusCounts {
		if count > 0 {
			return true
		}
	}
	return false
}

func (s City) GetActive() VirusCounts {
	out := VirusCounts{}
	for name, count := range s.VirusCounts {
		if count > 0 {
			out[name] = count
		}
	}
	return out
}

func (s Cities) FindCityByName(name string) *City {
	for _, v := range s {
		if v.Name == name {
			return v
		}
	}

	panic(fmt.Errorf("cant find city %s", name))
	// return nil
}

func (s Cities) Strings() []string {
	cities := []string{}
	for _, city := range s {
		cities = append(cities, city.Name)
	}
	return cities
}

func (s Cities) IsEradicated(v VirusType) bool {
	for _, c := range s {
		if c.VirusCounts[v] > 0 {
			return false
		}
	}

	return true
}
