package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"app.com/registry"
)

func Start(ctx context.Context, host, port string, service registry.Registration, registerHandlersFunc func()) (context.Context, error) {
	registerHandlersFunc()
	ctx = startService(ctx, service.ServiceName, host, port)
	err := registry.RegisterService(service)
	if err != nil {
		return ctx, err
	}

	return ctx, nil

}

func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = ":" + port

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Println("%v started. Press any key to stop. \n", serviceName)
		var s string
		fmt.Scanln(&s)
		err := registry.ShutdownService(fmt.Sprintf("http://%v:%v", host, port))
		if err != nil {
			log.Println(err)
		}
		srv.Shutdown(ctx)
		cancel()

	}()
	return ctx
}