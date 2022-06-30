package transducer

func NewPedestrianLightTransducer(setupConfig *Config) *Transducer {
	transitionTable := TransitionTable{
		{Stop, PedestrianTimer}: func() *Outputs {
			return CreateOutputs().
				SetState(Walk).
				AddEffect(UpdatePedestrianSymbol)
		},

		{Walk, PedestrianTimer}: func() *Outputs {
			return CreateOutputs().
				SetState(Wait).
				AddEffect(UpdatePedestrianSymbol)
		},

		{Wait, PedestrianTimer}: func() *Outputs {
			return CreateOutputs().
				SetState(Stop).
				AddEffect(UpdatePedestrianSymbol)
		},
	}

	return &Transducer{
		Name:            PedestrianLightTransducerName,
		TransitionTable: transitionTable,
	}
}

func NewLightMachine(trafficLightStateStr string, pedestrianLightStateStr string) (*Config, *Transducer) {
	pedestrianLightState := State(pedestrianLightStateStr)
	pedestrianLightConfig := CreateConfig().SetState(pedestrianLightState)
	pedestrianLightTransducer := NewPedestrianLightTransducer(pedestrianLightConfig)

	trafficLightState := State(trafficLightStateStr)
	config := CreateConfig().
		SetState(trafficLightState).
		SetChildConfig(pedestrianLightTransducer.Name, pedestrianLightConfig)
	trafficLightTransducer := NewTrafficLightTransducer(config, *pedestrianLightTransducer)

	return config, trafficLightTransducer
}
