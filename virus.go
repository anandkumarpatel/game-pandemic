package main

type VirusType string

const (
	Blue   VirusType = "Blue"
	Red    VirusType = "Red"
	Yellow VirusType = "Yellow"
	Black  VirusType = "Black"
)

type VirusStatus int

const (
	NoneVirusStatus VirusStatus = iota
	CuredVirusStatus
	EradicatedVirusStatus
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
