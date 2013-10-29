(function($){

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
				var oncetoken = $('[name=_once]').filter(':last').val();
				var headers = options.headers || {};
				var domain = document.domain.replace(/\./ig, '\\.');
				if (!/^(http:|https:).*/.test(url) || eval('/^(http:|https:)\\/\\/(.+\\.)*' + domain + '.*/').test(url)) {
					headers = $.extend(headers, {'X-Xsrftoken':xsrftoken, 'X-Form-Once':oncetoken});
				}
				options.headers = headers;
				var callback = options.success;
				options.success = function(data){
					if(data.once){
						// change all _once value if ajax data.once exist
						$('[name=_once]').val(data.once);
					}
					if(callback){
						callback.apply(this, arguments);
					}
				};
				return ajax(url, options);
			}
		});

	})();
})(jQuery);