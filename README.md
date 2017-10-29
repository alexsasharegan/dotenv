# dotenv

![build status](https://travis-ci.org/alexsasharegan/dotenv.svg?branch=master)

A Go (golang) implementation of dotenv _(inspired by: [https://github.com/joho/godotenv](https://github.com/joho/godotenv))_.

## Installation

As a **Library**:

```sh
go get github.com/alexsasharegan/dotenv
```

## Usage

In your environment file (canonically named `.env`):

```environ
S3_BUCKET=YOURS3BUCKET
SECRET_KEY=YOURSECRETKEYGOESHERE

MESSAGE="A message containing important spaces."
ESCAPED='You can escape you\'re strings too.'

# A comment line that will be ignored
GIT_PROVIDER=github.com
LIB=${GIT_PROVIDER}/alexsasharegan/dotenv # variable interpolation (plus ignored trailing comment)
```

```go
package main

import (
    "github.com/alexsasharegan/dotenv"
    "fmt"
    "log"
    "os"
)

func main() {
  err := dotenv.Load()
  if err != nil {
    log.Fatalf("Error loading .env file: %v", err)
  }

  s3Bucket := os.Getenv("S3_BUCKET")
  secretKey := os.Getenv("SECRET_KEY")

  fmt.Println(os.Getenv("MESSAGE"))
}
```

## Documentation

[https://godoc.org/github.com/alexsasharegan/dotenv](https://godoc.org/github.com/alexsasharegan/dotenv)
