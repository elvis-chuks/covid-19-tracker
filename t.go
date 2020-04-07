package main

import (
	"os"
	"fmt"
)
const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "@123elvischuks"
	dbname = "covidtracker"
)
func main(){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+" password=%s dbname=%s sslmode=disable",
	host,port,user,password,dbname)

	os.Setenv("DATABASE_URL",psqlInfo)
	fmt.Println("done")
}