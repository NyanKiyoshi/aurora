package aurora

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/event"
	"math"
	"strings"
	"time"
)

type Event struct {
	Name string                      // the event name
	Run  func(a *Aurora) interface{} // the handler called on the event
}

var Events = make(map[string]*Event)

func NewEvent(name string, handler func(a *Aurora) interface{}) *Event {
	return &Event{name, handler}
}

// Default Aurora events
// the message handler event is used for command routing
func defaultMessageHandler(a *Aurora) interface{} {
	return func(s disgord.Session, msg *disgord.MessageCreate) {
		m := msg.Message
		prefix := a.GetPrefix(m)
		// edge case but we should be handling this
		if prefix == "" {
			return
		}

		isCommand := strings.HasPrefix(m.Content, prefix)
		if !isCommand {
			return
		}
		raw := m.Content[len(prefix):]
		rawArgs := strings.Fields(raw)

		// only the prefix was sent, no command
		if raw == "" || len(rawArgs) < 1 {
			return
		}

		cmd := rawArgs[0]
		var args []string

		if len(rawArgs) > 1 {
			args = rawArgs[1:]
		}
		command, ok := Commands[cmd]
		if !ok {
			return
		}

		ctx := Context{m, a, args}
		t := time.Now()
		command.Run(ctx)
		diff := int(math.Round(float64(time.Now().Sub(t).Nanoseconds() / 1e6)))

		a.Logger.Info(fmt.Sprintf("%s used command %s (%f)", m.Author.Username, command.Name, diff))
	}
}

func init() {
	Use(NewEvent(event.MessageCreate, defaultMessageHandler))
}
