package main

type VirusStatus int

const (
	NoneVirusStatus VirusStatus = iota
	CuredVirusStatus
	EradicatedVirusStatus
)

type Viruses struct {
	Blue   VirusStatus
	Red    VirusStatus
	Yellow VirusStatus
	Black  VirusStatus
}
