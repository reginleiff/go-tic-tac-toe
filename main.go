package main

import (
	"net/http"
	"strconv"
	"fmt"
	"database/sql"
	"github.com/reginleiff/go-tic-tac-toe/models"
	"encoding/json"
	_ "github.com/lib/pq"
) 

var db *sql.DB

const (
	DB_USER = "m012-hb"
	DB_NAME = "ttt_dev"
	DB_HOST = "localhost"
	DB_PORT = "5432"
)

const (
	queryGetRooms = "SELECT * FROM rooms"
)

const (
	playersTable = "players"
	roomsTable = "rooms"
	tilesTable = "tiles"	
	boardsTable = "boards"
)

const (
	//	playerIdNotProvidedError = "Bad Request: Player ID not provided"
	roomIdNotProvidedError = "Bad Request: Room ID not provided"
	//	tileIdNotProvidedError = "Bad Request: Tile ID not provided"
	boardIdNotProvidedError = "Bad Request: Board ID not provided"
	paramsNotProvidedError = "Bad Request: Multiple Parameters not provided"

	playerDoesNotExistError = "Bad Request: Player does not exist"
	roomDoesNotExistError = "Bad Request: Room does not exist"
	boardDoesNotExistError = "Bad Request: Board does not exist"
	invalidStatusCodeError = "Bad Request: Invalid status code given"	
	
	playerRoomUpdateSuccess = "Success: Room for Player has been updated"
	roomStatusUpdateSuccess = "Success: Status for Room has been updated"
)

func doesIdExist(id string, tableName string) bool {
	fmt.Printf("debug (doesIdExist): checking id->%s is in table->%s\n", id, tableName)
	var entryExists bool

	query := fmt.Sprintf("SELECT 1 FROM %s WHERE id = %s", tableName, id)
	fmt.Printf("debug (doesIdExist): query formed - %s\n", query)
	err := db.QueryRow(query).Scan(&entryExists)

	if err != nil {
		fmt.Printf("debug (doesIdExist): error checking id - %s\n", err)
	}

	return entryExists
}

func readCookies(w http.ResponseWriter, r *http.Request) {
	var playerId string

	for _, cookie := range r.Cookies() {
		fmt.Printf("debug (readCookies): cookie name - " + cookie.Name + ", cookie value - " + cookie.Value + "\n")
		if (cookie.Name == "player_id") {
			playerId = cookie.Value
		}
	}

	if !doesIdExist(playerId, playersTable) {
		fmt.Printf("debug (readCookies): need to set cookie\n")
		newPlayerId, err := createPlayer()

		if err != nil {
			fmt.Printf("debug (readCookies): error creating player\n")
		}

		fmt.Printf("debug (readCookies): id created is %v\n", newPlayerId)
		setCookies(w, newPlayerId)
	} else {
		fmt.Printf("debug (readCookies): cookies already set\n")
	}
	return
}

func setCookies(w http.ResponseWriter, playerID int) {
	idString := strconv.Itoa(playerID)
	fmt.Printf("debug (setCookies): setting cookie value to %s\n", idString)
	cookie := http.Cookie{Name: "player_id", Value: idString}
	http.SetCookie(w, &cookie)
	return
}

func createPlayer() (int, error) {
	var id int
	fmt.Printf("debug (createPlayer): creating new player\n")	
	err := db.QueryRow("INSERT INTO players (created_at, updated_at) VALUES (NOW(), NOW()) RETURNING id").Scan(&id);

	if err != nil {
		fmt.Println("debug (createPlayer): error retrieving id - %s\n", err)
		return 0, err
	}

	fmt.Printf("debug (createPlayer): id created is %v\n", id)
	return id, nil
}

func getRooms(w http.ResponseWriter, r *http.Request) {
	readCookies(w, r) // need to check cookies in lobby

	rows, err := db.Query(queryGetRooms)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)	
		return
	}

	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		room := models.Room{}
		err := rows.Scan(
			&room.ID,
			&room.BoardID,
			&room.Status,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)	
			return
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(rooms)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) 
		return
	}

	fmt.Fprintf(w, string(out))
	return
}

func getTiles(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("debug (getTiles): getting tiles\n")
	tileParams, tileParamsOk := r.URL.Query()["boardid"]

	if !tileParamsOk || len(tileParams) < 1 {
		http.Error(w, boardIdNotProvidedError, http.StatusBadRequest)
		return
	}

	boardId := tileParams[0]

	if !doesIdExist(boardId, boardsTable) {
		http.Error(w, boardDoesNotExistError, http.StatusBadRequest)
		return
	}	

	fmt.Printf("debug (queryTiles): querying tiles with board id - %s\n", boardId)
	query := fmt.Sprintf("SELECT * FROM tiles WHERE board_id = %s", boardId)
	rows, err := db.Query(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) 
		return
	}

	defer rows.Close()

	var tiles []models.Tile
	for rows.Next() {
		tile := models.Tile{}
		err := rows.Scan(
			&tile.ID,
			&tile.BoardID,
			&tile.CreatedAt,
			&tile.UpdatedAt,
			&tile.PlayerID,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tiles = append(tiles, tile)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(tiles)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) 
		return
	}

	fmt.Fprintf(w, string(out))
	return
}

func queryTiles(boardId string) ([]models.Tile, error) {
	fmt.Printf("debug (queryTiles): querying tiles with board id - %s\n", boardId)
	tiles := []models.Tile{}
	query := fmt.Sprintf("SELECT * FROM tiles WHERE board_id = %s", boardId)
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		tile := models.Tile{}
		err := rows.Scan(
			&tile.ID,
			&tile.BoardID,
			&tile.CreatedAt,
			&tile.UpdatedAt,
			&tile.PlayerID,
		)
		if err != nil {
			return nil, err
		}

		tiles = append(tiles, tile)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tiles, nil
}

func updateTile (w http.ResponseWriter, r *http.Request) {
	var err error	

	tileParams, tileParamsOk := r.URL.Query()["tileid"]	
	playerParams, playerParamsOk := r.URL.Query()["playerid"]	

	if !tileParamsOk || !playerParamsOk || len(tileParams) < 1 || len(playerParams) < 1 {
		fmt.Printf("debug (updateTile): bad parameters\n")	
		http.Error(w, err.Error(), http.StatusBadRequest)	
		return
	}

	tileId := tileParams[0]
	playerId := playerParams[0]

	if !doesIdExist(playerId, playersTable) {
		fmt.Printf("debug (updateTile): player id does not exist\n")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if doesIdExist(tileId, tilesTable) {
		query := fmt.Sprintf("UPDATE tiles SET player_id=%s WHERE id=%s", playerId, tileId)		
		_, err := db.Exec(query)

		if err != nil {
			fmt.Printf("debug (updateTile): error updating tile - %s\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "200 success", http.StatusOK)	
		return
	}

	http.Error(w, err.Error(), http.StatusBadRequest)
	return
}

func updatePlayerRoom(w http.ResponseWriter, r *http.Request) {
	playerParams, playerParamsOk := r.URL.Query()["playerid"] 
	roomParams, roomParamsOk := r.URL.Query()["roomid"]

	if !playerParamsOk || !roomParamsOk || len(playerParams) < 1 || len(roomParams) < 1 {
		http.Error(w, paramsNotProvidedError, http.StatusBadRequest)
		return
	}

	playerId := playerParams[0]
	roomId := roomParams[0]

	if !doesIdExist(playerId, playersTable) {
		http.Error(w, playerDoesNotExistError, http.StatusBadRequest)
		return
	}

	if !doesIdExist(roomId, roomsTable) {
		http.Error(w, roomDoesNotExistError, http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("UPDATE players SET room_id=%s WHERE id=%s", roomId, playerId)
	fmt.Printf("debug (updatePlayerRoom): query - %s\n", query)
	_, err := db.Exec(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, playerRoomUpdateSuccess, http.StatusOK)
	return
}

func getPlayersInRoom(w http.ResponseWriter, r *http.Request) {
	roomParams, roomParamsOk := r.URL.Query()["roomid"]

	if !roomParamsOk || len(roomParams) < 1 {
		http.Error(w, roomIdNotProvidedError, http.StatusBadRequest)
		return
	}

	roomId := roomParams[0]

	if !doesIdExist(roomId, roomsTable) {
		http.Error(w, roomDoesNotExistError, http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("SELECT * FROM players WHERE room_id=%s", roomId)	
	rows, err := db.Query(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var players []models.Player	

	for rows.Next() {
		player := models.Player{}
		err := rows.Scan(
			&player.ID,
			&player.RoomID,
			&player.CreatedAt,
			&player.UpdatedAt,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		players = append(players, player)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(players)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) 
		return
	}

	fmt.Fprintf(w, string(out))
	return
}

func updateRoomStatus(w http.ResponseWriter, r *http.Request) {
	roomParams, roomParamsOk := r.URL.Query()["roomid"]
	statusParams, statusParamsOk := r.URL.Query()["status"]

	if !roomParamsOk || !statusParamsOk || len(roomParams) < 1 || len(statusParams) < 1 {
		http.Error(w, paramsNotProvidedError, http.StatusBadRequest)
		return
	}

	roomId := roomParams[0]
	status := statusParams[0]

	if !doesIdExist(roomId, roomsTable) {
		http.Error(w, roomDoesNotExistError, http.StatusBadRequest)
		return
	}

	switch(status) {
		case "0": break
		case "1": break
		case "2": break
		default: http.Error(w, invalidStatusCodeError, http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("UPDATE rooms SET status=%s WHERE id=%s", status, roomId)
	_, err := db.Exec(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, roomStatusUpdateSuccess, http.StatusOK)
	return
}

func initDB() {
	psqlInfo := fmt.Sprintf("user=%s dbname=%s host=%s port=%s sslmode=disable", DB_USER, DB_NAME, DB_HOST, DB_PORT)
	fmt.Println(psqlInfo)

	var err error
	db, err = sql.Open("postgres", psqlInfo);

	if err != nil{
		fmt.Printf("debug (initDB): couldn't open database\n")
		panic(err)
	} 
}

func main() {

	initDB()
	defer db.Close()

	http.Handle("/", http.FileServer(http.Dir("./assets")))
	http.HandleFunc("/api/get/rooms", getRooms)
	http.HandleFunc("/api/get/tiles/", getTiles)
	http.HandleFunc("/api/get/players/", getPlayersInRoom)
	http.HandleFunc("/api/put/tiles/", updateTile)	
	http.HandleFunc("/api/put/players/", updatePlayerRoom)
	http.HandleFunc("/api/put/rooms/", updateRoomStatus)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("debug (main): failed to listen on port 8080\n")
		panic(err)
		return
	}
}

