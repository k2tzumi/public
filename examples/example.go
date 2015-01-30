package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dericofilho/firego"
)

type Application struct {
	*firego.FireGo
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Message(firego.Log, "Message from FireGo")
	a.Table([][]string{
		[]string{"test", "test2"},
		[]string{"test3", "test4"},
	})
	a.GroupStart("Group Name")
	a.Log("Message from FireGo in Group")
	a.GroupEnd()
	a.Flush(w, r)

	fmt.Fprintf(w, "APPa: Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	app := &Application{
		firego.New(),
	}
	http.Handle("/", app)
	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
