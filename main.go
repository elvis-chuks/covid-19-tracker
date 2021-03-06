package main 

import (
	"fmt"
	"os"
	"time"
	"net/http"
	"database/sql"
	"encoding/json"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"

)

var jwtKey = []byte("purevibes")

var db *sql.DB

type Resp map[string]interface{}

type User struct{
	FirstName string `json:"firstname,omitempty"`
	LastName string `json:"lastname,omitempty"`
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type Symptom struct {
	Email string `json:"email,omitempty"`
	Day string `json:"day,omitempty"`
	Month string `json:"month,omitempty"`
	Year string `json:"year,omitempty"`
	Score string `json:"score,omitempty"`
	Prognosis string `json:"prognosis,omitempty"`
	Date string `json:"theDate,omitempty"`
	Token string `json:"token,omitempty"`
}

type Question struct {
	Question string `json:"question,omitempty"`
	Point string `json:"point,omitempty"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}


func InitDB() *sql.DB{
	db, err := sql.Open("postgres",os.Getenv("DATABASE_URL"))

	if err != nil{
		panic(err)
	}
	
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS users(id SERIAL,email VARCHAR PRIMARY KEY, firstname VARCHAR,lastname VARCHAR,password VARCHAR)")
	
	_, err = db.Exec(query)

	query1 := fmt.Sprintf("CREATE TABLE IF NOT EXISTS symptoms(id SERIAL,email VARCHAR, day VARCHAR, month VARCHAR, year VARCHAR,date DATE, score VARCHAR,prognosis varchar)")

	_, _ = db.Exec(query1)

	query2 := fmt.Sprintf("CREATE TABLE IF NOT EXISTS questions(id SERIAL,question VARCHAR,point VARCHAR)")

	_, _ = db.Exec(query2)

	if err != nil {
		panic(err)
	}
	return db

}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
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
	setupResponse(&w,r)
	if (*r).Method == "OPTIONS" {
		return
	}


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


func Login(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type","application/json")
	setupResponse(&w,r)
	if (*r).Method == "OPTIONS" {
		return
	}
	if r.Method == "POST"{
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

				expirationTime := time.Now().Add(120 * time.Minute)

				claims := &Claims{
					Email: email,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: expirationTime.Unix(),
					},
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)


				tokenString,err := token.SignedString(jwtKey)

				if err != nil{
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name: "token",
					Value: tokenString,
					Expires: expirationTime,
				})

				res := Resp{"status":"success","token":tokenString,"user_details":Resp{"firstname":firstname,"lastname":lastname}}
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

func Symptoms(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type","application/json")
	setupResponse(&w,r)
	if (*r).Method == "OPTIONS" {
		return
	}
	var symptom Symptom
	
	_ = json.NewDecoder(r.Body).Decode(&symptom)
	
	if symptom.Token == ""{
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Resp{"status":"failed","msg":"Bad Request"})
		return
	}

	tknStr := symptom.Token

	fmt.Println("tknstr",tknStr)

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error){
		return jwtKey, nil
	})
	if err != nil{
		if err ==jwt.ErrSignatureInvalid {
			fmt.Println("unauthorized 2")
			w.WriteHeader(http.StatusUnauthorized)

			json.NewEncoder(w).Encode(Resp{"status":"failed","msg":"Unauthorized"})
			return
		}
		fmt.Println("bad req 2")
		
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Resp{"status":"failed","msg":"Token expired"})
		return
	}

	if !tkn.Valid {
		fmt.Println("unauthorized 3")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Resp{"status":"failed","msg":"Invalid token"})
		return
	}
	fmt.Println(claims.Email)
	if r.Method == "POST"{
		
		query := fmt.Sprintf("INSERT INTO symptoms(email,day,month,year,date,score,prognosis) VALUES('%s','%s','%s','%s','%s','%s','%s')",
		claims.Email,symptom.Day,symptom.Month,symptom.Year,symptom.Date,symptom.Score,symptom.Prognosis)
		
		
		db := InitDB()
		defer db.Close()

		_, err := db.Exec(query)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error",500)
		}else{
			res := Resp{"status":"success"}
			json.NewEncoder(w).Encode(res)
		}
	}else if r.Method == "GET"{

		db := InitDB()
		defer db.Close()

		query := fmt.Sprintf("SELECT day,month,year,date,score,prognosis FROM symptoms WHERE email = '%s'",claims.Email)

		rows,err := db.Query(query)
		
		if err != nil {
			fmt.Println(err)
			http.Error(w, `{'status':'error','msg':'Internal Server Error in query'}`,500)
		}else{
			defer rows.Close()

			var respList []Resp

			for rows.Next(){
				var day,month,year,date,score,prognosis string

				rows.Scan(&day,&month,&year,&date,&score,&prognosis)

				resMap := Resp{"day":day,"month":month,"year":year,"score":score,"prognosis":prognosis}

				respList = append(respList,resMap)
			}
			
			res := Resp{"status":"success","symptoms":respList}
			json.NewEncoder(w).Encode(res)
		}
	}else{
		http.Error(w, `{'status':'error','msg':'Method Not Allowed'}`,400)
	}
}


func Questions(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type","application/json")
	setupResponse(&w,r)
	if (*r).Method == "OPTIONS" {
		return
	}
	
	var question Question

	_ = json.NewDecoder(r.Body).Decode(&question)

	fmt.Println(question)
	if r.Method == "POST"{

		db := InitDB()

		defer db.Close()

		query := fmt.Sprintf("INSERT INTO questions(question,point) VALUES('%s','%s')",question.Question,question.Point)

		_, err := db.Exec(query)

		if err != nil {
			fmt.Println(err)
			res := Resp{"status":"error"}
			json.NewEncoder(w).Encode(res)
		}else{
			res := Resp{"status":"success"}
			json.NewEncoder(w).Encode(res)
		}
		
	}else if r.Method == "GET"{
		db := InitDB()

		defer db.Close()

		query := fmt.Sprintf("SELECT question,point from questions")

		rows, err := db.Query(query)

		if err != nil {
			fmt.Println(err)
			res := Resp{"status":"error","msg":"server error"}
			json.NewEncoder(w).Encode(res)
		}else{
			defer rows.Close()
			var respList []Resp
			for rows.Next(){
				var question,point string

				rows.Scan(&question,&point)

				respList = append(respList,Resp{"question":question,"point":point})

			}

			res := Resp{"status":"success","questions":respList}

			json.NewEncoder(w).Encode(res)
		}
	}else{
		http.Error(w,"Method not Allowed",400)
	}
}


func Refresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type","application/json")
	setupResponse(&w,r)
	if (*r).Method == "OPTIONS" {
		return
	}
	
	var symptom Symptom
	
	_ = json.NewDecoder(r.Body).Decode(&symptom)
	
	if symptom.Token == ""{
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Resp{"status":"failed","msg":"Bad Request"})
		return
	}

	tknStr := symptom.Token

	fmt.Println("tknstr",tknStr)
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if !tkn.Valid {
		fmt.Println("invalid token")
		w.WriteHeader(http.StatusUnauthorized)
		res  := Resp{"status":"failed","msg":"invalid token"}
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			res  := Resp{"status":"failed","msg":"invalid token signature"}
			json.NewEncoder(w).Encode(res)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		res  := Resp{"status":"failed","msg":"error parsing token"}
		json.NewEncoder(w).Encode(res)
		return
	}
	

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	// if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Minute {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	res  := Resp{"status":"failed","msg":"too soon"}
	// 	json.NewEncoder(w).Encode(res)
	// 	return
	// }

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	res  := Resp{"status":"success","token":tokenString}
	json.NewEncoder(w).Encode(res)
}


func main(){

	http.HandleFunc("/v1/register",Register)
	http.HandleFunc("/v1/login",Login)
	http.HandleFunc("/v1/symptoms",Symptoms)
	http.HandleFunc("/v1/questions",Questions)
	http.HandleFunc("/v1/refresh",Refresh)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	port = fmt.Sprintf(":"+"%s",port)
	

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println(err)
	}
}