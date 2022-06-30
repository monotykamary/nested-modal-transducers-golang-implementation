package transducer

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRehearseTrafficLightTransducer(t *testing.T) {
	config := CreateConfig().SetState(Green)
	transducer := NewTrafficLightTransducer(config)

	nextState, nextEffects := transducer.Rehearse(Green, []Input{Timer})
	if nextState != Yellow {
		t.Errorf("characterTransducer.Rehearse() failed, expected %v, got %v", Yellow, nextState)
	}
	targetEffects := []Effect{UpdateTrafficColor}
	if !reflect.DeepEqual(nextEffects, targetEffects) {
		t.Errorf("characterTransducer.Rehearse() failed, expected %v, got %v", targetEffects, nextEffects)
	}

	nextState, nextEffects = transducer.Rehearse(Yellow, []Input{Timer})
	if nextState != PedestrianRed {
		t.Errorf("characterTransducer.Rehearse() failed, expected %v, got %v", PedestrianRed, nextState)
	}
	targetEffects = []Effect{UpdateTrafficColor}
	if !reflect.DeepEqual(nextEffects, targetEffects) {
		t.Errorf("characterTransducer.Rehearse() failed, expected %v, got %v", targetEffects, nextEffects)
	}

	nextState, nextEffects = transducer.Rehearse(Red, []Input{Timer})
	if nextState != Green {
		t.Errorf("characterTransducer.Rehearse() failed, expected %v, got %v", Green, nextState)
	}
	targetEffects = []Effect{UpdateTrafficColor}
	if !reflect.DeepEqual(nextEffects, targetEffects) {
		t.Errorf("characterTransducer.Rehearse() failed, expected %v, got %v", targetEffects, nextEffects)
	}
}

func TestTrafficLightTransducerToSQL(t *testing.T) {
	_, transducer := NewLightMachine(PedestrianRed.String(), Wait.String())

	transitionQuery, aggregateQuery := transducer.ToSQL(Green)
	fmt.Printf("Transition Query: %v\n", transitionQuery)
	fmt.Printf("Aggregate Query: %v\n", aggregateQuery)
}

func TestTrafficLightTransducerToDigraph(t *testing.T) {
	_, transducer := NewLightMachine(PedestrianRed.String(), Wait.String())

	digraph := transducer.ToDiGraph()
	t.Logf("Digraph: %v\n", digraph)
}
