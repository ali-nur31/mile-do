package main

import (
	"context"

	"github.com/ali-nur31/taskee/pkg/postgres"
)

func main() {
	ctx := context.Background()

	postgres.InitializeDatabaseConnection(ctx)
}
