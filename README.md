# Nested Modal Transducers on Golang

This is a sample implementation of nested modal transducers from Chris's specification: https://github.com/cpressey/Nested-Modal-Transducers. This is not a complete 1:1 implementation of the transducer, as there is a difference in outputs. This is due to the fact that the original implementation for this is currently in use in production.project.


## What example are we using

This repository attempts to replicate the traffic light machine as designed on the [XState docs](https://xstate.js.org/docs-v3/#/)

![traffic light machine](https://imgur.com/GDZAeB9.png)

### Differences

Since we don't have implicit side effects, our transition table also includes a filler state `PEDESTRIAN_RED` to represent that the current parent state has its child state in transition.
