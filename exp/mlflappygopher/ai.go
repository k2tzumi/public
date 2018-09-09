// Copyright 2017 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/sjwhitworth/golearn/neural"
	"gonum.org/v1/gonum/mat"
)

type ai struct {
	bird              *bird
	population        []*round
	currentRound      int
	currentGeneration int
	msgCount          int
}

func newAI(b *bird) *ai {
	ai := &ai{
		bird: b,
	}
	for i := 0; i < 8; i++ {
		randNet := newRandomNetwork()
		round := &round{neuralNet: randNet}
		ai.population = append(ai.population, round)
	}
	ai.currentRound = 0
	return ai
}

func (ai *ai) feedInputLayer(dx, dy int32) {
	if ai.msgCount%20 == 0 {
		go fmt.Println(ai.currentGeneration, ":", ai.currentRound, ":", dy)
	}
	ai.msgCount++
	round := ai.population[ai.currentRound]
	if round.start.IsZero() {
		round.start = time.Now()
	}
	flap := round.neuralNet.activate(dx, dy) > 0.5
	if flap {
		ai.bird.jump()
	}
}

func (ai *ai) iterate() {
	r := ai.population[ai.currentRound]
	r.duration = time.Since(r.start)
	fmt.Println(ai.currentGeneration, ":", ai.currentRound, ":", r.fitness())
	ai.currentRound++
	if ai.currentRound >= len(ai.population) {
		sort.Slice(ai.population, func(i, j int) bool {
			return ai.population[i].fitness() > ai.population[j].fitness()
		})
		fmt.Println("generation", ai.currentGeneration, "fitness:")
		for _, round := range ai.population {
			fmt.Println(round.fitness())
		}
		ai.currentGeneration++
		ai.currentRound = 0
		pop := ai.population[0:4]
		pop = append(pop, genCross(pop[0], pop[1])...)
		pop = append(pop, genCross(pop[0], pop[2])...)
		pop = append(pop, genCross(pop[0], pop[3])...)
		for i := 0; i < 3; i++ {
			pop = append(pop, genFakeMutate(pop[0])...)
		}
		for i := range pop {
			pop[i].start = time.Time{}
			pop[i].duration = 0
		}
		ai.population = pop
		fmt.Println("new population size", len(ai.population))
		fmt.Println("new generation:", ai.currentGeneration)
	}
}

type round struct {
	start     time.Time
	duration  time.Duration
	neuralNet *neuralNet
}

func (r *round) fitness() time.Duration {
	return r.duration
}

const (
	inputLayer   = 2
	hiddenLayer1 = 8
	hiddenLayer2 = 8
	output       = 1
	totalNodes   = inputLayer + hiddenLayer1 + hiddenLayer2 + output
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
			w := rand.Float64()
			if rand.Float64() > 0.5 {
				w *= -1
			}
			nn.weights = append(nn.weights, weight{src: j, dst: i, w: w})
		}
	}
	for j := inputLayer + 1; j <= inputLayer+hiddenLayer1; j++ {
		for i := inputLayer + hiddenLayer1 + 1; i <= inputLayer+hiddenLayer1+hiddenLayer2; i++ {
			w := rand.Float64()
			if rand.Float64() > 0.5 {
				w *= -1
			}
			nn.weights = append(nn.weights, weight{src: j, dst: i, w: w})
		}
	}
	for j := inputLayer + hiddenLayer1 + 1; j <= inputLayer+hiddenLayer1+hiddenLayer2; j++ {
		for i := inputLayer + hiddenLayer1 + hiddenLayer2 + 1; i <= inputLayer+hiddenLayer1+hiddenLayer2+output; i++ {
			w := rand.Float64()
			if rand.Float64() > 0.5 {
				w *= -1
			}
			nn.weights = append(nn.weights, weight{src: j, dst: i, w: w})
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

func (nn *neuralNet) activate(dx, dy int32) float64 {
	a := mat.NewDense(totalNodes, 1, make([]float64, totalNodes))
	a.Set(0, 0, float64(dx))
	a.Set(1, 0, float64(dy))
	robot := nn.create()
	const numLayers = 4
	robot.Activate(a, numLayers)
	return a.At(totalNodes-1, 0)
}

func genCross(x, y *round) []*round {
	if len(x.neuralNet.weights) != len(x.neuralNet.weights) {
		panic("genoma lengths don't match. bug found.")
	}
	var ret []*round
	xy := &round{neuralNet: &neuralNet{}}
	for i := 0; i < len(x.neuralNet.weights); i++ {
		if i%2 == 0 {
			xy.neuralNet.weights = append(xy.neuralNet.weights, x.neuralNet.weights[i])
		} else {
			xy.neuralNet.weights = append(xy.neuralNet.weights, y.neuralNet.weights[i])
		}
	}
	ret = append(ret, xy)
	yx := &round{neuralNet: &neuralNet{}}
	for i := 0; i < len(x.neuralNet.weights); i++ {
		if i%2 == 0 {
			yx.neuralNet.weights = append(yx.neuralNet.weights, y.neuralNet.weights[i])
		} else {
			yx.neuralNet.weights = append(yx.neuralNet.weights, x.neuralNet.weights[i])
		}
	}
	ret = append(ret, yx)
	return ret
}

func genFakeMutate(x *round) []*round {
	round := &round{neuralNet: newRandomNetwork()}
	return genCross(x, round)
}
