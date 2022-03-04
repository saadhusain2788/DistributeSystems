package main

import (
	"context"
	"fmt"
	stlog "log"

	"app.com/registry"
	s "app.com/service"

	l "app.com/logger"
)

func main() {
	l.Run("./app.log")

	host, port := "localhost", "4000"

	var r registry.Registration
	r.ServiceName = registry.LogService
	r.ServiceURL = fmt.Sprintf("http://%v:%v", host, port)
	ctx, err := s.Start(context.Background(), host, port, r, l.RegisterHandler)
	if err != nil {
		stlog.Fatal(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down log service")
}
