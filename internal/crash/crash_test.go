package crash

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolvePath_UsesXDGStateHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)
	t.Setenv(envCrashLog, "")

	got, err := resolvePath()
	if err != nil {
		t.Fatalf("resolvePath: %v", err)
	}
	want := filepath.Join(dir, "vibez", "crash.log")
	if got != want {
		t.Fatalf("resolvePath = %q, want %q", got, want)
	}
}

func TestResolvePath_EnvOverride(t *testing.T) {
	override := filepath.Join(t.TempDir(), "custom.log")
	t.Setenv(envCrashLog, override)

	got, err := resolvePath()
	if err != nil {
		t.Fatalf("resolvePath: %v", err)
	}
	if got != override {
		t.Fatalf("resolvePath = %q, want %q", got, override)
	}
}

func TestWriteReport_AppendsToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vibez", "crash.log")
	logPath = path
	appVersion = "test"

	writeReport("panic", "unit", "boom", "goroutine 1 [running]:\n")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	text := string(data)
	for _, want := range []string{
		"=== vibez panic report ===",
		"version: test",
		"component: unit",
		"message: boom",
		"stack:",
		"goroutine 1 [running]:",
		"=== end ===",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("crash log missing %q:\n%s", want, text)
		}
	}
}

func TestRecover_WritesReportAndExits(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "crash.log")
	logPath = path
	appVersion = "test"

	exited := false
	oldExit := osExit
	osExit = func(code int) {
		exited = true
		if code != 2 {
			t.Fatalf("exit code = %d, want 2", code)
		}
	}
	t.Cleanup(func() { osExit = oldExit })

	func() {
		defer Recover("test")
		panic("test panic")
	}()

	if !exited {
		t.Fatal("Recover did not exit")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(data), "test panic") {
		t.Fatalf("crash log missing panic message: %s", data)
	}
}

func TestReportError_SkipsNil(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "crash.log")
	logPath = path

	ReportError("tui", nil)
	if _, err := os.Stat(path); err == nil {
		t.Fatal("expected no crash log for nil error")
	}
}

func TestInstall_CreatesLogAndSetsPath(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)
	t.Setenv(envCrashLog, "")

	if err := Install("0.0.0"); err != nil {
		t.Fatalf("Install: %v", err)
	}
	want := filepath.Join(dir, "vibez", "crash.log")
	if Path() != want {
		t.Fatalf("Path() = %q, want %q", Path(), want)
	}
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("crash log not created: %v", err)
	}
}
