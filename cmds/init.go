package cmds

import (
	"fmt"
	"strings"
	"github.com/joho/godotenv"
	"os"
	"net/smtp"
	"strconv"
	"math/rand"
	"time"
)
var doctor_email string
func is_validEmail(email string) bool {
	if email == "" {
		return false
	}
	supported_email := []string{"@gmail.com","@outlook.com","@icloud.com","@yahoo.com"} // supported email services
	flag := false
	for _,val := range supported_email {
		if strings.Contains(email,val) && strings.HasSuffix(email,val) {
			flag = true
			break
		}
	}
	return flag
}
func send_AuthenticationCode(email string,body int) error {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("FROM_EMAIL"),
		os.Getenv("FROM_EMAIL_PASSWORD"),
		os.Getenv("FROM_EMAIL_SMTP"),
	)
	message := "Subject:Verification Code\n"+strconv.Itoa(body)
	return smtp.SendMail(
		os.Getenv("SMTP_ADDR"),
		auth,
		os.Getenv("FROM_EMAIL"),
		[]string{email},
		[]byte(message),
	)
}
func Init(){
	fmt.Print("Enter Doctor Email: ")
	fmt.Scanf("%s",&doctor_email)
	if !is_validEmail(doctor_email) {
		panic("Enter a valid email")
	}
	err := godotenv.Load(".env")
	if err != nil {
		panic("An internal error has occurred please try again")
	}
	rand.Seed(time.Now().UnixNano())
	verification_code := 100000 +rand.Intn(900000)
	err = send_AuthenticationCode(doctor_email,verification_code)
	if err != nil {
		panic("Error while sending an verification code to the given email ")
	}
	fmt.Print("Enter verification code: ")
	var code int
	fmt.Scanf("%d",&code)
	if code == verification_code {
		fmt.Println("Initialised Successfully")
		os.WriteFile(".wgsl",[]byte("1"),0644)
	} else{
		fmt.Println("Verification was unsuccessfull, unable to initialse init in the current folder please try again")
	}
}
