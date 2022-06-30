package main

import (
	"fmt"

	t "github.com/monotykamary/nested-modal-transducers-golang-implementation/transducer"
)

func main() {
	config, lightTransducer := t.NewLightMachine(t.PedestrianRed.String(), t.Stop.String())
	outputs := lightTransducer.Transduce(config, t.PedestrianTimer)

	nextState := outputs.GetState()
	nextChildState := outputs.GetChildState(t.PedestrianLightTransducerName)

	fmt.Println(nextState.String(), nextChildState.String())
}
