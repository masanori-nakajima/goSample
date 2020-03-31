package main

import ( //Goで使いたいパッケージをインポート
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Todo struct { //structとは構造体を意味している。
	ID   int    `json:"id"`
	Name string `json:"name"`
	Todo string `json:"todo"`
}

func getTodo(w http.ResponseWriter, r *http.Request) { //Dbから名前とTodo呼び出してJｓに返す
	rows, err := db.Query("select * from todos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		var (
			id   int
			name string
			todo string
		)

		if err := rows.Scan(&id, &name, &todo); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		todos = append(todos, Todo{id, name, todo})
	}

	if err := json.NewEncoder(w).Encode(&todos); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func createTodo(w http.ResponseWriter, r *http.Request) { //Jsから送られてきた情報をDBに送って保管する。
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if _, err := db.Exec("insert into todos (name, todo) values (?, ?)", todo.Name, todo.Todo); err != nil { //values(?)にはtodo.Nameとtodo.Todoの値が入る
		http.Error(w, err.Error(), 500)
		return
	}
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "no id", 400)
		return
	}

	if _, err := db.Exec("delete from todos where id = ?", id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec("create table if not exists todos (id integer primary key autoincrement, name varchar(255), todo varchar(255))"); err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTodo(w, r)
		case http.MethodPost:
			createTodo(w, r)
		case http.MethodDelete:
			deleteTodo(w, r)
		}
	})
	log.Println("start http server :8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) //ここでサーバを起動
}
