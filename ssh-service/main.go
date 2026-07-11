package main

import (
	"fmt"
	"log"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		s.Write(gossh.MarshalAuthorizedKey(s.PublicKey()))

		for term := term.NewTerminal(s, fmt.Sprintf("%s: ", s.User())); ; {
			line, err := term.ReadLine()
			if err != nil {
				log.Println(err)
				break
			}
			log.Println(line)
		}
	})
	publicKeyOption := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		return true
	})
	ssh.ListenAndServe(":2222", nil, publicKeyOption)
}
