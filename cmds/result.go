package cmds

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func sendEmail() {
	// 1. Read and Parse the .wgsl file to get the stored email
	configData, fileErr := os.ReadFile(".wgsl")
	if fileErr != nil {
		log.Fatal("Cannot find the registered email. Please run 'init' first.")
	}

	// Split content (Line 1: Status, Line 2: Email)
	lines := strings.Split(strings.TrimSpace(string(configData)), "\n")
	if len(lines) < 2 {
		log.Fatal("Invalid configuration file format. Please re-run 'init'.")
	}
	recipient := strings.TrimSpace(lines[1])

	from := os.Getenv("FROM_EMAIL")
	password := os.Getenv("FROM_EMAIL_PASSWORD")
	smtpHost := os.Getenv("FROM_EMAIL_SMTP")
	smtpAddr := os.Getenv("SMTP_ADDR")

	// 2. Build the multipart message
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	// RFC 5322 Headers
	fmt.Fprintf(buf, "From: %s\r\n", from)
	fmt.Fprintf(buf, "To: %s\r\n", recipient)
	fmt.Fprintf(buf, "Subject: Leukemia Pattern Discovery Results\r\n")
	fmt.Fprintf(buf, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(buf, "Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary())
	fmt.Fprintf(buf, "\r\n")

	// 3. Add Text Body
	bodyPart, _ := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/plain; charset=utf-8"},
	})
	bodyPart.Write([]byte("The AML pattern discovery pipeline has completed. Visualizations are attached below."))

	// 4. Attach all .png files from the Images folder
	imageDir := "./Images"
	files, err := ioutil.ReadDir(imageDir)
	if err == nil {
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".png" {
				filePath := filepath.Join(imageDir, file.Name())
				fileData, _ := os.ReadFile(filePath)

				partHeader := make(map[string][]string)
				partHeader["Content-Type"] = []string{http.DetectContentType(fileData)}
				partHeader["Content-Transfer-Encoding"] = []string{"base64"}
				partHeader["Content-Disposition"] = []string{fmt.Sprintf("attachment; filename=\"%s\"", file.Name())}

				partWriter, _ := writer.CreatePart(partHeader)

				// Use a Base64 encoder that writes to the part
				encoder := base64.NewEncoder(base64.StdEncoding, partWriter)
				encoder.Write(fileData)
				encoder.Close()
			}
		}
	}

	// Close writer to finalize boundary
	writer.Close()

	// 5. Authenticate and Send
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Use recipient from file, not the empty global variable
	err = smtp.SendMail(smtpAddr, auth, from, []string{recipient}, buf.Bytes())

	if err != nil {
		log.Fatalf("Failed to send results: %v", err)
	}
}

func Result() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}
	sendEmail()
}
