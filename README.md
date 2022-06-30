# Nested Modal Transducers on Golang

This is a sample implementation of nested modal transducers from Chris's specification: https://github.com/cpressey/Nested-Modal-Transducers. This is not a complete 1:1 implementation of the transducer, as there is a difference in how we create outputs.

## What example state machine is this?

This repository attempts to replicate the traffic light machine as designed on the [XState docs](https://xstate.js.org/docs-v3/#/):

<img src="https://imgur.com/GDZAeB9.png" alt="traffic light" width="400px">

Since we don't have implicit side effects and nested labelling of hierarchical state machines, our implementation transition table must include a filler state `PEDESTRIAN_RED` to represent that the current parent state has its child state in transition. We need this because otherwise, a `TIMER` event could completely ignore our pedestrian light state machine and leave us in a logically invalid state. This also means that our starting state for the pedestrian light has to start at `STOP` if the traffic light starts at `GREEN`, which makes sense.

## Generated Dot Graphs

Transducers have access to a method called `.ToDiGraph()` to generate simple dot graphs. You can see an example of this in the `*_test.go` files inside the `transducer` folder:

![traffic_light](https://user-images.githubusercontent.com/1130103/176672884-8a51febe-c3be-4abd-a5d9-66227031e0ee.png)
![pedestrian_light](https://user-images.githubusercontent.com/1130103/176672891-f5a2235f-7e42-4309-96aa-773d024ccdbe.png)

## Generated SQL

Transducers have access to a method called `.ToSQL(initialState state)` to generate simple SQL state machines. You can see an example of this in the `*_test.go` files inside the `transducer` folder:

This idea came from Felix's post on state machines in PostgreSQL: https://felixge.de/2017/07/27/implementing-state-machines-in-postgresql/. After reading the article, I thought it would also be a nice QoL to include generation of SQL for ad-hoc audits:

#### Traffic Light

```sql
CREATE OR REPLACE FUNCTION traffic_light_transition(state text, event text) RETURNS text LANGUAGE sql AS $$
SELECT
    CASE state
        WHEN 'RED' THEN CASE
            EVENT
            WHEN 'TIMER' THEN 'GREEN'
            ELSE state
        END
        WHEN 'GREEN' THEN CASE
            EVENT
            WHEN 'TIMER' THEN 'YELLOW'
            ELSE state
        END
        WHEN 'YELLOW' THEN CASE
            EVENT
            WHEN 'TIMER' THEN 'PEDESTRIAN_RED'
            ELSE state
        END
        WHEN 'PEDESTRIAN_RED' THEN CASE
            EVENT
            WHEN 'PED_TIMER' THEN 'RED'
            ELSE state
        END
    END
$$;

CREATE AGGREGATE traffic_light_fsm(text) (
    SFUNC = traffic_light_transition,
    STYPE = text,
    INITCOND = 'GREEN'
);
```

#### Pedestrian Light

```sql
CREATE OR REPLACE FUNCTION pedestrian_light_transition(state text, event text) RETURNS text LANGUAGE sql AS $$
SELECT
    CASE state
        WHEN 'STOP' THEN CASE
            EVENT
            WHEN 'PED_TIMER' THEN 'WALK'
            ELSE state
        END
        WHEN 'WALK' THEN CASE
            EVENT
            WHEN 'PED_TIMER' THEN 'WAIT'
            ELSE state
        END
        WHEN 'WAIT' THEN CASE
            EVENT
            WHEN 'PED_TIMER' THEN 'STOP'
            ELSE state
        END
    END
$$;

CREATE AGGREGATE pedestrian_light_fsm(text) (
    SFUNC = pedestrian_light_transition,
    STYPE = text,
    INITCOND = 'STOP'
);
```

## Generating shortest paths

Transducers have access to a method called `.GetShortestPaths()` to generate the shortest path between states. This path is generated with the Floyd-Warshall algorithm (for simplicity); these paths have not been deduped and truncated. You can see an example of this in the `*_test.go` files inside the `transducer` folder:

#### Traffic light

```
'RED' state
		RED                            + TIMER                          -> GREEN
		RED                            + TIMER                          -> GREEN
		GREEN                          + TIMER                          -> YELLOW
		RED                            + TIMER                          -> GREEN
		GREEN                          + TIMER                          -> YELLOW
		YELLOW                         + TIMER                          -> PEDESTRIAN_RED

'GREEN' state
		GREEN                          + TIMER                          -> YELLOW
		GREEN                          + TIMER                          -> YELLOW
		YELLOW                         + TIMER                          -> PEDESTRIAN_RED

'YELLOW' state
		YELLOW                         + TIMER                          -> PEDESTRIAN_RED
```

#### Pedestrian light
```
'STOP' state
		STOP                           + PED_TIMER                      -> WALK
		WALK                           + PED_TIMER                      -> WAIT
		STOP                           + PED_TIMER                      -> WALK

'WALK' state
		WALK                           + PED_TIMER                      -> WAIT
		WALK                           + PED_TIMER                      -> WAIT
		WAIT                           + PED_TIMER                      -> STOP

'WAIT' state
		WAIT                           + PED_TIMER                      -> STOP
		WAIT                           + PED_TIMER                      -> STOP
		STOP                           + PED_TIMER                      -> WALK
```

## Drawbacks and remaining concerns

We allow for output composition instead of coupling the outputs with the transition table. This means for cases (like the `PEDESTRIAN_RED` filler state), we can have logic that isn't trivial to generate SQL or dot graphs with without adding special cases for them. For instance, our nested machine must start at `PEDESTRIAN_RED` and `STOP` for the traffic light and pedestrian light machines respectively, if we wanted to make sure we generate the entire machine on the diagram/SQL. Not ideal, but at least we maintain strict consistency.

The initial design of the transducer only considered mirrored events/inputs across nested state machines such that it would encourage flat composition of outputs (without using native conditions inside the transition function). In addition, the orthoganal state machine story isn't as flushed out yet as our hierarchical example. Probably an upcoming TODO. The method to transitioning multiple state machines in parallel is quite trivial, but it isn't trivial to generate SQL or diagrams from.
