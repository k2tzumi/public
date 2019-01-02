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

// Gunning-fog index analyzer written in Go. This analyzer processes an English
// text and produces its Gunning Fox index score. Refer to its logic in
// https://en.wikipedia.org/wiki/Gunning_fog_index - it does not analyse word
// endings (-es, -ed, or -ing), or discriminate proper nouns, familiar jargon or
// compound words.
//
//
//     $ go get cirello.io/gunning-fog/...
//     $ cat LICENSE | $GOPATH/bin/gunning-fog
//     16
//
// `gunning-fog` will always wait content from STDIN.
package main // import "cirello.io/gunning-fog"

import (
	"fmt"
	"os"

	"cirello.io/gunning-fog/fogcount"
)

func main() {
	fmt.Println(fogcount.Analyze(os.Stdin))
}
