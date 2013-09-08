{{define "navbar"}}
<div class="navbar navbar-static-top">
    <div class="navbar-inner navbar-fixed-top ">
        <div class="container">
        	<div class="brand">
	            <a class="logo" href="/">
	            	<img src="/static/img/bee.gif" style="height: 60px;">
	            </a>
        	</div>
            <div style="padding-left: 80px;">
                <form class="navbar-search" action="/">
                    <input id="navbar_search_box" class="search-query" type="text" placeholder="{{i18n .Lang "Search"}}" name="q">
                </form>
            </div>

            <ul class="nav pull-right">
                <li {{if .IsHome}}class="active"{{end}}><a href="/">{{i18n .Lang "Home"}}</a></li>
                <li {{if .IsResource}}class="active"{{end}}><a href="/resource">{{i18n .Lang "Resource"}}</a></li>
                <li><a target="_blank" href="http://beego.me">{{i18n .Lang "Official website"}}</a></li>
                {{if .IsLogin}}
                {{else}}
                <li><a href="/login">{{i18n .Lang "Login"}}</a></li>
                {{end}}
            </ul>
        </div>
    </div>
</div>
{{end}}