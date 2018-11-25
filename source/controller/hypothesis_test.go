package controller

import (
	"testing"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

func TestSomeStuff(t *testing.T) {
	domainName := `mada.wp.co.uk`
	domainClean, err := publicsuffix.Domain(domainName)
	if err != nil {
		t.Error(err)
	}
	if domainName != domainClean {
		t.Error(
			"Expected " + domainName +
				" got " + domainClean)
	}
}
