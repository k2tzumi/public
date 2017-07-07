// FireGo
// FirePHP ported to Go.
//
// It partially implements the FirePHP Protocol, supporting:
//
// Log
// Info
// Warn
// Error
// The TRACE, EXCEPTION, TABLE and GROUP are not implemented - I still need to understand whether it is desirable and possible to port these message types.
//
// Also, it does not analyse the backtrace to feed the information with extra information such filename and line. http://golang.org/pkg/runtime/#Stack should do the trick.
//
// Check the example to see it working:
//
// # go run examples/example.go
package firego // import "cirello.io/FireGo"

//go:generate stringer -type=MessageType
import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	chunkSize = 4096
)

type message struct {
	t       MessageType
	content interface{}
	label   string
}

type FireGo struct {
	mu       sync.Mutex
	messages []message
}

func New() *FireGo {
	msgs := make([]message, 0)
	fireGo := &FireGo{
		messages: msgs,
	}
	return fireGo
}

func (f *FireGo) Message(t MessageType, content interface{}) {
	msg := message{
		t: t,
	}
	if t == GroupStart {
		msg.label = content.(string)
	} else {
		msg.content = content
	}
	f.mu.Lock()
	f.messages = append(f.messages, msg)
	f.mu.Unlock()
}

func (f *FireGo) Log(content string) {
	f.Message(Log, content)
}
func (f *FireGo) Info(content string) {
	f.Message(Info, content)
}
func (f *FireGo) Warn(content string) {
	f.Message(Warn, content)
}
func (f *FireGo) Error(content string) {
	f.Message(Error, content)
}
func (f *FireGo) Table(content [][]string) {
	f.Message(Table, content)
}
func (f *FireGo) GroupStart(label string) {
	f.Message(GroupStart, label)
}
func (f *FireGo) GroupEnd() {
	f.Message(GroupEnd, nil)
}

func (f *FireGo) Flush(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	messages := f.messages
	f.messages = make([]message, 0)
	f.mu.Unlock()

	if -1 == strings.Index(r.UserAgent(), "FirePHP") && "" != r.Header.Get("X-FirePHP-Version") {
		return
	}

	headerCount := newHeaderCounter()

	headers := w.Header()
	headers.Set(`X-Wf-Protocol-1`, `http://meta.wildfirehq.org/Protocol/JsonStream/0.2`)
	headers.Set(`X-Wf-1-Plugin-1`, `http://meta.firephp.org/Wildfire/Plugin/FirePHP/Library-FirePHPCore/0.3`)
	headers.Set(`X-Wf-1-Structure-1`, `http://meta.firephp.org/Wildfire/Structure/FirePHP/FirebugConsole/0.1`)

	for _, v := range messages {
		msgType := &struct {
			Type  string
			Label string
		}{
			Type:  v.t.String(),
			Label: v.label,
		}
		response := []interface{}{msgType, v.content}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			continue
		}

		lenResponse := len(responseBytes)
		log.Printf("%s %v %v", responseBytes, lenResponse, chunkSize)
		if lenResponse < chunkSize {
			headers.Set(headerCount.generate(), strconv.Itoa(lenResponse)+`|`+string(responseBytes)+`|`)
			continue
		}

		buf := bytes.NewBuffer(responseBytes)
		chunk := buf.Next(chunkSize)
		headers.Set(headerCount.generate(), strconv.Itoa(lenResponse)+`|`+string(chunk)+`|\`)
		for {
			chunk := buf.Next(chunkSize)
			if len(chunk) == 0 {
				break
			}
			body := `|` + string(chunk) + `|`
			if len(chunk) == chunkSize {
				body = body + `\`
			}
			headers.Set(headerCount.generate(), body)
		}
	}
}

type headCounter int

func newHeaderCounter() *headCounter {
	return new(headCounter)
}

func (h *headCounter) generate() string {
	*h++
	return fmt.Sprintf("X-Wf-1-1-1-%d", *h)
}
