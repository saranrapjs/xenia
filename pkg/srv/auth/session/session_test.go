package session_test

import (
	"testing"
	"time"

	"github.com/coralproject/shelf/pkg/srv/auth/session"
	"github.com/coralproject/shelf/pkg/srv/mongo"
	"github.com/coralproject/shelf/pkg/tests"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	collection = "sessions"
	publicID   = "6dcda2da-92c3-11e5-8994-feff819cdc9f"
	context    = "testing"
)

func init() {
	tests.Init()
}

// removeSessions is used to clear out all the test sessions that are
// created from tests.
func removeSessions(ses *mgo.Session) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		return c.Remove(q)
	}

	err := mongo.ExecuteDB(context, ses, collection, f)
	if err != mgo.ErrNotFound {
		return err
	}

	return nil
}

func TestIsExpired(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to validate a session has expired.")
	{
		s := session.Session{
			DateExpires: time.Now().Add(-time.Hour),
		}

		t.Log("\tWhen using an expired session.")
		{
			if !s.IsExpired() {
				t.Fatalf("\t%s\tShould be expired.", tests.Failed)
			}
			t.Logf("\t%s\tShould be expired", tests.Success)
		}

		s = session.Session{
			DateExpires: time.Now().Add(time.Hour),
		}

		t.Log("\tWhen using an valid session")
		{
			if s.IsExpired() {
				t.Fatalf("\t%s\tShould Not be expired.", tests.Failed)
			}
			t.Logf("\t%s\tShould Not be expired", tests.Success)
		}
	}
}

// TestCreate tests the creation of sessions.
func TestCreate(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	ses := mongo.GetSession()
	defer ses.Close()

	defer func() {
		if err := removeSessions(ses); err != nil {
			t.Errorf("\t%s\tShould be able to remove all sessions : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove all sessions.", tests.Success)
	}()

	t.Log("Given the need to create sessions in the DB.")
	{
		t.Logf("\tWhen using PublicID %s", publicID)
		{
			if err := removeSessions(ses); err != nil {
				t.Fatalf("\t%s\tShould be able to remove all sessions : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to remove all sessions.", tests.Success)

			s1, err := session.Create(context, ses, publicID, 10*time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a session.", tests.Success)

			s2, err := session.Get(context, ses, s1.SessionID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the session.", tests.Success)

			if s1.SessionID != s2.SessionID {
				t.Fatalf("\t%s\tShould be able to get back the same session.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same session.", tests.Success)
			}

			if s1.PublicID != s2.PublicID {
				t.Fatalf("\t%s\tShould be able to get back the same user.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same user.", tests.Success)
			}
		}
	}
}

// TestGetNotFound tests when a session is not found.
func TestGetNotFound(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	ses := mongo.GetSession()
	defer ses.Close()

	t.Log("Given the need to test finding a session and it is not found.")
	{
		t.Logf("\tWhen using SessionID %s", "NOT EXISTS")
		{
			if _, err := session.Get(context, ses, "NOT EXISTS"); err == nil {
				t.Fatalf("\t%s\tShould Not be able to retrieve the session.", tests.Failed)
			}
			t.Logf("\t%s\tShould Not be able to retrieve the session.", tests.Success)
		}
	}
}

// TestNoSession tests when a nil session is used.
func TestNoSession(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to test calls with a bad session.")
	{
		t.Log("\tWhen using a nil session")
		{
			if _, err := session.Create(context, nil, publicID, 10*time.Second); err == nil {
				t.Errorf("\t%s\tShould be able to create a session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a session.", tests.Success)

			if _, err := session.Get(context, nil, "NOT EXISTS"); err == nil {
				t.Errorf("\t%s\tShould Not be able to retrieve the session.", tests.Failed)
			}
			t.Logf("\t%s\tShould Not be able to retrieve the session.", tests.Success)
		}
	}
}
