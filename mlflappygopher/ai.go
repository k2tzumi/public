package main

import (
	"time"
)

type ai struct {
	bird *bird

	currentRound *round
	rounds       []*round
}

func newAI(b *bird) *ai {
	return &ai{
		bird: b,
	}
}

func (ai *ai) feedInputLayer(dx, dy int32) {
	if ai.currentRound == nil {
		ai.currentRound = &round{start: time.Now()}
	}
	// fmt.Println("perceived distances", dx, dy)
	// TODO: trigger bird flapping.
	// if rand.Float32() > 0.95 {
	// 	ai.bird.jump()
	// 	fmt.Println("flapped the wings")
	// }
}

func (ai *ai) iterate() {
	// TODO: reset neural network and evolve.
	ai.currentRound.duration = time.Since(ai.currentRound.start)
	ai.rounds = append(ai.rounds, ai.currentRound)
	ai.currentRound = nil
}

type round struct {
	start    time.Time
	duration time.Duration

	// TODO: add neural net info.
}
