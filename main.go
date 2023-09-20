package main

import (
	"encoding/json"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type User struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	PhoneNumber string `json:"phone_number"`
}

func main() {
	http.HandleFunc("/users", getUsers)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	// Установите параметры подключения к вашей базе данных PostgreSQL
	connStr := "user=postgres password=password host=localhost port=5432 dbname=movies sslmode=disable"

	// Откройте соединение с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка при открытии соединения с базой данных:", err)
	}
	defer db.Close()

	// Выполните запрос на получение пользователей из таблицы Users
	rows, err := db.Query("SELECT name, age, phonenumber FROM Users")
	if err != nil {
		log.Fatal("Ошибка при выполнении запроса:", err)
	}
	defer rows.Close()

	// Создайте слайс пользователей
	users := []User{}

	// Обработайте результаты запроса
	for rows.Next() {
		var name string
		var age int
		var phoneNumber string
		err := rows.Scan(&name, &age, &phoneNumber)
		if err != nil {
			log.Fatal("Ошибка при сканировании строки:", err)
		}
		user := User{Name: name, Age: age, PhoneNumber: phoneNumber}
		users = append(users, user)
	}

	// Верните результаты в формате JSON
	json.NewEncoder(w).Encode(users)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	// Получите имя пользователя из параметра запроса
	name := r.FormValue("name")

	// Получите новые значения для поля age и phone_number из тела запроса
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Установите параметры подключения к вашей базе данных PostgreSQL
	connStr := "user=your_username password=your_password host=your_host port=5432 dbname=your_database sslmode=disable"

	// Откройте соединение с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Обновите значения поля age и phone_number для пользователя с указанным именем
	updateUser := `
		UPDATE Users
		SET age = $1, phonenumber = $2
		WHERE name = $3
	`
	_, err = db.Exec(updateUser, user.Age, user.PhoneNumber, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Пользователь %s успешно обновлен", name)
}

