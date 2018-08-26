package main

import (
	"fmt"
)

type ai struct {
	bird *bird
}

func newAI(b *bird) *ai {
	return &ai{
		bird: b,
	}
}

func (ai *ai) feedInputLayer(dx, dy int32) {
	fmt.Println("perceived distances", dx, dy)
	// TODO: trigger window flapping
	// if rand.Float32() > 0.95 {
	// 	ai.bird.jump()
	// 	fmt.Println("flapped the wings")
	// }
}
