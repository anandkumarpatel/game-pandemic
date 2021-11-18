package main

type VirusType string

const (
	Blue   VirusType = "Blue"
	Red    VirusType = "Red"
	Yellow VirusType = "Yellow"
	Black  VirusType = "Black"
)

type VirusStatus string

const (
	NoneVirusStatus       VirusStatus = "exists"
	CuredVirusStatus      VirusStatus = "cured"
	EradicatedVirusStatus VirusStatus = "eradicated"
)

type Viruses map[VirusType]VirusStatus

func (s Viruses) AllCured() bool {
	for _, status := range s {
		if status == NoneVirusStatus {
			return false
		}
	}

	return true
}
