/**
* @program: go-watch
*
* @description:
*
* @author: lemo
*
* @create: 2020-06-06 22:54
**/

package app

import (
	"net/url"
	path2 "path"
	"strings"

	"github.com/Lemo-yxk/lemo/console"
	"github.com/Lemo-yxk/lemo/exception"
	"github.com/Lemo-yxk/lemo/http"
	"github.com/Lemo-yxk/lemo/http/server"
	"github.com/fsnotify/fsnotify"
)

func (w *Watch) StartServer(host string) {

	console.Bold.Println(host)

	u, err := url.Parse(host)
	exception.Assert(err)

	if u.Path == "" {
		u.Path = "/"
	}

	var httpServer = server.Server{Host: u.Host}

	// httpServer.OnError = func(stream *http.Stream, err exception.Error) {
	// 	console.Error(err)
	// }

	var router = server.Router{IgnoreCase: true}

	router.Route("GET", u.Path).Handler(func(stream *http.Stream) exception.Error {

		stream.AutoParse()

		var name = stream.Query.Get("name").String()
		if name == "" {
			name = w.listenPath
		}

		if strings.HasSuffix(name, "*") {
			name = name[0 : len(name)-1]
		}

		if !path2.IsAbs(name) {
			name = path2.Join(w.listenPath, name)
		}

		w.Task(fsnotify.Event{Name: name})

		return stream.JsonFormat("SUCCESS", 200, name)
	})

	httpServer.SetRouter(&router).Start()
}
