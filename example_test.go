package dotenv_test

import (
	"fmt"
	"log"
	"os"

	"github.com/alexsasharegan/dotenv"
)

func ExampleLoad() {
	err := dotenv.Load("fixtures/example.env")
	if err != nil {
		log.Fatal(err)
	}

	envKeys := []string{"S3_BUCKET", "SECRET_KEY", "MESSAGE"}
	for _, key := range envKeys {
		fmt.Printf("%s : %s", key, os.Getenv(key))
	}
}

func ExampleReadFile() {
	env, err := dotenv.ReadFile("fixtures/example.env")
	if err != nil {
		log.Fatal(err)
	}

	for key, val := range env {
		fmt.Printf("%s : %s", key, val)
	}
}

func ExampleParseString() {
	envStrs := []string{
		`FOO=bar`,
		`FOO="escaped\"bar"`,
		`FOO="bar\nbaz`,
		`FOO.BAR=foobar`,
		`FOO="bar#baz" # comment`,
		`INVALID LINE`,
	}

	for _, s := range envStrs {
		k, v, err := dotenv.ParseString(s)
		if err != nil {
			fmt.Printf("parsing error: %v", err)
			continue
		}
		fmt.Printf("%s : %s", k, v)
	}
}
