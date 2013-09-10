{{define "bodyHeader"}}
	</head>	
	<body>
		<div id="wrapper">
			{{template "navbar" .}}
			<div id="main" class="container">
{{end}}
{{define "bodyFooter"}}
		    </div>
		    <div class="wrapper-push"></div>
		</div>
		{{template "footer" .}}
	</body>
</html>
{{end}}