// public wrapper for some helpful stuff
package driver

import (
	"net/http"

	"github.com/google/pprof/internal/plugin"
	"github.com/google/pprof/profile"
)

func NewProfileCopier(src *profile.Profile) profileCopier {
	return makeProfileCopier(src)
}

type WebUI struct {
	web  *webInterface
	host string
	port int
}

func (w *WebUI) Handlers() map[string]http.Handler {
	panic("implement me")
}

func (w *WebUI) Handler() http.HandlerFunc {
	args := plugin.HTTPServerArgs{
		Handlers: map[string]http.Handler{
			"/":              http.HandlerFunc(w.web.dot),
			"/top":           http.HandlerFunc(w.web.top),
			"/disasm":        http.HandlerFunc(w.web.disasm),
			"/source":        http.HandlerFunc(w.web.source),
			"/peek":          http.HandlerFunc(w.web.peek),
			"/flamegraph":    http.HandlerFunc(w.web.stackView),
			"/flamegraph2":   redirectWithQuery("flamegraph", http.StatusMovedPermanently), // Keep legacy URL working.
			"/flamegraphold": redirectWithQuery("flamegraph", http.StatusMovedPermanently), // Keep legacy URL working.
			"/saveconfig":    http.HandlerFunc(w.web.saveConfig),
			"/deleteconfig":  http.HandlerFunc(w.web.deleteConfig),
			"/download": http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
				wr.Header().Set("Content-Type", "application/vnd.google.protobuf+gzip")
				wr.Header().Set("Content-Disposition", "attachment;filename=profile.pb.gz")
				w.web.prof.Write(wr)
			}),
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := args.Handlers[r.URL.Path]
		if h == nil {
			panic("unexpected path " + r.URL.Path)
			h = http.DefaultServeMux
		}
		h.ServeHTTP(w, r)
	})

	return handler

}

func NewWebUI(p *profile.Profile) (*WebUI, error) {
	copier := NewProfileCopier(p)
	web, err := makeWebInterface(p, copier, &plugin.Options{})
	if err != nil {
		return nil, err
	}

	for n, c := range pprofCommands {
		web.help[n] = c.description
	}
	for n, help := range configHelp {
		web.help[n] = help
	}
	web.help["details"] = "Show information about the profile and this view"
	web.help["graph"] = "Display profile as a directed graph"
	web.help["flamegraph"] = "Display profile as a flame graph"
	web.help["reset"] = "Show the entire profile"
	web.help["save_config"] = "Save current settings"
	return &WebUI{
		web: web,
	}, nil
}
