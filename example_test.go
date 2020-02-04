package errgroup

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ExampleGroup_Run_simple() {
	var g Group

	// Add a member which immediately returns nil.
	g.Add(
		func() error {
			fmt.Println("member 1: returning nil immediately")
			return nil
		},
		func(e error) {
			fmt.Println("member 1: terminating")
		},
	)

	// Add a member which can be cancelled or times out.
	cancel := make(chan struct{})
	g.Add(
		func() error {
			select {
			case <-time.After(1 * time.Second):
				fmt.Println("member 2: timed out after 1s")
				return nil
			case <-cancel:
				fmt.Println("member 2: canceled")
				return nil
			}

		},
		func(e error) {
			fmt.Println("member 2: terminating")
			close(cancel)
		},
	)

	// Add a member which errors after a short wait.
	g.Add(
		func() error {
			time.Sleep(500 * time.Millisecond)
			fmt.Println("member 3: erroring after 500ms")
			return errors.New("tearing down")
		},
		func(e error) {
			fmt.Println("member 3: terminating")
		},
	)

	err := g.Run()
	fmt.Printf("error: %s\n", err)
	// Output:
	// member 1: returning nil immediately
	// member 3: erroring after 500ms
	// member 1: terminating
	// member 2: terminating
	// member 3: terminating
	// member 2: canceled
	// error: tearing down
}

func ExampleGroup_Run_context() {
	var g Group

	// Add a member which will run until a context is canceled.
	ctx, cancel := context.WithCancel(context.Background())
	g.Add(
		func() error {
			<-ctx.Done()
			return ctx.Err()
		},
		func(e error) {
			cancel()
		},
	)

	// Cancel the context
	go cancel()
	err := g.Run()

	fmt.Printf("error: %s\n", err)
	// Output:
	// error: context canceled
}

func ExampleGroup_Run_signals() {
	var g Group

	// Add a member which terminates on signal.
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGHUP)
	g.Add(
		func() error {
			<-sig
			return errors.New("terminated on SIGHUP")
		},
		func(e error) {
			close(sig)
		},
	)

	// Add a member which runs a web server.
	listener, _ := net.Listen("tcp", ":0")
	g.Add(
		func() error {
			defer fmt.Println("server stopped")
			return http.Serve(listener, http.NewServeMux())
		},
		func(e error) {
			_ = listener.Close()
		},
	)

	// Simulate keyboard ^C
	time.Sleep(250 * time.Millisecond)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	err := g.Run()

	fmt.Printf("error: %s\n", err)
	// Output:
	// server stopped
	// error: terminated on SIGHUP
}

func ExampleGroup_Run_withErrorHandler() {
	var g Group

	// Add a member which immediately returns nil.
	g.Add(
		func() error {
			fmt.Println("member 1: returning nil immediately")
			return nil
		},
		func(e error) {
			fmt.Println("member 1: terminating")
		},
	)

	// Add a member which errors after a short wait.
	g.Add(
		func() error {
			time.Sleep(500 * time.Millisecond)
			fmt.Println("member 2: erroring after 500ms")
			return errors.New("tearing down")
		},
		func(e error) {
			fmt.Println("member 2: terminating")
		},
	)

	// Register an error handler which is called prior to
	// terminating all members.
	g.OnError(func(err error) {
		fmt.Printf("on error: %v\n", err)
	})

	err := g.Run()
	fmt.Printf("error: %s\n", err)
	// Output:
	// member 1: returning nil immediately
	// member 2: erroring after 500ms
	// on error: tearing down
	// member 1: terminating
	// member 2: terminating
	// error: tearing down
}
