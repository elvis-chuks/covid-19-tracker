package main 

import (
	"fmt"
	"os"
	// "errors"
	"net/http"
	"database/sql"
	"encoding/json"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

)



// const (
// 	host = "localhost"
// 	port = 5432
// 	user = "postgres"
// 	password = "@123elvischuks"
// 	dbname = "covidtracker"
// )
var db *sql.DB

type Resp map[string]interface{}

type User struct{
	FirstName string `json:"firstname,omitempty"`
	LastName string `json:"lastname,omitempty"`
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}


func InitDB() *sql.DB{
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+" password=%s dbname=%s sslmode=disable",
	// host,port,user,password,dbname)


	// db, err := sql.Open("postgres",psqlInfo)

	db, err := sql.Open("postgres",os.Getenv("DATABASE_URL"))

	// db, err := sql.Open("postgres","postgres://urfaimyxtxvoov:d0dd7b17a756c0ff932f6db616711b04600fe4d3dde32449b279cfe7b8e83e75@ec2-18-210-51-239.compute-1.amazonaws.com:5432/d88051tuei2207")

	if err != nil{
		panic(err)
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS users(id SERIAL,email VARCHAR, firstname VARCHAR,lastname VARCHAR,password VARCHAR)")
	// defer db.Close()

	_, err = db.Exec(query)

	if err != nil {
		panic(err)
	}
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

		res := Resp{"status":"error","msg":err.Error()}

		json.NewEncoder(w).Encode(res)
	}else{

		res := Resp{"status":"success"}

		json.NewEncoder(w).Encode(res)
	}

}

// login endpoint

func Login(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type","application/json")
	if r.Method == "GET"{
		var user User

		_ = json.NewDecoder(r.Body).Decode(&user)

		db := InitDB()
		defer db.Close()

		query := fmt.Sprintf("SELECT email,firstname,lastname,password from users where email = '%s'",user.Email)

		rows, err := db.Query(query)

		if err != nil{
			fmt.Println(err)

		}else{
			defer rows.Close()
			var email,firstname,lastname,password string
			for rows.Next(){
				rows.Scan(&email,&firstname,&lastname,&password)
			}
			if CheckPasswordHash(user.Password,password)== true {
				fmt.Println("correct password")
				res := Resp{"status":"success","user_details":Resp{"firstname":firstname,"lastname":lastname}}
				json.NewEncoder(w).Encode(res)
			}else{
				fmt.Println("Incorrect password")
				res := Resp{"status":"failed","msg":"Incorrect user credentials"}
				json.NewEncoder(w).Encode(res)
			}
		}


	}else{
		http.Error(w, "Method Not Allowed",400)
	}

}

func main(){

	http.HandleFunc("/v1/register",Register)
	http.HandleFunc("/v1/login",Login)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("this is port",port)
	port = fmt.Sprintf(":"+"%s",port)
	fmt.Println(port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println(err)
	}
}