package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sensorapi/src/configuration"
	"sensorapi/src/domain"
	"sensorapi/src/persistence"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	conn := persistence.NewSqlxConnection(configuration.Load())
	userRepo := persistence.NewSqlxUserRepo(&conn)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Let's create the admin user")
	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	rawPass, _ := reader.ReadString('\n')
	password := strings.TrimSpace(rawPass)
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln("Error while hashing password")
	}
	userRepo.Save(domain.User{
		Name:     strings.TrimSpace(name),
		Password: string(hashed),
		Role:     domain.AdminUserRole,
	})
}
