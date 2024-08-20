package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	valkeydocker "github.com/testcontainers/testcontainers-go/modules/valkey"
	"github.com/valkey-io/valkey-go"
)

const (
	keyData   = "new-redis"
	valueData = "Valkey"
)

func main() {
	ctx := context.Background()

	// Start Valkey with Testcontainers
	container, err := valkeydocker.Run(ctx,
		"valkey/valkey:7.2.6",
		valkeydocker.WithLogLevel(valkeydocker.LogLevelVerbose),
		valkeydocker.WithConfigFile(filepath.Join("conf", "valkey.conf")),
	)
	if err != nil {
		log.Fatalf("Unable to create the container: %s", err)
	}
	defer func() {
		if err := container.Container.Terminate(ctx); err != nil {
			log.Fatalf("Unable to stop Valkey: %s", err)
		}
	}()
	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		log.Fatalf("Unable to get the container endpoint: %s", err)
	}

	// Create a new Valkey client
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{endpoint},
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Write data using the SET command
	err = client.Do(ctx, client.B().Set().Key(keyData).Value(valueData).Build()).Error()
	if err != nil {
		panic(err)
	}

	// Read data using the GET command
	value, err := client.Do(ctx, client.B().Get().Key(keyData).Build()).ToString()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}
