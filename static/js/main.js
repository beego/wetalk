(function($){

	(function(){
	    var ajax = $.ajax;
	    $.extend({
	        ajax: function(url, options) {
	            if (typeof url === "object") {
	                options = url;
	                url = undefined;
	            }
	            options = options || {};
	            url = options.url;
	            var xsrftoken = Cookies.get('_xsrf');
	            var data = options.data || {};
	            var domain = document.domain.replace(/\./ig, '\\.');
	            if (!/^(http:|https:).*/.test(url) || eval('/^(http:|https:)\\/\\/(.+\\.)*' + domain + '.*/').test(url)) {
	                data = $.extend(data, {'_xsrf':xsrftoken});
	            }
	            options.data = data;
	            return ajax(url, options);
	        }
        });
	})();

	// btn checked box toggle
	$(document).on("click", ".btn-checked", function(){
		var $e = $(this);
		var $i = $e.siblings("[name="+$e.data("name")+"]");
		if($e.hasClass("active")) {
			$i.val("true");
		} else {
			$i.val("false");
		}
	});

	// change locale and reload page
	$(document).on("click", ".lang-changed", function(){
		var $e = $(this);
		var lang = $e.data("lang");
		Cookies.set("lang", lang);
		window.location.reload();
	});

	(function(){
		var unload = false;
		var submited = "submited";

		$(window).unload(function(){
			if(unload){
				// skip first unload
				unload = false;
				return;
			}
			$('form').each(function(k, e){
				var $e = $(e);
				if($e.data(submited)){
					$.data(submited, false);
				}
			});
		});

		// avoid form re-submit
		$(document).on("submit", "form", function(){
			var $e = $(this);
			if($e.data(submited)){
				return false;
			}
			$e.data(submited, true);
			unload = true
		});
	})()

	$(function(){
		// on dom ready

	    $('[data-show=tooltip]').each(function(k, e){
	        var $e = $(e);
	        $e.tooltip({placement: $e.data('placement'), title: $e.data('tooltip-text')});
	        $e.tooltip('show');
	    });
	});

})(jQuery);