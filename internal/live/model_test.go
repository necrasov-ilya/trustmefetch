package live

import (
	"strings"
	"testing"

	"github.com/necrasov-ilya/trustmefetch/internal/system"
	"github.com/necrasov-ilya/trustmefetch/internal/theme"
)

func TestQuestionModePinsAnswerInHeader(t *testing.T) {
	model := New(system.Info{}, theme.Must("rgb-linux"), false, true, true)
	if view := model.View().Content; !strings.Contains(view, "YES · 100% LINUX") {
		t.Fatalf("question answer is missing from live header:\n%s", view)
	}
}

func TestRegularLiveModeHasNoQuestionAnswer(t *testing.T) {
	model := New(system.Info{}, theme.Must("rgb-linux"), false, true, false)
	if view := model.View().Content; strings.Contains(view, "YES · 100% LINUX") {
		t.Fatalf("regular live mode contains question answer:\n%s", view)
	}
}
