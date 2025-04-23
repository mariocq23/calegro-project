package utilities

import (
	"scripter/entities"
)

type Security struct {
}

func (security Security) ValidateSecurity(signal entities.Signal) bool {
	if signal.BypassSecurity {
		return true
	}
	if signal.AuthenticationHub == "" || signal.AuthorizationHub == "" || signal.CertificationHub == "" {
		return false
	}
	return false
}
