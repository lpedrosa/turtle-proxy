package delay

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type matchRule struct {
	method string
	path   string
}

type gorillaMatcher struct {
	router *mux.Router
	rules  map[string]*matchRule
}

func newGorillaMatcher() *gorillaMatcher {
	return &gorillaMatcher{
		router: mux.NewRouter(),
		rules:  make(map[string]*matchRule)}
}

func matchIdFor(method string, path string) string {
	// FIXME maybe we should base64 encode the key (method+path)
	return path
}

func (g *gorillaMatcher) Add(method string, path string) string {
	g.router.Handle(path, nil).Methods(method)

	// store in the map to help removal
	matchId := matchIdFor(method, path)
	g.rules[matchId] = &matchRule{method: method, path: path}

	return matchId
}

func (g *gorillaMatcher) Match(method string, path string) bool {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		// TODO use a proper logger
		log.Printf("Error matching path: %s\n", err)
		return false
	}

	var routeMatch mux.RouteMatch
	ok := g.router.Match(req, &routeMatch)

	return ok
}

func (g *gorillaMatcher) Remove(id string) {
	_, ok := g.rules[id]

	if !ok {
		// nothing found
		return
	}

	// delete entry
	delete(g.rules, id)

	// clear the router
	g.router = mux.NewRouter()

	// add remaining rules
	for _, matchRule := range g.rules {
		method := matchRule.method
		path := matchRule.path
		g.router.Handle(path, nil).Methods(method)
	}
}

func (g *gorillaMatcher) Clear() {
	g.router = mux.NewRouter()
	g.rules = make(map[string]*matchRule)
}
