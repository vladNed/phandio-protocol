package signalling

type Marketplace struct {
	offers map[string][]byte
}

func NewMarketplace() *Marketplace {
	return &Marketplace{
		offers: make(map[string][]byte),
	}
}

func (m *Marketplace) AddOffer(offerID string, offer []byte) {
	m.offers[offerID] = offer
}

func (m *Marketplace) GetOffer(offerID string) ([]byte, bool) {
	offer, ok := m.offers[offerID]
	return offer, ok
}

// TODO: Delete this function
func (m *Marketplace) Display(pageLogger func(string)) {
	for k, _ := range m.offers {
		pageLogger("Offer ID: " + k)
	}
}