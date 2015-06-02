package main

import "golang.org/x/crypto/ssh"
import "fmt"
import "bytes"

func main() {
	config := &ssh.ClientConfig{
		User: "mirko",
		Auth: []ssh.AuthMethod{
			ssh.Password("yourpassword"),
		},
	}
	client, err := ssh.Dial("tcp", "bithalde.de:22", config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		panic("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())

}
