package transducer

func NewTrafficLightTransducer(setupConfig *Config, childTransducers ...Transducer) *Transducer {
	childTransducersMap := MapChildTransducers(childTransducers...)

	transitionTable := TransitionTable{
		{Green, Timer}: func() *Outputs {
			return CreateOutputs().
				SetState(Yellow).
				AddEffect(UpdateTrafficColor)
		},

		{Yellow, Timer}: func() *Outputs {
			return CreateOutputs().
				SetState(PedestrianRed).
				AddEffect(UpdateTrafficColor)
		},

		{PedestrianRed, PedestrianTimer}: func() *Outputs {
			outputs := CreateOutputs().
				SetState(PedestrianRed).
				TransduceChild(childTransducersMap[PedestrianLightTransducerName], setupConfig, PedestrianTimer)

			childState := outputs.Config.GetChildState(PedestrianLightTransducerName)
			switch childState {
			case Walk:
			case Wait:
			case Stop:
				outputs.SetState(Red)
			}

			return outputs
		},

		{Red, Timer}: func() *Outputs {
			return CreateOutputs().
				SetState(Green).
				AddEffect(UpdateTrafficColor)
		},
	}

	return &Transducer{
		Name:            TrafficLightTransducerName,
		TransitionTable: transitionTable,
	}
}

func NewTrafficLightMachine(stateStr string) (*Config, *Transducer) {
	state := State(stateStr)
	config := CreateConfig().SetState(state)
	transducer := NewTrafficLightTransducer(config)

	return config, transducer
}
