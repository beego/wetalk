(function($){

	// Avoid embed thie site in an iframe of other WebSite
	if (top.location != location) {
		top.location.href = location.href;
	}

	(function(){

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
		$e.blur();
	});

	// change locale and reload page
	$(document).on('click', '.lang-changed', function(){
		var $e = $(this);
		var lang = $e.data('lang');
		$.cookie('lang', lang);
		window.location.reload();
	});

	// avoid re-submit
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
			unload = true;
		});
	})();

	// for ajax dropdown login
	$(document).on('submit', '#dropdown-login', function(){
		var $form = $(this);
		var $alert = $form.find('.alert');
		var url = $form.attr('action');
		var data = $form.find('input').fieldSerialize();
		if($.trim($form.find("[name=UserName]").val()) === '' ||
			$.trim($form.find("[name=Password]").val()) === '') {
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


	// file upload
	(function(){
		function uploadFileChange($e, $file, $field, flag){
			var current = $file.val();
			var last = $e.data('last') || '';
			if(current != last){
				$e.data('last', current);
				$field.val(current.replace(/.*(\\|\/)/, ''));
				if(!flag){
					$file.trigger('change');
				}
			}
		}

		$(document).on('click', '[data-dismiss=upload]', function(e){
			var $e = $(this);
			var $btn = $(e.target);
			if($btn.attr('rel') != 'button' && $btn.attr('rel') != 'filename') {
				return;
			}
			var $file = $e.find('input[type=file]');
			var $field = $e.find('[rel=filename]');
			$file.click();
			setTimeout(uploadFileChange, 0, $e, $file, $field);
		});

		$(document).on('change', '[data-dismiss=upload] input[type=file]', function(){
			var $file = $(this);
			var $e = $file.parents('[data-dismiss=upload]');
			var $field = $e.find('[rel=filename]');
			uploadFileChange($e, $file, $field, true);
		});
	}());

	(function(){
		var v = $.cookie('JsStorage');
		if(v){
			var values = v.split(':::');
			if(values.length > 1){
				$.jStorage[values[0]].apply(this, values.splice(1));
			}
			$.removeCookie('JsStorage', {path: '/'});
		}
	})();

	(function(){
		$(document).on('submit', '#navbar-search', function(){
			var q = $(this).find('[name=q]').val();
			if($.trim(q) !== ''){
				var host = window.location.hostname + ":" + window.location.port;
				var url = 'http://www.google.com/search?q=' + 'site:' + host + '/p%20' +  escape($.trim(q));
				window.open(url, "_blank");
			}
            return false;
		});
	})();

	(function(){

		$.fn.mdFilter = function(){
			var $e = $(this);
			$e.find('img').each(function(_,img){
				var $img = $(img);
				$img.addClass('img-responsive');
				var src = $img.attr('src');
				var url = src.replace(/(\/img\/.+\.)(\d+)(\.(jpg|png))/, function(_,p1,p2,p3){
					return p1 + "full" + p3;
				});
				if(url !== src){
					$img.wrap('<a target="_blank" href="'+url+'"></a>');
				}
			});

			$e.find('pre > code').parent().addClass("prettyprint linenums");
			prettyPrint();
		};

	})();

	$(function(){
		// on dom ready

		$('[data-show=tooltip]').each(function(k, e){
			var $e = $(e);
			$e.tooltip({placement: $e.data('placement'), title: $e.data('tooltip-text')});
			$e.tooltip('show');
		});

		$('[rel=select2]').select2();

		$('.markdown').mdFilter();
	});

})(jQuery);