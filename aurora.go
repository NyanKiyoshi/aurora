package aurora

import (
	"fmt"
	"github.com/andersfylling/disgord"
)

func New(opts *Options) *Aurora {
	aurora := &Aurora{Options: opts}
	client, err := disgord.NewClient(opts.DisgordOptions)
	if err != nil {
		panic(fmt.Sprintf("failure to initialize disgord: %s", err.Error()))
	}

	aurora.Disgord = client
	return aurora
}

func Use(a interface{}) {
	switch t := a.(type) {
	case *Command:
		if len(t.Aliases) > 0 {
			for i := range t.Aliases {
				Commands[t.Aliases[i]] = t
			}
		}
		Commands[t.Name] = t

	case *Event:
		Events[t.Name] = t
	}
}

func (a *Aurora) Init() error {
	err := a.Connect()
	if err != nil {
		return err
	}
	a.Logger.Info(fmt.Sprintf("Loaded %d commands", len(Commands)))

	for k := range Events {
		event := Events[k]
		fmt.Printf("%t\n", event.Run(a))
		err := a.On(event.Name, event.Run(a))
		if err != nil {
			a.Logger.Error(fmt.Sprintf("Failed to load event %s: %v", event.Name, err))
		}
	}
	a.Logger.Info(fmt.Sprintf("Loaded %d events", len(Events)))
	a.DisconnectOnInterrupt()
	return nil
}

func (a *Aurora) Use(v interface{}) {
	switch t := v.(type) {
	case Logger:
		a.Logger = t
	}
}
