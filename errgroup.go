// Package errgroup provides a means to run goroutines as a logical
// unit where an error from one goroutine will trigger the termination
// of all other goroutines in the group. Non-error (nil) returns do not
// cause the group to terminate.
package errgroup

// member is a member of a group. It defines the function which will
// be run within a goroutine and a function which will be called on
// group termination.
type member struct {
	routine   func() error
	terminate func(error)
}

// Group holds a collection of members which whose routines are run
// concurrently. Any non-nil error from a member routine will cause the
// Group to terminate.
type Group struct {
	members []*member

	onError func(err error)
}

// Add a new member to the Group.
//
// All members should define a routine to be called and a termination
// function. The termination function should cause the member's routine
// to return. Additionally, it should be safe to call the terminate function
// after the routine has returned.
func (g *Group) Add(routine func() error, terminate func(error)) {
	g.members = append(g.members, &member{routine, terminate})
}

// OnError registers an error handler with the Group.
//
// The error handler is optional and is run prior to terminating members of the
// group. It can be used for things like logging out the trapped error.
func (g *Group) OnError(handler func(err error)) {
	g.onError = handler
}

// Run the routines of all Group members concurrently.
//
// If a routine terminates with a nil error, the other members will continue
// to run. When the first non-nil error is returned from a member routine, all
// members of the Group will be terminated. This function does not return until
// all members have terminated. Once all members terminate, this will return
// the error which triggered the group termination.
//
// Note that if a member routine returns a nil error, its terminate function
// will not be called until a non-nil error is returned by another member of
// the group.
func (g *Group) Run() error {
	// If there are no members of the group, there is nothing to do.
	if len(g.members) == 0 {
		return nil
	}

	// Run the goroutine for each member of the group.
	errors := make(chan error, len(g.members))
	for _, m := range g.members {
		go func(m *member) {
			errors <- m.routine()
		}(m)
	}

	// Wait for the first non-nil error returned.
	var terminated int
	var err error
	for e := range errors {
		terminated++
		if e != nil {
			err = e
			break
		}
		if terminated == cap(errors) {
			break
		}
	}

	// If an error handler is specified and there is an error,
	// execute the handler function.
	if err != nil && g.onError != nil {
		g.onError(err)
	}

	// Terminate all group members.
	for _, member := range g.members {
		member.terminate(err)
	}

	// Wait for all the members to terminate.
	for i := terminated; i < cap(errors); i++ {
		<-errors
	}

	return err
}
