(function($){

	$(function(){
		$('[rel=select2-admin-model]').each(function(_,e){
			var $e = $(e);
			var model = $e.data('model');
			$e.select2({
				minimumInputLength: 3,
				ajax: {
					url: '/admin/model/select',
					type: 'POST',
					data: function (query, page) {
						return {
							'search': query,
							'model': model
						};
					},
					results: function(d){
						var results = [];
						if(d.success && d.data){
							var data = d.data;
							$.each(data, function(i,v){
								results.push({
									'id': v[0],
									'text': v[1]
								});
							});
						}
						return {'results': results};
					}
				},
				initSelection: function(elm, cbk){
					var id = parseInt($e.val(), 10);
					if(id){
						$.post('/admin/model/get', {'id': id, 'model': model}, function(d){
							if(d.success){
								if(d.data && d.data.length){
									cbk({
										'id': d.data[0],
										'text': d.data[1]
									});
								}
							}
						});
					}
				}
			});
		});
	});

})(jQuery);