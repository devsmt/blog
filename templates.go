package main

import (
	"fmt"
	"text/template"
)

var (
	base = `<!DOCTYPE html>
<html>
	<head>
		<style>
body {
        font-family: helvetica, verdana, sans-serif;
        background-color: #eee;
        font-size: 1em;
        color: #333;
        margin: 0;
}

.main {
}

.content {
        max-width: 1000px;
        min-width: 500px;
        margin: auto;
        line-height: 1.5;
}

p, ul, ol {
        font-size: 1.25em;
        line-height: 1.7em;
        font-family: "Droid Serif", serif;
}

article {
        background-color: #fff;
        border: 1px solid #ccc;
        border-radius: 3px;
        padding: 2em 2em 0 2em;
        margin: 0 1em 1em 1em;
}

a, u {
        text-decoration: none;
        color: rgb(22, 85, 126);
}

a:hover {
        text-decoration: underline;
}

.main-header a:visited {
        color: #fff;
        font-weight: 200;
}

.main-header a:hover {
        color: #ccc;
        text-decoration: none;
}

@font-face {
        font-family: "Lato";
        font-style: normal;
        font-weight: 400;
        src: local("Lato Bold");
}

.main-header {
        margin: 0;
        margin-bottom: 20px;
        padding: 20px;
        font-family:Lato, sans-serif;
        font-size: 2em;
        background-color: rgb(22, 85, 126);
        color: #fff;
}

p {
        margin: 0;
        margin-bottom: 1em;
}

h1, h2 {
        margin: 0 0 .25em 0;
        font-size: 2em;
}

h2 a {
        color: black;
}

h2 a:hover {
        color: rgb(22, 85, 126);
        text-decoration: none;
}

		</style>
		<script>
		  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
		  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
		  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
		  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

		  ga('create', 'UA-46738088-1', 'weberc2.com');
		  ga('send', 'pageview');

		</script>
		<link rel="stylesheet" type="text/css" href="//fonts.googleapis.com/css?family=Lato:300,400,700,900" />
		<link rel="stylesheet" type="text/css" href="//fonts.googleapis.com/css?family=Droid+Serif:400,700,400italic,700italic" />
	</head>
	<body>
		<div class="main">
			<h1 class="main-header"><a href="/">CRAIG WEBER</a></h1>
			<div class="content">
			%s
			</div>
		</div>
	</body>
</html>`
)

func templify(content string) *template.Template {
	content = fmt.Sprintf(base, content)
	return template.Must(template.New("").Parse(content))
}

var (
	HOME_TEMPLATE = templify(`{{range .}}<article><h2><a href="{{.Path}}">{{.Metadata.Title}}</a></h2>{{.Snippet}}<p><a class="more" href="{{.Path}}">Read more...</a></p></article>{{end}}`)
	DOC_TEMPLATE  = templify("<article><h2>{{.Metadata.Title}}</h2>{{.Text}}</article>")
)
