package main

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func connectViaSsh(user, host string, password string) (*ssh.Client, *ssh.Session) {
	GlobalPassWord = password

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.KeyboardInteractive(SshInteractive), ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host, config)
	fmt.Println(err)
	session, err := client.NewSession()
	fmt.Println(err)

	return client, session
}

func SshInteractive(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
	answers = make([]string, len(questions))
	// The second parameter is unused
	for n := range questions {
		answers[n] = GlobalPassWord
	}
	return answers, nil
}
