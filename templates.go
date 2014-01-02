package main

import (
	"fmt"
	"text/template"
)

var (
	style = `
body {
	font-family: helvetica, verdana, sans-serif;
	background-color: #eee;
	font-size: 1em;
	color: #555;
	margin: 0;
}

.content {
	line-height: 1.5;
}

article {
	min-width: 500px;
	width: 45%;
	float: left;
	display: inline-block;
	background-color: #fff;
	border: 1px solid #ccc;
	border-radius: 3px;
	padding: 1px 15px;
	margin: 0 1em 1em 1em;
}

a, u {
	text-decoration: none;
}

.main-header a:visited {
	color: black;
}

.main-header {
	margin: .5em;
	padding: 0;
	font-size: 2em;
}

p {
	margin: 0;
	margin-bottom: 1em;
}

h1, h2 {
	margin: .5em 0;
}

`
	base = `<html>
	<head>
		<style>
			%s
		</style>
		<script>
		  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
		  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
		  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
		  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

		  ga('create', 'UA-46738088-1', 'weberc2.com');
		  ga('send', 'pageview');

		</script>
	</head>
	<body>
		<div class="main">
			<h1 class="main-header"><a href="/">Craig Weber's blog</a></h1>
			<div class="content">
			%s
			</div>
		</div>
	</body>
</html>`
)

func templify(content string) *template.Template {
	content = fmt.Sprintf(base, style, content)
	return template.Must(template.New("").Parse(content))
}

var (
	HOME_TEMPLATE = templify(`{{range .}}<article><h2><a href="{{.Path}}">{{.Title}}</a></h2>{{.Text}}</article>{{end}}`)
	DOC_TEMPLATE  = templify("<article>{{.}}</article>")
)
