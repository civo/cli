package common

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestGithubError(t *testing.T) {
	errorResponse := &github.ErrorResponse{Response: nil}
	twoFactor := (*github.TwoFactorAuthError)(errorResponse)
	if !IsGHError(twoFactor) {
		t.Fail()
	}
	rateLimit := github.RateLimitError{}
	if !IsGHError(&rateLimit) {
		t.Fail()
	}
	abuseRateLimit := github.AbuseRateLimitError{}
	if !IsGHError(&abuseRateLimit) {
		t.Fail()
	}
}
