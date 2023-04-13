# awscrondoc

A tool to list up cron expressions registered in Amazon EventBridge for the company's internal wiki.

## Installation

### Library

```
go get github.com/shmokmt/awscrondoc@latest
```

### CLI

```
go install github.com/shmokmt/awscrondoc/cmd/awscrondoc@latest
```

## Usage

### Library

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

### CLI

```
awscrondoc
```
