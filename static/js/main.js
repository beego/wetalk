(function($){

	$(document).on("click", ".btn-checked", function(){
		var $e = $(this);
		if($e.val() == "1") {
			$e.val(0);
		} else {
			$e.val(1);
		}
	});

	$(document).on("click", ".lang-changed", function(){
		var $e = $(this);
		var lang = $e.data("lang");
		$.cookie("lang", lang);
		window.location.reload();
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