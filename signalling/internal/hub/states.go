package hub

type State int

const (
	New State = iota
	Registered
	OfferCreated
	OfferAccepted
	Connected
)
