package main

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/valkey-io/valkey-go"
)

const (
	keyData   = "new-redis"
	valueData = "Valkey"
)

type valkeyContainer struct {
	testcontainers.Container
	Endpoint string
}

func main() {
	ctx := context.Background()

	container, err := createContainerWithValkey(ctx)
	if err != nil {
		log.Fatalf("Unable to create the container: %s", err)
	}
	defer func() {
		if err := container.Container.Terminate(ctx); err != nil {
			log.Fatalf("Unable to stop Valkey: %s", err)
		}
	}()

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{container.Endpoint},
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

// Best way to use Valkey with Docker containers. Everything is managed
// by the TestContainers library. You can find more information about
// TestContainers here: https://golang.testcontainers.org
func createContainerWithValkey(ctx context.Context) (*valkeyContainer, error) {
	containerRequest := testcontainers.ContainerRequest{
		Name:         "valkey",
		Image:        "valkey/valkey:7.2.5",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		log.Fatalf("Could not start Valkey: %s", err)
		return nil, err
	}
	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		log.Fatalf("Unable to retrieve the endpoint: %s", err)
	}
	return &valkeyContainer{
		Container: container,
		Endpoint:  endpoint,
	}, nil
}
