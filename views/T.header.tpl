{{define "header"}}<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8" />
		<link rel="shortcut icon" href="/static/img/favicon.png" />
		<meta name="author" content="Unknown, slene" />
		{{if .IsHome}}
		<meta name="description" content="{{.AppDescription}}" />
		<meta name="keywords" content="{{.AppKeywords}}">
		{{end}}

		<!-- Stylesheets -->
		<link href="/static/css/bootstrap.min.css" rel="stylesheet" />
		<link href="/static/css/bootstrap-theme.min.css" rel="stylesheet" />
		<link href="/static/css/font-awesome.min.css" rel="stylesheet" />
		<!--[if IE 7]>
		<link href="/static/css/font-awesome-ie7.min.css" rel="stylesheet" />
		<![endif]-->
		<link href="/static/css/main.css" rel="stylesheet" />
		<!-- <link href="/static/css/link.css" rel="stylesheet" /> -->
		<!-- <link href="/static/css/github.css" rel="stylesheet" /> -->

		<script type="text/javascript" href="/static/js/jquery.min.js"></script>
		<script type="text/javascript" href="/static/js/bootstrap.min.js"></script>
{{end}}