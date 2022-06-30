package transducer

// States
const (
	// an invalid state for empty states
	Invalid State = "invalid"

	Green         State = "GREEN"
	Yellow        State = "YELLOW"
	Red           State = "RED"
	PedestrianRed State = "PEDESTRIAN_RED"

	Walk State = "WALK"
	Wait State = "WAIT"
	Stop State = "STOP"
)

// Inputs
const (
	Timer           Input = "TIMER"
	PedestrianTimer Input = "PED_TIMER"
)

// Effects
const (
	UpdateTrafficColor Effect = iota
	UpdatePedestrianSymbol
)

// Names
const (
	TrafficLightTransducerName    = "traffic_light"
	PedestrianLightTransducerName = "pedestrian_light"
)
