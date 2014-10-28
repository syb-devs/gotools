package auth_test

import (
	"testing"

	"bitbucket.org/syb-devs/gotools/auth"
)

func TestCheck(t *testing.T) {
	a := auth.New()
	username, plain := "john doe", []byte("7h1$ 1$ 50m37h!n6")
	a.SetUsername(username)
	a.GeneratePassword(plain)

	err := a.Check(username, plain)

	if err != nil {
		t.Errorf("checking user / password: %v", err)
	}
}

func TestCheckInvalidUser(t *testing.T) {

	a := auth.New()
	username, plain := "john doe", []byte("7h1$ 1$ 50m37h!n6")
	a.SetUsername(username)
	a.GeneratePassword(plain)

	err := a.Check("john williams", plain)

	if err != auth.ErrInvalidUserName {
		t.Errorf("expecting invalid username error, got: %v", err)
	}
}

func TestCheckInvalidPassword(t *testing.T) {

	a := auth.New()
	username, plain := "john doe", []byte("7h1$ 1$ 50m37h!n6")
	a.SetUsername(username)
	a.GeneratePassword(plain)

	err := a.Check(username, []byte("invalid password"))

	if err != auth.ErrInvalidPassword {
		t.Errorf("expecting invalid password error, got: %v", err)
	}
}

func TestCan(t *testing.T) {

	auth.RegisterRole("sysadmin", "manage servers", "blame developers")
	auth.RegisterRole("developer", "code", "blame sysadmins")

	a := auth.New()
	a.AddRole("sysadmin")

	if a.Can("code") {
		t.Error("a sysadmin should not be able to code")
	}

	if !a.Can("manage servers") {
		t.Error("a sysadmin should be able to manage servers")
	}

	a.AddRole("developer")
	if !a.CanAll("manage servers", "blame developers", "code", "blame sysadmins") {
		t.Error("a sysadmin/developer should be able to do almost everything you can imagine")
	}

	a.AddRole("dummy")
	if a.CanAll("travel in time", "fly") {
		t.Error("are you saying a dev/sysadmin can travel in time and fly? seriously?? ")
	}
}
