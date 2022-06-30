package transducer

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRehearsePedestrianLightTransducer(t *testing.T) {
	config := CreateConfig().SetState(Stop)
	transducer := NewPedestrianLightTransducer(config)

	nextState, nextEffects := transducer.Rehearse(Stop, []Input{PedestrianTimer})
	if nextState != Walk {
		t.Errorf("characterTransducer.Rehearse() failed, expected %v, got %v", Walk, nextState)
	}
	targetEffects := []Effect{UpdatePedestrianSymbol}
	if !reflect.DeepEqual(nextEffects, targetEffects) {
		t.Errorf("characterTransducer.Rehearse() failed, expected %v, got %v", targetEffects, nextEffects)
	}

	nextState, nextEffects = transducer.Rehearse(Walk, []Input{PedestrianTimer})
	if nextState != Wait {
		t.Errorf("characterTransducer.Rehearse() failed, expected %v, got %v", Wait, nextState)
	}
	targetEffects = []Effect{UpdatePedestrianSymbol}
	if !reflect.DeepEqual(nextEffects, targetEffects) {
		t.Errorf("characterTransducer.Rehearse() failed, expected %v, got %v", targetEffects, nextEffects)
	}

	nextState, nextEffects = transducer.Rehearse(Wait, []Input{PedestrianTimer})
	if nextState != Stop {
		t.Errorf("characterTransducer.Rehearse() failed, expected %v, got %v", Stop, nextState)
	}
	targetEffects = []Effect{UpdatePedestrianSymbol}
	if !reflect.DeepEqual(nextEffects, targetEffects) {
		t.Errorf("characterTransducer.Rehearse() failed, expected %v, got %v", targetEffects, nextEffects)
	}
}

func TestLightMachineNested(t *testing.T) {
	config, transducer := NewLightMachine(PedestrianRed.String(), Stop.String())
	outputs := transducer.Transduce(config, PedestrianTimer)

	nextState := outputs.GetState()
	if nextState != PedestrianRed {
		t.Errorf("characterTransducerNested.Rehearse() failed, expected %v, got %v", PedestrianRed, nextState)
	}
	nextChildState := outputs.GetChildState(PedestrianLightTransducerName)
	if nextChildState != Walk {
		t.Errorf("characterTransducerNested.Rehearse() failed, expected %v, got %v", Walk, nextChildState)
	}
	nextEffects := outputs.Effects
	targetEffects := []Effect{UpdatePedestrianSymbol}
	if !reflect.DeepEqual(nextEffects, targetEffects) {
		t.Errorf("characterTransducer.Rehearse() failed, expected %v, got %v", targetEffects, nextEffects)
	}
}

func TestPedestrianLightTransducerToSQL(t *testing.T) {
	config := CreateConfig().SetState(Stop)
	transducer := NewPedestrianLightTransducer(config)

	transitionQuery, aggregateQuery := transducer.ToSQL(Stop)
	fmt.Printf("Transition Query: %v\n", transitionQuery)
	fmt.Printf("Aggregate Query: %v\n", aggregateQuery)
}

func TestPedestrianLightTransducerToDigraph(t *testing.T) {
	config := CreateConfig().SetState(Stop)
	transducer := NewPedestrianLightTransducer(config)

	digraph := transducer.ToDiGraph()
	t.Logf("Digraph: %v\n", digraph)
}
