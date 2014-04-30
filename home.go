package main

import (
	"html/template"
	"net/http"
)

var homeTmpl = template.Must(template.New("").Parse(`<html>
	<head>
		<meta http-equiv="Refresh" content="5" />
	</head>
	<body>
		<h1>MessageDroid</h1>
		<h2>Upcoming Services</h2>
		<div>
			<ul>
{{range .Services}}				<li><b>{{.ServiceId}}:</b> {{.Text}}</li>
{{end}}			</ul>
		</div>
	</body>
</html>`))

func home(w http.ResponseWriter, r *http.Request) {
	state.RLock()
	defer state.RUnlock()

	homeTmpl.Execute(w, state)
}

func init() {
	http.HandleFunc("/", home)
}
