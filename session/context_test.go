package session_test

import (
	"context"
	"testing"
	"time"

	"taylz.io/http/session"
)

func TestContext(t *testing.T) {

	ctx := context.Background()

	_, ok := session.FromContext(ctx)

	if ok {
		t.Log("expected context.Background().Value(ctxKey).(*session.T) is NOT ok")
		t.Fail()
	}

	s := session.New("1", "test", time.Now().Add(time.Millisecond))

	ctx = s.NewContext(ctx)

	ss, ok := session.FromContext(ctx)

	if !ok {
		t.Log("expected ctx.Value(ctxKey).(*session.T) is ok")
		t.Fail()
	} else if ss == nil {
		t.Log("expected ctx.Value(ctxKey).(*session.T) exists")
		t.Fail()
	} else if s != ss {
		t.Log("expected session.FromContext(session.NewContext()) equality")
		t.Fail()
	}

	ctx = session.NewContext(context.Background(), nil)

	ss, ok = session.FromContext(ctx)

	if !ok {
		t.Log("expected session.FromContext(session.NewContext(nil)) is ok")
		t.Fail()
	} else if ss != nil {
		t.Log("expected session.FromContext(session.NewContext(nil)) does NOT exist")
		t.Fail()
	}
}
