# awscrondoc

[![Go Reference](https://pkg.go.dev/badge/github.com/shmokmt/awscrondoc.svg)](https://pkg.go.dev/github.com/shmokmt/awscrondoc)

A tool to list up cron expressions registered in following AWS services for the company's internal wiki.

- EventBridge Rule
- Glue Trigger

## Installation

### Library

```
go get github.com/shmokmt/awscrondoc@latest
```

### CLI

```
go install github.com/shmokmt/awscrondoc/cmd/awscrondoc@latest
```

## Permissions

It Requires the following minimum set of permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["events:ListRules"],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": ["glue:ListTriggers", "glue:GetTrigger"],
      "Resource": "*"
    }
  ]
}
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
