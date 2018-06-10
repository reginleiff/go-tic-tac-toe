$(function() {
  var player = 1; 
  var playerId;
  var gamePlayable = false;
  var isGameMode = false;
  var rooms = $("#rooms"); 
  var table = $("table");
  var messages = $("#messages");
  var turn = $("#turn");

  retrieveRooms();
  setInterval(function() {
    console.log("lmao");
  }, 1000);

  $('td').click(function() {
    console.log("gamePlayableState is %s", gamePlayable)
    if (!gamePlayable) {
      return;
    }

    td = $(this);
    var state = getState(td);
    if(!state) {
      var pattern = definePatternForCurrentPlayer(player);
      changeState(td, pattern);
      if(checkIfPlayerWon(table, pattern)) {
	messages.html('Player '+ player +' has won.');
	turn.html('');
	gamePlayable = false;
      } else {
	player = setNextPlayer(player);
	displayNextPlayer(turn, player);
      }
    } else {
      messages.html('This box is already checked.');
    }
  });

  $(".reset").click(function() {
    console.log("quit");
    exitGameMode();
    player = 1;
    messages.html('');
    reset(table);
    displayNextPlayer(turn, player);
    gamePlayable = false;
    retrieveRooms();
  });

  function addClickFunction(id) {
    var tag = "#" + id; 
    $(document).one("click", tag, function() {
      console.log("room id is %s", id);
      enterGameMode();
      rooms.empty();
      gamePlayable = true;
    });
    console.log("click function added for room %s", id);
  }

  function enterGameMode() {
    $("#game").removeClass("hidden");
    $("#lobby").addClass("hidden");
  }

  function exitGameMode() {
    $("#game").addClass("hidden");
    $("#lobby").removeClass("hidden");
  }

  function toggleMode() {
    var gameDiv = $("#game");
    var lobbyDiv = $("#lobby"); 
    if(isGameMode) {
      gameDiv.addClass("hidden");
      lobbyDiv.removeClass("hidden");
    } else {
      gameDiv.removeClass("hidden");
      lobbyDiv.addClass("hidden");
    }	
    isGameMode = !isGameMode;
  }

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
    if(table.find('#tile1').hasClass(pattern) && table.find('#tile2').hasClass(pattern) && table.find('#tile3').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#tile1').hasClass(pattern) && table.find('#tile4').hasClass(pattern) && table.find('#tile7').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#tile1').hasClass(pattern) && table.find('#tile5').hasClass(pattern) && table.find('#tile9').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#tile4').hasClass(pattern) && table.find('#tile5').hasClass(pattern) && table.find('#tile6').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#tile7').hasClass(pattern) && table.find('#tile8').hasClass(pattern) && table.find('#tile9').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#tile2').hasClass(pattern) && table.find('#tile5').hasClass(pattern) && table.find('#tile8').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#tile3').hasClass(pattern) && table.find('#tile6').hasClass(pattern) && table.find('#tile9').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#tile3').hasClass(pattern) && table.find('#tile5').hasClass(pattern) && table.find('#tile7').hasClass(pattern)) {
      won = 1;
    }
    return won;
  }

  function reset(table) {
    table.find('td').each(function() {
      $(this).removeClass("cross").removeClass("circle");
    });
  }

  function retrieveRooms() {
    $.ajax({
      url: "http://localhost:8080/api/get/rooms",
      async: true
    }).then(function(data) {
      console.log("retrieving rooms");
      $.each(JSON.parse(data), function(i, room) {
	var roomElem;
	var status;
	switch (room.status) 
	{
	  case 0:
	    status = "Empty";
	    break;
	  case 1:
	    status = "Partially Empty";
	    break;
	  case 2:
	    status = "Full";
	    break;
	  default:
	    console.log("Error: invalid status code")
	    break;
	}
	roomElem = ("<li><button type='button' id='" + room.id + "'> Room " + i + " Status:" + status + "</button></li>");
	console.log("room %s is retrieved", room.id);
	addClickFunction(room.id);
	rooms.append(roomElem);
      })
    });
  }
});


