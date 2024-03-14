package peer

/// The peer state is used to determine the current state of the peer
/// within an swap transaction.
/// The state helps identify how the transaction has evolved
type PeerState int

const (
	PeerStateDefault PeerState = iota
	PeerState0
	PeerState1
)

