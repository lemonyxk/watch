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
	"github.com/Lemo-yxk/lemo/utils"
	"github.com/fsnotify/fsnotify"
)

func (w *Watch) StartServer(host string) {

	console.Bold.Println(host)

	u, err := url.Parse(host)
	if err != nil {
		console.FgRed.Println(err)
	}

	var ip = strings.Split(u.Host, ":")[0]
	var port = u.Port()
	var path = u.Path

	if port == "" {
		port = "80"
	}

	if path == "" {
		path = "/"
	}

	var httpServer = server.Server{Host: ip, Port: utils.Conv.Atoi(port)}

	// httpServer.OnError = func(stream *http.Stream, err exception.Error) {
	// 	console.Error(err)
	// }

	var router = server.Router{IgnoreCase: true}

	router.Route("GET", path).Handler(func(stream *http.Stream) exception.Error {

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
