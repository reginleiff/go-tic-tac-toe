package main

import (
	"net/http"
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
	QUERY_GET_ROOMS = "SELECT * FROM rooms"
)

func getRooms(w http.ResponseWriter, r *http.Request) {
	
	rooms := []models.Room{}
	
	if err := queryRooms(&rooms); err != nil {
		http.Error(w, err.Error(), 500) // server error if failed to retrieve rooms
		return
	}
	
	out, err := json.Marshal(rooms)

	if err != nil {
		http.Error(w, err.Error(), 500) // server error if failed to retrieve rooms
		return
	}

	fmt.Fprintf(w, string(out))
}

func queryRooms(rooms *[]models.Room) error {
	rows, err := db.Query(QUERY_GET_ROOMS)

	if err != nil {
		return err
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
			return err
		}

		*rooms = append(*rooms, room)
	}

	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func initDB() {
	psqlInfo := fmt.Sprintf("user=%s dbname=%s host=%s port=%s sslmode=disable", DB_USER, DB_NAME, DB_HOST, DB_PORT)
	fmt.Println(psqlInfo)
	
	var err error
	db, err = sql.Open("postgres", "postgres://m012-hb@localhost/ttt_dev?sslmode=disable");

	if err != nil{
		fmt.Printf("couldn't open database\n")
		panic(err)
	} }

func main() {
	
	initDB()
	defer db.Close()

	http.HandleFunc("/api/get/rooms", getRooms)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("failed to listen on port 8080\n")
		panic(err)
	}
}

