package theme

import "fmt"

type Theme struct {
	ID        string
	Name      string
	Distro    string
	Logo      string
	Primary   string
	Secondary string
	Accent    string
	Tagline   string
	Desktop   string
	Joke      bool
	Rainbow   bool
}

func All() []Theme {
	return append([]Theme(nil), themes...)
}

func ByID(id string) (Theme, bool) {
	for _, item := range themes {
		if item.ID == id {
			return item, true
		}
	}
	return Theme{}, false
}

func Must(id string) Theme {
	item, ok := ByID(id)
	if !ok {
		panic(fmt.Sprintf("unknown theme %q", id))
	}
	return item
}
