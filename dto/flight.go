package dto

type Flight struct {
	Name     string
	IsStart  bool
	IsEnd    bool
	Visited  bool
	Outgoing []*Flight
	Incoming []*Flight
}
