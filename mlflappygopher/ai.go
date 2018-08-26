package main

import (
	"math/rand"
	"time"

	"github.com/sjwhitworth/golearn/neural"
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
	start     time.Time
	duration  time.Duration
	neuralNet neuralNet

	// TODO: add neural net info.
}

const (
	inputLayer   = 2
	hiddenLayer1 = 6
	hiddenLayer2 = 6
	output       = 1

	totalNodes = inputLayer + hiddenLayer1 + hiddenLayer2 + output
)

type weight struct {
	src, dst int
	w        float64
}
type neuralNet struct {
	weights []weight
}

func newRandomNetwork() *neuralNet {
	nn := &neuralNet{}

	for j := 1; j <= inputLayer; j++ {
		for i := inputLayer + 1; i <= inputLayer+hiddenLayer1; i++ {
			nn.weights = append(nn.weights, weight{src: j, dst: i, w: rand.Float64()})
		}
	}

	for j := inputLayer + 1; j <= inputLayer+hiddenLayer1; j++ {
		for i := inputLayer + hiddenLayer1 + 1; i <= inputLayer+hiddenLayer1+hiddenLayer2; i++ {
			nn.weights = append(nn.weights, weight{src: j, dst: i, w: rand.Float64()})
		}
	}

	for j := inputLayer + hiddenLayer1 + 1; j <= inputLayer+hiddenLayer1+hiddenLayer2; j++ {
		for i := inputLayer + hiddenLayer1 + hiddenLayer2 + 1; i <= inputLayer+hiddenLayer1+hiddenLayer2+output; i++ {
			nn.weights = append(nn.weights, weight{src: j, dst: i, w: rand.Float64()})
		}
	}

	return nn
}

func (nn *neuralNet) create() *neural.Network {
	n := neural.NewNetwork(totalNodes, inputLayer, neural.Sigmoid)
	for _, w := range nn.weights {
		n.SetWeight(w.src, w.dst, w.w)
	}
	return n
}
