package cmds

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/mdns"
	"github.com/joho/godotenv"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)
func readSecret(prompt string) (string, error) {
	fmt.Print(prompt)
	byteSecret, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return strings.TrimSpace(string(byteSecret)), nil
}
func downloadImage(IP, image string) {
	_ = godotenv.Load(".env")
	port := os.Getenv("SSH_PORT")
	if port == "" {
		port = "22"
	}
	userName, err := readSecret("Enter Username: ")
	if err != nil {
		fmt.Println("Error reading username")
		return
	}
	password, err := readSecret("Enter Password: ")
	if err != nil {
		fmt.Println("Error reading password")
		return
	}
	sshClientConfig := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         20 * time.Second,
	}
	target := fmt.Sprintf("%s:%s", IP, port)
	client, err := ssh.Dial("tcp", target, sshClientConfig)
	if err != nil {
		fmt.Printf("Authentication Failed: %v\n", err)
		return
	}
	defer client.Close()
	sftpc, err := sftp.NewClient(client)
	if err != nil {
		fmt.Println(" SFTP Session Failed")
		return
	}
	defer sftpc.Close()
	homeDir, _ := os.UserHomeDir()
	localPath := filepath.Join(homeDir, image)
	remoteFile, err := sftpc.Open(image)
	if err != nil {
		fmt.Printf("Remote Error: %s not found.\n", image)
		return
	}
	defer remoteFile.Close()
	localFile, err := os.Create(localPath)
	if err != nil {
		fmt.Printf("Local Error: Could not create %s\n", localPath)
		return
	}
	defer localFile.Close()
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		fmt.Println(" Transfer interrupted")
		return
	}
	fmt.Printf("Downloaded to: %s\nReady for analysis.\n", localPath)
}

func Get(imagePath string) {
	if !strings.HasSuffix(strings.ToLower(imagePath), ".jpeg") &&
	   !strings.HasSuffix(strings.ToLower(imagePath), ".png") {
		fmt.Println("Error: Unsupported format.")
		return
	}

	entrieschannel := make(chan *mdns.ServiceEntry, 4)

	var serviceName string
	fmt.Print("Enter the service name (e.g., Images): ")
	fmt.Scan(&serviceName)
	param := &mdns.QueryParam{
		Service: "_" + serviceName + "._tcp",
		Domain:  "local",
		Timeout: time.Second * 5,
		Entries: entrieschannel,
	}
	go func() {
		mdns.Query(param)
		close(entrieschannel)
	}()
	var foundIP string
	for val := range entrieschannel {
		if strings.Contains(val.Name, serviceName) {
			foundIP = val.AddrV4.String()
			break
		}
	}

	if foundIP == "" {
		fmt.Println(" Service not found on this network.")
		return
	}

	downloadImage(foundIP, imagePath)
}
