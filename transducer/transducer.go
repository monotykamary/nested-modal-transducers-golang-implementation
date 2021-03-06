package transducer

import (
	"fmt"

	"github.com/monotykamary/nested-modal-transducers-golang-implementation/transducer/graph"
)

type State string
type Input string
type Effect int

func (s State) String() string {
	return string(s)
}

func (i Input) String() string {
	return string(i)
}

func (e Effect) Int() int {
	return int(e)
}

type TransitionTable map[StateInputTuple]func() *Outputs
type StateInputTuple struct {
	State State
	Input Input
}

type Transducer struct {
	Name            string
	TransitionTable TransitionTable
}

func (t *Transducer) Transduce(config *Config, input Input) *Outputs {
	state := config.State
	stateInputTuple := StateInputTuple{State: state, Input: input}
	transitionTable := t.TransitionTable

	f, exists := transitionTable[stateInputTuple]
	if !exists {
		outputs := Outputs{Config: config}
		return &outputs
	}
	outputs := f()
	return outputs
}

func (t *Transducer) TransduceAll(configs []*Config, input Input) *Outputs {
	composedConfigs := []*Config{}
	for _, config := range configs {
		outputs := t.Transduce(config, input)
		composedConfigs = append(composedConfigs, outputs.Config)
	}
	metadata := Metadata{ParallelConfigs: composedConfigs}
	config := &Config{Metadata: metadata}
	outputs := Outputs{Config: config}
	return &outputs
}

func (t *Transducer) TransduceByBlockNumber(config *Config, input Input, txHashBN int64, currentBN int64) *Outputs {
	if txHashBN >= currentBN {
		outputs := t.Transduce(config, input)
		return outputs
	}
	state := config.State
	outputs := Outputs{Config: &Config{State: state}}
	return &outputs
}

func (t *Transducer) Rehearse(state State, inputs []Input) (State, []Effect) {
	finalState := Invalid
	effects := []Effect{}
	config := &Config{State: state}

	for _, input := range inputs {
		outputs := t.Transduce(config, input)
		finalState = outputs.Config.State
		config.State = finalState
		effects = append(effects, outputs.Effects...)
	}

	return finalState, effects
}

func (t *Transducer) ToSQL(initialState State) (string, string) {
	stateMap := map[State]map[Input]State{}
	transitionQuery := fmt.Sprintf("CREATE OR REPLACE FUNCTION %v_transition(state text, event text) RETURNS text ", t.Name)
	transitionQuery += "LANGUAGE sql AS $$ "
	transitionQuery += "SELECT CASE state "

	for stateInputTuple, f := range t.TransitionTable {
		state := stateInputTuple.State
		input := stateInputTuple.Input
		outputs := f()
		_, exists := stateMap[state]
		if !exists {
			stateMap[state] = map[Input]State{}
		}
		stateMap[state][input] = outputs.Config.State
	}

	for state, inputMap := range stateMap {
		transitionQuery += fmt.Sprintf("WHEN '%v' THEN CASE EVENT ", state)
		for input, nextState := range inputMap {
			transitionQuery += fmt.Sprintf("WHEN '%v' THEN '%v' ", input, nextState)
		}
		transitionQuery += "ELSE state END "
	}
	transitionQuery += "END $$;"

	aggregateQuery := fmt.Sprintf("CREATE AGGREGATE %v_fsm(text) (", t.Name)
	aggregateQuery += fmt.Sprintf("SFUNC = %v_transition, ", t.Name)
	aggregateQuery += "STYPE = text, "
	aggregateQuery += fmt.Sprintf("INITCOND = '%v');", initialState)

	return transitionQuery, aggregateQuery
}

func (t *Transducer) ToDiGraph() string {
	stateMap := map[State]map[Input]State{}
	digraph := "digraph {\n"

	for stateInputTuple, f := range t.TransitionTable {
		state := stateInputTuple.State
		input := stateInputTuple.Input
		outputs := f()
		_, exists := stateMap[state]
		if !exists {
			stateMap[state] = map[Input]State{}
		}
		stateMap[state][input] = outputs.Config.State
	}

	for state, inputMap := range stateMap {
		for input, nextState := range inputMap {
			digraph += "\t" + state.String() + " -> "
			label := fmt.Sprintf(`[label="%v"]`, input)
			digraph += nextState.String() + label + ";\n"
		}
	}

	digraph += "}"
	return digraph
}

func (t *Transducer) GetShortestPaths() ([][]graph.Vertex, map[graph.Vertex]map[graph.Vertex]Input) {
	vertexes := []graph.Vertex{}
	edges := make(map[graph.Vertex]map[graph.Vertex]Input)
	edgeWeights := make(map[graph.Vertex]map[graph.Vertex]float64)
	weight := 0.0

	for stateInputTuple, f := range t.TransitionTable {
		state := graph.Vertex(stateInputTuple.State)
		input := stateInputTuple.Input
		outputs := f()
		nextState := graph.Vertex(outputs.GetState())

		vertexes = append(vertexes, state)

		_, exists := edges[state]
		if !exists {
			edges[state] = make(map[graph.Vertex]Input)
		}
		edges[state][nextState] = input

		_, exists = edgeWeights[state]
		if !exists {
			edgeWeights[state] = make(map[graph.Vertex]float64)
			weight += 1
		}
		_, exists = edgeWeights[state][nextState]
		if !exists {
			weight += 1
		}
		edgeWeights[state][nextState] = weight
	}

	g := graph.NewGraph(vertexes, edgeWeights)
	dist, next := graph.FloydWarshall(g)
	paths := [][]graph.Vertex{}

	for u, m := range dist {
		for v := range m {
			if u != v {
				nextPaths := graph.GetPaths(u, v, next)
				paths = append(paths, nextPaths)
			}
		}
	}

	return paths, edges
}

type ChildTransducer map[string]Transducer

func MergeChildOutputs(outerOutputs *Outputs, innerOutputs *Outputs, childName string) *Outputs {
	effects := []Effect{}
	effects = append(effects, outerOutputs.Effects...)
	effects = append(effects, innerOutputs.Effects...)

	if outerOutputs.Config.Metadata.ChildConfig == nil {
		outerOutputs.Config.Metadata.ChildConfig = map[string]*Config{}
	}
	outerOutputs.Config.Metadata.ChildConfig[childName] = innerOutputs.Config
	outerOutputs.Effects = effects
	return outerOutputs
}

func MergeParallelOutputs(outerOutputs *Outputs, innerOutputs Outputs) Outputs {
	effects := []Effect{}
	effects = append(effects, outerOutputs.Effects...)
	effects = append(effects, innerOutputs.Effects...)

	outerOutputs.Config.Metadata.ParallelConfigs = innerOutputs.Config.Metadata.ParallelConfigs
	outerOutputs.Effects = effects
	return *outerOutputs
}

func MapChildTransducers(transducers ...Transducer) ChildTransducer {
	childTransducers := ChildTransducer{}
	for _, transducer := range transducers {
		childTransducers[transducer.Name] = transducer
	}
	return childTransducers
}
