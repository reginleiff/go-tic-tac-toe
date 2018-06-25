$(function() {
  var player = 1; 
  var playerId = 1; //remember to set as cookies
  var opponentId = 2; //set to opponent in room
  var boardId;
  var currentRoomId;
  var gameUpdateHandler;
  var lobbyUpdateHandler;
  var gamePlayable = false;
  var isGameMode = false;
  var rooms = $("#rooms"); 
  var table = $("table");
  var messages = $("#messages");
  var turn = $("#turn");


  enableLobbyUpdateHandler();
  getPlayerId();

  $('td').click(function() {
    console.log("gamePlayableState is %s", gamePlayable)
    if (!gamePlayable) {
      return;
    }

    td = $(this);
    var state = getState(td);
    if(!state) {
      changeState(td, "cross");
      updateDbGameState(td.attr("id"), playerId);
    } else {
      messages.html("This box is already checked.");
    }
  });

  $(".reset").click(function() {
    console.log("quit");
    exitGameMode();
    gameEnd(null);
    boardId = null;
    updateDbRoomExit(currentRoomId);
    currentRoomId = null;
    messages.html("");
    reset(table);
    displayNextPlayer(turn, player);
    gamePlayable = false;
  });

  $(".clear").click(function() {
    gamePlayable = false;
    clearBoardState();
  });

  function addClickFunction(roomId) {
    var tag = "#" + roomId; 
    $(document).one("click", tag, function() {
      enterGameMode();
      rooms.empty();
      gamePlayable = true;
      boardId = $(this).attr("class");
      currentRoomId = $(this).attr("id");
      console.log("boardId set to %s, roomId set to %s", boardId, currentRoomId);
      disableLobbyUpdateHandler();
      enableGameUpdateHandler();
      console.log("HELLO the ROOM ID IS %s", $(this).attr("id"));
      updateDbRoomEnter($(this).attr("id"));
    });
    //console.log("click function added for room %s", roomId);
  }

  function enterGameMode() {
    $("#game").removeClass("hidden");
    $("#lobby").addClass("hidden");
  }

  function exitGameMode() {
    $("#game").addClass("hidden");
    $("#lobby").removeClass("hidden");
  }

  function getState(td) {
    console.log("tile has class %s\n", td.attr("class")) 
    if(td.hasClass("cross") || td.hasClass("circle")) {
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
      return "cross";
    } else {
      return "circle";
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

  function checkIfPlayerWon(pattern) {
    var won = 0;
    if(table.find('#1').hasClass(pattern) && table.find('#2').hasClass(pattern) && table.find('#3').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#1').hasClass(pattern) && table.find('#4').hasClass(pattern) && table.find('#7').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#1').hasClass(pattern) && table.find('#5').hasClass(pattern) && table.find('#9').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#4').hasClass(pattern) && table.find('#5').hasClass(pattern) && table.find('#6').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#7').hasClass(pattern) && table.find('#8').hasClass(pattern) && table.find('#9').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#2').hasClass(pattern) && table.find('#5').hasClass(pattern) && table.find('#8').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#3').hasClass(pattern) && table.find('#6').hasClass(pattern) && table.find('#9').hasClass(pattern)) {
      won = 1;
    } else if (table.find('#3').hasClass(pattern) && table.find('#5').hasClass(pattern) && table.find('#7').hasClass(pattern)) {
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
    rooms.empty();
    $.ajax({
      url: "http://149.28.144.110:3000/api/get/rooms",
      async: true
    }).then(function(data) {
      //console.log("retrieving rooms");
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
	// board id stored as class, room id stored as id
	roomElem 
	  = ("<li><button type='button' class='" + room.board_id + "' id='" + room.id + "'> Room " + room.id + " Status: " + status + "</button></li>");
	//console.log("room %s retrieved, corresponding board is %s", room.id, room.board_id);
	addClickFunction(room.id);
	rooms.append(roomElem);
      })
    });
  }

  function retrieveTileStates(boardId) {
    var tilesUrlQuery = "http://149.28.144.110:3000/api/get/tiles/?boardid=" + boardId;
    var movesMade = 0;
    $.ajax({
      url: tilesUrlQuery,
      async: true
    }).then(function(data) {
      $.each(JSON.parse(data), function(i, tile) { 
	var tileElem; 
	//console.log("tile retrieved with id - %s, gametile - %s, playerId - %s", tile.id, tile.game_tile, tile.player_id); 
	if(tile.player_id == playerId) { 
	  setTileState(tile.game_tile, 1); 
	  movesMade++;
	} else if (tile.player_id == null) { 
	  setTileState(tile.game_tile, 0); 
	} else {
	  setTileState(tile.game_tile, 2);
	  movesMade++;
	}
      })
    })

    if(checkIfPlayerWon("cross")) {
      gameEnd(playerId);
    } else if (checkIfPlayerWon("circle")) {
      gameEnd(opponentId);	
    } else if (movesMade == 9) {
      gameEnd(null);
    }
  }

  function gameStart() {
    disableLobbyUpdateHandler();
    enableGameUpdateHandler();
    // tell first player to start
  }

  function gameEnd(winner) {
    gamePlayable = false;

    if(winner == null) {
      messages.html("It's a draw!");
    } else {
      messages.html("Player " + winner + " has won!");
    }

    disableGameUpdateHandler();
    enableLobbyUpdateHandler();
    turn.html("");
  }

  function setTileState(gameTile, status) {
    var targetTile = "#" + gameTile;
    switch(status) {
      case 0: break;
	//console.log("tile %s is not occupied", gameTile);
      case 1: $(targetTile).addClass("cross");
	//console.log("tile %s set to cross", gameTile);
	break;
      case 2: $(targetTile).addClass("circle");
	//console.log("tile %s set to circle", gameTile);
	break;
      default:
	break;
    }
  }

  function getPlayerId() {
    // TODO: set player id to his cookie player id
    return;
  }

  function enableLobbyUpdateHandler() {
    retrieveRooms() //do 1 time first
    if(lobbyUpdateHandler) {
      disableLobbyUpdateHandler();
    }

    if(gameUpdateHandler) {
      disableGameUpdateHandler();
    }

    lobbyUpdateHandler = setInterval(function() {
      retrieveRooms();
    }, 5000);
    return;
  }

  function disableLobbyUpdateHandler() {
    clearInterval(lobbyUpdateHandler);
  }

  function enableGameUpdateHandler() {
    if (lobbyUpdateHandler) {
      disableLobbyUpdateHandler();
    }

    if (gameUpdateHandler) {
      disableGameUpdateHandler();
    }

    gameUpdateHandler = setInterval(function() {
      retrieveTileStates(boardId)
    }, 500);
    return;
  }

  function disableGameUpdateHandler() {
    console.log("Game update handler disabled!");
    clearInterval(gameUpdateHandler);
    return;
  }

  function updateDbGameState(tileId, playerId) {
    var updateGameTileQuery = "http://149.28.144.110:3000/api/put/tiles/?tileid=" + tileId + "&playerid=" + playerId;
    $.ajax({
      url: updateGameTileQuery
    }).then(function(res) {
      console.log(res);
    }).catch(function(err) {
      console.log(err);
    })
    return;
  }

  function clearBoardState() {
    if (boardId == null) {
      return;
    }

    var clearBoardStateQuery = "http://149.28.144.110:3000/api/put/boards/?boardid=" + boardId;
    $.ajax({
      url: clearBoardStateQuery
    }).then(function(res) {
      console.log(res);
    }).catch(function(err) {
      console.log(err);
    })
    return;
  }

  function updateDbRoomEnter(roomId) {
    var updatePlayersQuery = "http://149.28.144.110:3000/api/put/players/?playerid=" + playerId + "&roomid=" + roomId;
    var getPlayersInRoomQuery = "http://149.28.144.110:3000/api/get/players/?roomid=" + roomId;
    var newStatus;

    $.ajax({
      url: getPlayersInRoomQuery,
      async: true
    }).then(function(data) {
      players = JSON.parse(data);
      if (players == null) {
	console.log("number of players in room is 0");
	newStatus = 1;
      } else {
	console.log("number of players in room is %s", players.length);
	switch(players.length) {
	  case 1:
	    newStatus = 2;
	    break;
	  default:
	    retrieveRooms();
	    return;
	}
      }


      var updateRoomStatusQuery = "http://149.28.144.110:3000/api/put/rooms/?roomid=" + roomId + "&status=" + newStatus; 

      $.ajax({
	url: updateRoomStatusQuery
      }).then(function(res) {
	console.log(res);
      }).catch(function(err) {
	console.log(err);
      })

      $.ajax({
	url: updatePlayersQuery
      }).then(function(res) {
	console.log(res);
      }).catch(function(err) {
	console.log(err);
      })
    })

    return;
  }

  function updateDbRoomExit(roomId) {
    var updatePlayersQuery = "http://149.28.144.110:3000/api/put/players/?playerid=" + playerId + "&roomid=NULL";
    var getPlayersInRoomQuery = "http://149.28.144.110:3000/api/get/players/?roomid=" + roomId;
    var newStatus;

    $.ajax({
      url: getPlayersInRoomQuery,
      async: true
    }).then(function(data) {
      players = JSON.parse(data);
      if (players == null) {
	console.log("number of players in room is 0");
	newStatus = 0;
      } else {
	console.log("number of players in room is %s", players.length);
	switch(players.length) {
	  case 1:
	    newStatus = 0;
	    break;
	  case 2:
	    newStatus = 1;
	    return;
	}
      }

      var updateRoomStatusQuery = "http://149.28.144.110:3000/api/put/rooms/?roomid=" + roomId + "&status=" + newStatus; 

      $.ajax({
	url: updateRoomStatusQuery
      }).then(function(res) {
	console.log(res);
      }).catch(function(err) {
	console.log(err);
      })

      $.ajax({
	url: updatePlayersQuery
      }).then(function(res) {
	console.log(res);
      }).catch(function(err) {
	console.log(err);
      })
    })

    return;  
  }

});


