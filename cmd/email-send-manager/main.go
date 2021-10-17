package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"email-send-manager/internal/app"
	"email-send-manager/pkg/logger"
)

// VERSION 版本号，可以通过编译的方式指定版本号：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "1.0.0"

func main() {
	rand.Seed(time.Now().Unix())
	ctx := logger.NewTagContext(context.Background(), "__main__")

	opts := []app.Option{app.SetVersion(VERSION)}

	if len(os.Args) > 2 {
		fmt.Printf("Usage: %s [path/to/config.yaml]\n", os.Args[0])
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		opts = append(opts, app.SetConfigFile(os.Args[1]))
	}

	err := app.Run(ctx, opts...)
	if err != nil {
		logger.WithContext(ctx).Errorf(err.Error())
	}
}
