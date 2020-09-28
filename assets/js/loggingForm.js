$(document).ready(function () {
	$('.go-to-auth-form-bt').click(function() {
		$('.auth').toggle();
		$('.reg').toggle();
	})

	$('.go-to-reg-form-bt').click(function() {
		$('.reg').toggle();
		$('.auth').toggle();
	})

	/*$('#auth-bt').click(function() {
		e.preventDefault();
		var information = {
			Login: $(this).closest("input[name=login]").val(),
			Password: $(this).closest("input[name=password").val(),
		};

		$.ajax({
			url: "http://localhost:8181/authorisation",
			method: "POST",
			content-type: "application/json",
			data: JSON.stringify(information),
			success: (res) => {
				readyData = JSON.parse(res);
				console.log(readyData);
			},
			error: (err) => {
				console.log(err)
			},
		});

	}) */

})