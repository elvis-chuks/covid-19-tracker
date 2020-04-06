package main 

import (
	"fmt"
	// "log"
	"errors"
	"net/http"
	"database/sql"
	"encoding/json"

	_ "github.com/lib/pq"
	// "github.com/gorilla/mux"
	"github.com/ichtrojan/thoth"
	"golang.org/x/crypto/bcrypt"

)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "@123elvischuks"
	dbname = "covidtracker"
)
var db *sql.DB

type Resp map[string]interface{}

type User struct{
	FirstName string `json:"firstname,omitempty"`
	LastName string `json:"lastname,omitempty"`
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

var (
	file, _ = thoth.Init("log")
)

func InitDB() *sql.DB{
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+" password=%s dbname=%s sslmode=disable",
	host,port,user,password,dbname)

	db, err := sql.Open("postgres",psqlInfo)

	if err != nil{
		panic(err)
	}

	// defer db.Close()

	return db

}


func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func Register(w http.ResponseWriter, r *http.Request){

	w.Header().Set("content-type","application/json")

	var user User

	_ = json.NewDecoder(r.Body).Decode(&user)

	fmt.Println(user)

	db := InitDB()
	defer db.Close()

	pash, err := HashPassword(user.Password)
	
	if err != nil{
		fmt.Println(err)
	}

	query := fmt.Sprintf(
	"INSERT INTO users(firstname,lastname,email,password) VALUES('%s','%s','%s','%s');",user.FirstName,user.LastName,user.Email,pash)

	_, err = db.Exec(query)

	if err != nil{

		file.Log(errors.New(err.Error()))

		res := Resp{"status":"error","msg":err.Error()}

		json.NewEncoder(w).Encode(res)
	}else{

		res := Resp{"status":"success"}

		json.NewEncoder(w).Encode(res)
	}

}



func main(){

	// router := mux.NewRouter()

	http.HandleFunc("/register",Register)

	// file, err := thoth.Init("log")


	err := file.Serve("/logs","12345")

	if err != nil {
		file.Log(err)
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("application running on port 8080")
		file.Log(err)
	}
}