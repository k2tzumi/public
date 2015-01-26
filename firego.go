//go:generate stringer -type=MessageType
package firego

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type message struct {
	t       MessageType
	content string
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

func (f *FireGo) Message(t MessageType, content string) {
	msg := message{
		t:       t,
		content: content,
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

func (f *FireGo) Flush(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	messages := f.messages
	f.messages = make([]message, 0)
	f.mu.Unlock()

	if -1 == strings.Index(r.UserAgent(), "FirePHP") && "" != r.Header.Get("X-FirePHP-Version") {
		return
	}

	headers := w.Header()
	headers.Set(`X-Wf-Protocol-1`, `http://meta.wildfirehq.org/Protocol/JsonStream/0.2`)
	headers.Set(`X-Wf-1-Plugin-1`, `http://meta.firephp.org/Wildfire/Plugin/FirePHP/Library-FirePHPCore/0.3`)
	headers.Set(`X-Wf-1-Structure-1`, `http://meta.firephp.org/Wildfire/Structure/FirePHP/FirebugConsole/0.1`)

	headerCount := 1
	for _, v := range messages {
		header := fmt.Sprintf("X-Wf-1-1-1-%d", headerCount)

		msgType := &struct{ Type string }{Type: v.t.String()}
		response := []interface{}{msgType, v.content}

		responseBytes, _ := json.Marshal(response)
		finalJson := string(responseBytes)

		headers.Set(
			header,
			strconv.Itoa(len(finalJson))+`|`+finalJson+`|`,
		)

		headerCount++
	}
}
