package main

import (
	"log"

	"github.com/shmokmt/awscrondoc"
)

func main() {
	d, err := awscrondoc.New()
	if err != nil {
		log.Fatal(err)
	}
	_, err = d.MarkdownString()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(md)
}
