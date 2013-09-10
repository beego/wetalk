{{define "footer"}}
<footer id="footer">
	<div class="container footer-wrap">
		<p>
			<a href="/about"><strong>{{i18n .Lang "About"}}</strong></a>
			|
			<a href="/faq"><strong>FAQ</strong></a>
			|
			<a target="_blank" href="https://github.com/beego/beebbs" target="_blank"><strong>GitHub</strong></a>
		</p>
		<p>{{i18n .Lang "Copyright"}} © 2013 {{i18n .Lang "Beego Community"}}</p>
		<p class="desc">
			{{i18n .Lang "As an open source project, contribute is welcome!"}}
			<br>
			{{i18n .Lang "Based on"}} <a target="_blank" href="http://getbootstrap.com/">Bootstrap</a>. {{i18n .Lang "Icons from"}} <a target="_blank" href="http://fortawesome.github.io/Font-Awesome/">Font Awesome</a>.
		</p>
	    <div class="btn-group">
		    <button type="button" class="btn btn-default btn-xs dropdown-toggle" data-toggle="dropdown">{{i18n .Lang "Language"}}: {{.CurLang}} <i class="caret"></i></button>
		    <ul class="dropdown-menu">
			{{$keyword := .Keyword}}
		    	{{range .RestLangs}}
		    		<li><a href="javascript::" data-lang="{{.Lang}}" class="lang-changed">{{.Name}}</a></li>
		    	{{end}}
		    </ul>
	    </div>
	</div>
</footer>
{{end}}