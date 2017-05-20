package main

import (
	"io"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/dgageot/jenkins-cli/jenkins"
)

type args struct {
	User   string `arg:"required,help:User,env:JENKINS_USER"`
	Token  string `arg:"required,help:Token,env:JENKINS_TOKEN"`
	Server string `arg:"required,help:Jenkins server,env:JENKINS_URL"`
	Query  string `arg:"required,help:api query,positional"`
}

func main() {
	var args args
	arg.MustParse(&args)

	body, err := jenkins.Get(args.User, args.Token, args.Server, args.Query)
	if err != nil {
		log.Fatal(err)
	}
	defer body.Close()

	if _, err := io.Copy(os.Stdout, body); err != nil {
		log.Fatal(err)
	}
}
