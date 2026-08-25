package steam

type SteamSystem struct {
	header    *Header
	pressure  *Pressure
	purgeOpen bool
	trapOpen  bool
	valveOpen bool
	ventOpen  bool
}

func NewSystem(header *Header) *SteamSystem {
	if header == nil {
		header = NewHeader()
	}
	return &SteamSystem{header: header, pressure: NewPressure()}
}

func (s *SteamSystem) Header() *Header {
	return s.header
}

func (s *SteamSystem) Pressure() *Pressure {
	return s.pressure
}
