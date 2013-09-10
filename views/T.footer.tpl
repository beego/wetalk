{{define "footer"}}
<footer>
	<div class="footer_main">
		<div class="links">
			<a class="dark" href="/about"><strong>{{i18n .Lang "About"}}</strong></a>
			|
			<a class="dark" href="/faq"><strong>FAQ</strong></a>
			|
			<a class="dark" target="_blank" href="https://github.com/beego/beebbs" target="_blank"><strong>GitHub</strong></a>
		</div>
		{{i18n .Lang "Copyright"}} © 2013 Beego Community <br>
		{{i18n .Lang "As an open source project, contribute is welcome!"}} <br>
		{{i18n .Lang "Based on"}} <a target="_blank" href="http://getbootstrap.com/">Bootstrap</a>. {{i18n .Lang "Icons from"}} <a target="_blank" href="http://fortawesome.github.io/Font-Awesome/">Font Awesome</a>. <br>
		<strong>{{i18n .Lang "Language"}}:</strong>

		<script type="text/javascript" src="http://cdn.staticfile.org/jquery/1.10.1/jquery.min.js"></script>
		<script type="text/javascript" src="/static/js/bootstrap.min.js"></script>
	    <div class="btn-group dropup">
		    <button class="btn dropdown-toggle" data-toggle="dropdown">{{.CurLang}} <span class="caret"></span></button>
		    <ul class="dropdown-menu">
			{{$keyword := .Keyword}}
		    	{{range .RestLangs}}
		    	<li><a href="?lang={{.Lang}}">{{.Name}}</a></li>
		    	{{end}}
		    </ul>
	    </div>
	</div>
</footer>
{{end}}