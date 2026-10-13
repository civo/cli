package utility

import (
	"testing"

	"github.com/civo/cli/common"
)

func TestAskForConfirm_QuietModeFailsWithoutPrompting(t *testing.T) {
	originalQuiet := common.Quiet
	originalRetrieveUserInput := retrieveUserInput
	defer func() {
		common.Quiet = originalQuiet
		retrieveUserInput = originalRetrieveUserInput
	}()

	common.Quiet = true

	promptWasCalled := false
	retrieveUserInput = func(message string) (string, error) {
		promptWasCalled = true
		return "yes", nil
	}

	err := AskForConfirm("delete the thing")
	if err == nil {
		t.Fatal("expected AskForConfirm to fail in --quiet mode without --yes, got nil error")
	}
	if promptWasCalled {
		t.Fatal("AskForConfirm must not read interactive input at all when --quiet is set, but it did")
	}
}

func TestAskForConfirm_NonQuietModeStillPrompts(t *testing.T) {
	originalQuiet := common.Quiet
	originalRetrieveUserInput := retrieveUserInput
	defer func() {
		common.Quiet = originalQuiet
		retrieveUserInput = originalRetrieveUserInput
	}()

	common.Quiet = false

	tests := []struct {
		name      string
		answer    string
		wantError bool
	}{
		{name: "yes answer confirms", answer: "yes", wantError: false},
		{name: "y answer confirms", answer: "y", wantError: false},
		{name: "no answer does not confirm", answer: "no", wantError: true},
		{name: "empty answer does not confirm", answer: "", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			promptWasCalled := false
			retrieveUserInput = func(message string) (string, error) {
				promptWasCalled = true
				return tt.answer, nil
			}

			err := AskForConfirm("delete the thing")
			if !promptWasCalled {
				t.Fatal("expected AskForConfirm to prompt for input when --quiet is not set")
			}
			if (err != nil) != tt.wantError {
				t.Fatalf("AskForConfirm(%q) error = %v, wantError %v", tt.answer, err, tt.wantError)
			}
		})
	}
}

func TestUserConfirmedDeletion_QuietModeRespectsYesFlag(t *testing.T) {
	originalQuiet := common.Quiet
	defer func() { common.Quiet = originalQuiet }()
	common.Quiet = true

	// ignoringConfirmed=true simulates --yes being passed: should succeed
	// without ever touching stdin/AskForConfirm, matching "combine with
	// --yes to override" from the issue's acceptance criteria.
	if !UserConfirmedDeletion("instance", true, "my-instance") {
		t.Fatal("expected deletion to be confirmed automatically when --yes is set, even with --quiet")
	}

	// ignoringConfirmed=false (no --yes) with --quiet must fail, not block.
	if UserConfirmedDeletion("instance", false, "my-instance") {
		t.Fatal("expected deletion to be denied when --quiet is set without --yes")
	}
}
