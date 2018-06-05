$(document).ready(function() {
	$.ajax({
		url: "http://localhost:8080/api/get/rooms"
	}).then(function(data) {
		$.each(JSON.parse(data), function(i, room) {
			var room_elem = ("<li class='square room'>" + room.status + "</li>")
			$('#rooms').append(room_elem)	
		})
	});
});
