// Copyright 2019 github.com/ucirello and https://cirello.io. All rights reserved.
//
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

package fogcount // import "cirello.io/gunning-fog/fogcount"

import (
	"bufio"
	"io"
)

// Analyze processes an English text and produces its Gunning Fox index score.
// Refer to its logic in https://en.wikipedia.org/wiki/Gunning_fog_index - it
// does not analyse word endings (-es, -ed, or -ing), or discriminate proper
// nouns, familiar jargon or compound words.
func Analyze(rdr io.Reader) float64 {
	var (
		phraseSize          int
		phraseCount         int
		hardWords           int
		words               int
		totalSentenceLength int
	)

	scanner := bufio.NewScanner(rdr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		words++
		phraseSize++
		w := scanner.Text()

		lastRune := w[len(w)-1]
		switch lastRune {
		case '.', ',', ';':
			w = w[:len(w)-1]
			totalSentenceLength += phraseSize
			phraseCount++
			phraseSize = 0
		}

		if len(w) > 6 {
			hardWords++
		}
	}

	return 0.4 * (float64(totalSentenceLength/phraseCount + 100*hardWords/words))
}
