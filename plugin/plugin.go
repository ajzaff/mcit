package plugin

import (
	"fmt"
	"plugin"

	"github.com/ajzaff/mcts"
)

// LoadSearchFunc loads a search function named "Search" from the plugin source.
//
// Compile a plugin using the `go build -buildmode=plugin` feature of Go
// (https://pkg.go.dev/plugin.) Only supported on Linux, FreeBSD, and Mac.
func LoadSearchFunc(path string) (mcts.Func, error) {
	x, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	sym, err := x.Lookup("Search")
	if err != nil {
		return nil, err
	}
	fn, ok := sym.(mcts.Func)
	if !ok {
		return nil, fmt.Errorf("expected Search to be %T, but found %T", mcts.Func(nil), fn)
	}
	return fn, nil
}
