package main

import (
	"fmt"
	"log"

	"github.com/shmokmt/awscrondoc"
)

func main() {
	d, err := awscrondoc.New()
	if err != nil {
		log.Fatal(err)
	}
	rules, err := d.ListRules()
	if err != nil {
		log.Fatal(err)
	}
	md, err := d.MarkdownString(rules)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(md)
}
