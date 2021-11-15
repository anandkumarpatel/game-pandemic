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

func (s City) GetOutbreak() VirusType {
	for virus, count := range s.VirusCounts {
		if count > 3 {
			return virus
		}
	}

	panic("GetOutbreak: Could not find outbreak")
}

func (s Cities) FindByName(name string) *City {
	for _, city := range s {
		if city.Name == name {
			return city
		}
	}

	panic(fmt.Errorf("FindByName: cant find city %s", name))
}

func (s Cities) Contains(name string) bool {
	for _, city := range s {
		if city.Name == name {
			return true
		}
	}

	return false
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

func (s Cities) HasOutbreak() bool {
	for _, city := range s {
		for _, count := range city.VirusCounts {
			if count > 3 {
				return true
			}
		}
	}

	return false
}

func (s Cities) FindOutbreakCity() *City {
	for _, city := range s {
		for _, count := range city.VirusCounts {
			if count > 3 {
				return city
			}
		}
	}

	panic("FindOutbreak: cant find a city with outbreak")
}
