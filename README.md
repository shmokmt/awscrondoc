# awscrondoc

A tool to list up cron expressions registered in Amazon EventBridge for the company's internal wiki.

## Usage

```go
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
	md, err := d.MarkdownString()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(md)
}
```

or

```
awscrondoc
```
