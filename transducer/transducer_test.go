package transducer

import (
	"fmt"
	"testing"

	"github.com/monotykamary/nested-modal-transducers-golang-implementation/transducer/graph"
)

func TestTrafficLightShortestPaths(t *testing.T) {
	config := CreateConfig().SetState(Green)
	transducer := NewTrafficLightTransducer(config)

	paths, edges := transducer.GetShortestPaths()
	currentState := graph.Vertex(Green)

	for _, path := range paths {
		for j, state := range path {
			if path[0] != currentState {
				stateLog := fmt.Sprintf("\n'%v' state", path[0])
				currentState = path[0]
				fmt.Println(stateLog)
			}

			if j < len(path)-1 {
				nextState := path[j+1]

				input := edges[state][nextState]
				transitionLog := fmt.Sprintf("\t\t%-30v + %-30v -> %-30v", state, input, nextState)
				fmt.Println(transitionLog)
			}
		}
	}
}

func TestPedestrianLightShortestPaths(t *testing.T) {
	config := CreateConfig().SetState(Walk)
	transducer := NewPedestrianLightTransducer(config)

	paths, edges := transducer.GetShortestPaths()
	currentState := graph.Vertex(Walk)

	for _, path := range paths {
		for j, state := range path {
			if path[0] != currentState {
				stateLog := fmt.Sprintf("\n'%v' state", path[0])
				currentState = path[0]
				fmt.Println(stateLog)
			}

			if j < len(path)-1 {
				nextState := path[j+1]

				input := edges[state][nextState]
				transitionLog := fmt.Sprintf("\t\t%-30v + %-30v -> %-30v", state, input, nextState)
				fmt.Println(transitionLog)
			}
		}
	}
}
