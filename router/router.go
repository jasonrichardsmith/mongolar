package router

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/context"
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type HostSwitch map[string]http.Handler

//var Switch HostSwitch
var Switch = make(HostSwitch)

var Routers = make(map[string]*httprouter.Router)
var Aliases map[string]string
var Sites map[string]*site.SiteConfig
var HandlersMap map[string]httprouter.Handle

func Build(s map[string]*site.SiteConfig, a map[string]string, hm map[string]httprouter.Handle) HostSwitch {
	Sites = s
	Aliases = a
	HandlersMap = hm
	buildRouters()
	buildHostSwitch()
	return Switch
}

func buildRouters() {
	for k, s := range Sites {
		buildRouter(s, k)
	}
}

func buildRouter(s *site.SiteConfig, key string) {
	router := httprouter.New()
	for path, handle := range HandlersMap {
		router.GET(path, wrapHandler(handle, s))
	}
	Routers[key] = router
}

func buildHostSwitch() {
	for domain, key := range Aliases {
		Switch[domain] = Routers[key]
	}
}

func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403) // Or Redirect?
	}
}

func wrapHandler(h httprouter.Handle, s *site.SiteConfig) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "", s)
		h(w, r, ps)
	}
}
