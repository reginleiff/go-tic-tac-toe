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
	
	var rooms []models.Room
	var err error

	if rooms, err = queryRooms(); err != nil {
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

func queryRooms() ([]models.Room, error) {
	rooms := []models.Room{}
	rows, err := db.Query(queryGetRooms)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

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
			return nil, err
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}

func getTiles(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("debug (getTiles): getting tiles\n")

	var tiles []models.Tile
	var err error
	
	q, ok := r.URL.Query()["boardid"]
	
	if !ok || len(q) < 1 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	boardId := q[0]

	if tiles, err = queryTiles(boardId); err != nil {
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
	http.HandleFunc("/api/put/tiles/", updateTile)	

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("debug (main): failed to listen on port 8080\n")
		panic(err)
	}
}

