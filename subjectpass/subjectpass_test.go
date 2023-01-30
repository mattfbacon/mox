package subjectpass

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/mjl-/mox/smtp"
)

func TestSubjectPass(t *testing.T) {
	key := []byte("secret token")
	addr, _ := smtp.ParseAddress("mox@mox.example")
	sig := Generate(addr, key, time.Now())

	message := fmt.Sprintf("From: <mox@mox.example>\r\nSubject: let me in %s\r\n\r\nthe message", sig)
	if err := Verify(strings.NewReader(message), key, time.Hour); err != nil {
		t.Fatalf("verifyPassToken: %s", err)
	}

	if err := Verify(strings.NewReader(message), []byte("bad key"), time.Hour); err == nil {
		t.Fatalf("verifyPassToken did not fail")
	}

	sig = Generate(addr, key, time.Now().Add(-time.Hour-257))
	message = fmt.Sprintf("From: <mox@mox.example>\r\nSubject: let me in %s\r\n\r\nthe message", sig)
	if err := Verify(strings.NewReader(message), key, time.Hour); !errors.Is(err, ErrExpired) {
		t.Fatalf("verifyPassToken should have expired")
	}
}
