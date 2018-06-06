$(document).ready(function() {
  $.ajax({
    url: "http://localhost:8080/api/get/rooms"
  }).then(function(data) {
    $.each(JSON.parse(data), function(i, room) {
      var room_elem

      switch (room.status) 
      {
	case 0:
	  room_elem = ("<button type=`button` class='list-group-item list-group-item-action list-group-item-success'> Room " + i + " Status: Empty</li>") 
	  break
	case 1:
	  room_elem = ("<button type=`button` class='list-group-item list-group-item-action list-group-item-danger'> Room " + i + " Status: Partially Full</li>")
	  break
	case 2:
	  room_elem = ("<button type=`button` class='list-group-item list-group-item-action list-group-item-warning'> Room " + i + " Status: Full</li>")
	  break
	default:
	  break
      }

      $('#rooms').append(room_elem)	
    })
  });
});
