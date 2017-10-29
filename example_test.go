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
		fmt.Printf("%s : %s\n", key, os.Getenv(key))
	}
	// Output:
	// S3_BUCKET : YOURS3BUCKET
	// SECRET_KEY : YOURSECRETKEYGOESHERE
	// MESSAGE : A message containing important spaces.
}

func ExampleReadFile() {
	env, err := dotenv.ReadFile("fixtures/example.env")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s : %s\n", "LIB", env["LIB"])
	// Output:
	// LIB : github.com/alexsasharegan/dotenv
}

func ExampleParseString() {
	envStrs := []string{
		`FOO=bar`,
		`FOO="escaped\"bar with quote"`,
		`FOO="bar\nbaz"`,
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
		fmt.Printf("%s : %s\n", k, v)
	}
	// Output:
	// FOO : bar
	// FOO : escaped"bar with quote
	// FOO : bar
	// baz
	// FOO.BAR : foobar
	// FOO : bar#baz
	// parsing error: invalid line
}
