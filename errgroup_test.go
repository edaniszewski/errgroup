package errgroup

import (
	"errors"
	"testing"
	"time"
)

var errTest = errors.New("test error")

func TestGroup_Add(t *testing.T) {
	var g Group
	g.Add(
		func() error {
			return nil
		},
		func(e error) {

		},
	)

	if len(g.members) != 1 {
		t.Errorf("no members added")
	}
}

func TestGroup_RunEmpty(t *testing.T) {
	var g Group

	res := make(chan error)
	go func() {
		res <- g.Run()
	}()

	select {
	case err := <-res:
		if err != nil {
			t.Error(err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("test case timeout")
	}
}

func TestGroup_RunOneNil(t *testing.T) {
	var calledRoutine bool
	var calledTerminate bool

	var g Group
	g.Add(
		func() error {
			calledRoutine = true
			return nil
		},
		func(e error) {
			calledTerminate = true
		},
	)

	res := make(chan error)
	defer close(res)

	go func() {
		res <- g.Run()
	}()

	select {
	case err := <-res:
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
		}
		if !calledRoutine {
			t.Error("routine not called")
		}
		if !calledTerminate {
			t.Error("terminate not called")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("test case timeout")
	}
}

func TestGroup_RunOneError(t *testing.T) {
	var calledRoutine bool
	var calledTerminate bool

	var g Group
	g.Add(
		func() error {
			calledRoutine = true
			return errTest
		},
		func(e error) {
			calledTerminate = true
		},
	)

	res := make(chan error)
	defer close(res)

	go func() {
		res <- g.Run()
	}()

	select {
	case err := <-res:
		if err != errTest {
			t.Errorf("got unexpected error: %v", err)
		}
		if !calledRoutine {
			t.Error("routine not called")
		}
		if !calledTerminate {
			t.Error("terminate not called")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("test case timeout")
	}
}

func TestGroup_RunMultipleNil(t *testing.T) {
	var (
		calledRoutine1   bool
		calledRoutine2   bool
		calledRoutine3   bool
		calledTerminate1 bool
		calledTerminate2 bool
		calledTerminate3 bool
	)

	var g Group
	// First member
	g.Add(
		func() error {
			calledRoutine1 = true
			return nil
		},
		func(e error) {
			calledTerminate1 = true
		},
	)

	// Second member
	g.Add(
		func() error {
			calledRoutine2 = true
			return nil
		},
		func(e error) {
			calledTerminate2 = true
		},
	)

	// Third member
	g.Add(
		func() error {
			calledRoutine3 = true
			return nil
		},
		func(e error) {
			calledTerminate3 = true
		},
	)

	res := make(chan error)
	defer close(res)

	go func() {
		res <- g.Run()
	}()

	select {
	case err := <-res:
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
		}
		if !calledRoutine1 {
			t.Error("routine1 not called")
		}
		if !calledTerminate1 {
			t.Error("terminate1 not called")
		}
		if !calledRoutine2 {
			t.Error("routine2 not called")
		}
		if !calledTerminate2 {
			t.Error("terminate2 not called")
		}
		if !calledRoutine3 {
			t.Error("routine3 not called")
		}
		if !calledTerminate3 {
			t.Error("terminate3 not called")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("test case timeout")
	}
}

func TestGroup_RunMultipleError(t *testing.T) {
	var (
		calledRoutine1   bool
		calledRoutine2   bool
		calledRoutine3   bool
		calledTerminate1 bool
		calledTerminate2 bool
		calledTerminate3 bool
	)

	cancel1 := make(chan struct{})
	cancel2 := make(chan struct{})

	var g Group
	// First member
	g.Add(
		func() error {
			<-cancel1
			calledRoutine1 = true
			return nil
		},
		func(e error) {
			close(cancel1)
			calledTerminate1 = true
		},
	)

	// Second member
	g.Add(
		func() error {
			<-cancel2
			calledRoutine2 = true
			return nil
		},
		func(e error) {
			close(cancel2)
			calledTerminate2 = true
		},
	)

	// Third member
	g.Add(
		func() error {
			time.Sleep(50 * time.Millisecond)
			calledRoutine3 = true
			return errTest
		},
		func(e error) {
			calledTerminate3 = true
		},
	)

	res := make(chan error)
	defer close(res)

	go func() {
		res <- g.Run()
	}()

	select {
	case err := <-res:
		if err != errTest {
			t.Errorf("got unexpected error: %v", err)
		}
		if !calledRoutine1 {
			t.Error("routine1 not called")
		}
		if !calledTerminate1 {
			t.Error("terminate1 not called")
		}
		if !calledRoutine2 {
			t.Error("routine2 not called")
		}
		if !calledTerminate2 {
			t.Error("terminate2 not called")
		}
		if !calledRoutine3 {
			t.Error("routine3 not called")
		}
		if !calledTerminate3 {
			t.Error("terminate3 not called")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("test case timeout")
	}
}
