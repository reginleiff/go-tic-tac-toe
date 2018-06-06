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

func isNewPlayer(r *http.Request) bool {
	fmt.Printf("debug (isNewPlayer): number of cookies is %v\n", len(r.Cookies()))
	return len(r.Cookies()) == 0
}

func readCookies(w http.ResponseWriter, r *http.Request) {

	for _, cookie := range r.Cookies() {
		fmt.Printf("debug (readCookies): cookie name - " + cookie.Name + ", cookie value - " + cookie.Value + "\n")
	}

	if isNewPlayer(r) {
		fmt.Printf("debug (readCookies): need to set cookie\n")
		id, err := createPlayer()
		if err != nil {
			fmt.Printf("debug (readCookies): error creating player\n")
		}
		fmt.Printf("debug (readCookies): id created is %v\n", id)
		setCookies(w, id)
	} else {
		fmt.Printf("debug (readCookies): cookies already set\n")
	}
}

func setCookies(w http.ResponseWriter, playerID int) {
	idString := strconv.Itoa(playerID)
	fmt.Printf("debug (setCookies): setting cookie value to %s\n", idString)
	cookie := http.Cookie{Name: "player_id", Value: idString}
	http.SetCookie(w, &cookie)
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

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("debug (main): failed to listen on port 8080\n")
		panic(err)
	}
}

