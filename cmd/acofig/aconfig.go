package main

import (
	"fmt"
	"github.com/cristalhq/aconfig"
	"log"
	"os"
)

type MyConfig struct {
	HTTPPort int `default:"1111" usage:"just a number"`
	Auth     struct {
		User string `default:"def-user" usage:"your user"`
		Pass string `default:"def-pass" usage:"make it strong"`
	}
}

func main() {

	Example_Env()

}

func Example_Env() {
	os.Setenv("EXAMPLE_HTTP_PORT", "3333")
	os.Setenv("EXAMPLE_AUTH_USER", "env-user")
	os.Setenv("EXAMPLE_AUTH_PASS", "env-pass")
	defer os.Clearenv()

	var cfg MyConfig
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		EnvPrefix: "EXAMPLE",
	})
	if err := loader.Load(); err != nil {
		log.Panic(err)
	}

	fmt.Printf("HTTPPort:  %v\n", cfg.HTTPPort)
	fmt.Printf("Auth.User: %v\n", cfg.Auth.User)
	fmt.Printf("Auth.Pass: %v\n", cfg.Auth.Pass)

	// Output:
	//
	// HTTPPort:  3333
	// Auth.User: env-user
	// Auth.Pass: env-pass
}
