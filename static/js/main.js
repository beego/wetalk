(function($){

	// Avoid embed thie site in an iframe of other WebSite
	if (top.location != location) {
		top.location.href = location.href;
	}

	(function(){
		// extend jQuery ajax, set xsrf token value
	    var ajax = $.ajax;
	    $.extend({
	        ajax: function(url, options) {
	            if (typeof url === 'object') {
	                options = url;
	                url = undefined;
	            }
	            options = options || {};
	            url = options.url;
	            var xsrftoken = $('meta[name=_xsrf]').attr('content');
	            var headers = options.headers || {};
	            var domain = document.domain.replace(/\./ig, '\\.');
	            if (!/^(http:|https:).*/.test(url) || eval('/^(http:|https:)\\/\\/(.+\\.)*' + domain + '.*/').test(url)) {
	                headers = $.extend(headers, {'X-Xsrftoken':xsrftoken});
	            }
	            options.headers = headers;
	            var callback = options.success;
	            options.success = function(data){
	            	if(data.once){
	            		// change all _once value if ajax data.once exist
	            		$('[name=_once]').val(data.once);
	            	}
	            	if(callback){
	            		callback(data);
	            	}
	            }
	            return ajax(url, options);
	        }
        });

	    // shake a container box
	    $.fn.shake = function (options) {
	        // defaults
	        var settings = {
	            'shakes': 2,
	            'distance': 10,
	            'duration': 400
	        };
	        // merge options
	        if (options) {
	            $.extend(settings, options);
	        }
	        // make it so
	        var pos, shakes = settings.shakes, distance = settings.distance, duration = settings.duration;
	        return this.each(function () {
	            var $self = $(this), direction = 'left';
	            // position if necessary
	            pos = $self.css('position');
	            if (!pos || pos === 'static') {
	                $self.css('position', 'relative');
	            }else if(pos == 'absolute'){
	                if($self.css('left') == 'auto'){
	                    direction = 'right';
	                }
	            }
	            // shake it
	            for (var x = 1; x <= shakes; x++) {
	                var e1 = {}, e2 = {}, e3 = {};
	                e1[direction] = distance * -1;
	                e2[direction] = distance;
	                e3[direction] = 0;
	                $self.animate(e1, (duration / shakes) / 4)
	                    .animate(e2, (duration / shakes) / 2)
	                    .animate(e3, (duration / shakes) / 4);
	            }
	        });
	    };

	})();

	// btn checked box toggle
	$(document).on('click', '.btn-checked', function(){
		var $e = $(this);
		var $i = $e.siblings('[name='+$e.data('name')+']');
		if($e.hasClass('active')) {
			$i.val('true');
		} else {
			$i.val('false');
		}
	});

	// change locale and reload page
	$(document).on('click', '.lang-changed', function(){
		var $e = $(this);
		var lang = $e.data('lang');
		Cookies.set('lang', lang);
		window.location.reload();
	});

	(function(){
		var unload = false;
		var submited = 'submited';

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
		$(document).on('submit', 'form', function(){
			var $e = $(this);
			if($e.data(submited)){
				return false;
			}
			$e.data(submited, true);
			unload = true
		});
	})()

	// for ajax dropdown login
	$(document).on('submit', '#dropdown-login', function(){
		var $form = $(this);
	    var $alert = $form.find('.alert');
	    var url = $form.attr('action');
	    var data = $form.find('input').fieldSerialize();
	    if($.trim($form.find("[name=UserName]").val()) == '' 
	    	|| $.trim($form.find("[name=Password]").val()) == '') {
            $form.shake();
	    	return false;
	    }
	    $.post(url, data, function(data){
	        $alert.removeClass('alert-info alert-danger');
            $alert.text(data.message);
	        if(data.success){
	            $alert.addClass('alert-success');
	            setTimeout(function(){
	            	window.location.reload();
	            });
	        }else{
	            $alert.addClass('alert-danger');
	            $form.shake();
	        }
	    });
	    return false;
	});

	$(function(){
		// on dom ready

	    $('[data-show=tooltip]').each(function(k, e){
	        var $e = $(e);
	        $e.tooltip({placement: $e.data('placement'), title: $e.data('tooltip-text')});
	        $e.tooltip('show');
	    });
	});

})(jQuery);