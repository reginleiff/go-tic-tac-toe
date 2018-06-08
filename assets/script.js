$(document).ready(function() {
  var player = 1; 
  var table = $('table');
  var messages = $('.messages');
  var turn = $('.turn');

  $.ajax({
    url: "http://localhost:8080/api/get/rooms"
  }).then(function(data) {
    $.each(JSON.parse(data), function(i, room) {
      var room_elem;

      switch (room.status) 
      {
	case 0:
	  room_elem = ("<button type=`button` class='list-group-item'> Room " + i + " Status: Empty</li>");
	  break;
	case 1:
	  room_elem = ("<button type=`button` class='list-group-item'> Room " + i + " Status: Partially Full</li>");
	  break;
	case 2:
	  room_elem = ("<button type=`button` class='list-group-item'> Room " + i + " Status: Full</li>");
	  break;
	default:
	  break;
      }

      $('#rooms').append(room_elem);
    })
  });

  $('td').click(function() {
    td = $(this);
    var state = getState(td);
    if(!state) {
      var pattern = definePatternForCurrentPlayer(player);
      changeState(td, pattern);
      if(checkIfPlayerWon(table, pattern)) {
	messages.html('Player '+ player +' has won.');
	turn.html('');
      } else {
	player = setNextPlayer(player);
	displayNextPlayer(turn, player);
      }
    } else {
      messages.html('This box is already checked.');
    }
  });

  /* $('.reset').click(function() {
    player = 1;
    messages.html('');
    reset(table);
    displayNextPlayer(turn, player);
  }); */

});



function getState(td) {
  console.log("tile has class %s\n", td.attr('class')) 
  if(td.hasClass('cross') || td.hasClass('circle')) {
    return 1;
  } else {
    return 0;
  }
}

function changeState(td, pattern) {
  td.html = pattern; 
  return td.addClass(pattern);
}

function definePatternForCurrentPlayer(player) {
  if(player == 1) {
    return 'cross';
  } else {
    return 'circle';
  }
}

function setNextPlayer(player) {
  if(player == 1) {
    return player = 2;
  } else {
    return player = 1;
  }
}

function displayNextPlayer(turn, player) {
  turn.html('Player turn : '+player);
}

function checkIfPlayerWon(table, pattern) {
  var won = 0;
  if(table.find('.tile1').hasClass(pattern) && table.find('.tile2').hasClass(pattern) && table.find('.tile3').hasClass(pattern)) {
    won = 1;
  } else if (table.find('.tile1').hasClass(pattern) && table.find('.tile4').hasClass(pattern) && table.find('.tile7').hasClass(pattern)) {
    won = 1;
  } else if (table.find('.tile1').hasClass(pattern) && table.find('.tile5').hasClass(pattern) && table.find('.tile9').hasClass(pattern)) {
    won = 1;
  } else if (table.find('.tile4').hasClass(pattern) && table.find('.tile5').hasClass(pattern) && table.find('.tile6').hasClass(pattern)) {
    won = 1;
  } else if (table.find('.tile7').hasClass(pattern) && table.find('.tile8').hasClass(pattern) && table.find('.tile9').hasClass(pattern)) {
    won = 1;
  } else if (table.find('.tile2').hasClass(pattern) && table.find('.tile5').hasClass(pattern) && table.find('.tile8').hasClass(pattern)) {
    won = 1;
  } else if (table.find('.tile3').hasClass(pattern) && table.find('.tile6').hasClass(pattern) && table.find('.tile9').hasClass(pattern)) {
    won = 1;
  } else if (table.find('.tile3').hasClass(pattern) && table.find('.tile5').hasClass(pattern) && table.find('.tile7').hasClass(pattern)) {
    won = 1;
  }
  return won;
}
function reset(table) {
  table.find('td').each(function() {
    $(this).removeCross().removeCircle();
  });
}
