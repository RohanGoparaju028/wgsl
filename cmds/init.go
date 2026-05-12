package cmds

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/joho/godotenv"
	"golang.org/x/term"
)
var DoctorEmail string

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	supportedDomains := []string{
		"@gmail.com",
		"@outlook.com",
		"@icloud.com",
		"@yahoo.com",
	}
	for _, domain := range supportedDomains {
		if strings.HasSuffix(email, domain) {
			return true
		}
	}
	return false
}

func sendAuthenticationCode(email string, code int) error {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("FROM_EMAIL"),
		os.Getenv("FROM_EMAIL_PASSWORD"),
		os.Getenv("FROM_EMAIL_SMTP"),
	)
	message := "Subject: Verification Code\r\n\r\n" + strconv.Itoa(code)
	return smtp.SendMail(
		os.Getenv("SMTP_ADDR"),
		auth,
		os.Getenv("FROM_EMAIL"),
		[]string{email},
		[]byte(message),
	)
}

func hideEmail() (string, error) {
	rawInput, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return strings.TrimSpace(string(rawInput)), nil
}

func hideOTP() (int, error) {
	rawInput, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return 0, err
	}
	fmt.Println()
	otp, err := strconv.Atoi(strings.TrimSpace(string(rawInput)))
	return otp, err
}
func generateSecureOTP() (int, error) {
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return int(n.Int64()) + 100000, nil
}

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error: Missing .env file. Please configure the environment.")
		os.Exit(1)
	}

	fmt.Print("Enter Doctor Email: ")
	DoctorEmail, err = hideEmail()
	if err != nil {
		fmt.Println("Error reading email:", err)
		return
	}

	if len(DoctorEmail) > 254 {
		fmt.Println("Error: Email address is too long.")
		return
	}

	if !isValidEmail(DoctorEmail) {
		fmt.Println("Error: Please enter a valid email (gmail, outlook, icloud, or yahoo).")
		return
	}

	verificationCode, err := generateSecureOTP()
	if err != nil {
		fmt.Println("Error generating verification code. Please try again.")
		return
	}
	err = sendAuthenticationCode(DoctorEmail, verificationCode)
	if err != nil {
		fmt.Println("Error: Could not send verification code to the provided email.")
		return
	}
	sentAt := time.Now()

	fmt.Print("Enter verification code: ")
	code, err := hideOTP()
	if err != nil {
		fmt.Println("Error reading verification code:", err)
		return
	}

	if time.Since(sentAt) > 5*time.Minute {
		fmt.Println("Error: Verification code has expired. Please run init again.")
		return
	}

	if code == verificationCode {
		fmt.Println("Initialised successfully.")
		drEmail := fmt.Sprintf("1\n%s",DoctorEmail)
		os.WriteFile(".wgsl", []byte(drEmail), 0600)
	} else {
		fmt.Println("Incorrect verification code. Please run init again.")
	}
}
