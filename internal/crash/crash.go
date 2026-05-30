// Package crash records fatal errors and panics to stderr and a persistent log file.
package crash

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"
)

const envCrashLog = "VIBEZ_CRASH_LOG"

var (
	mu        sync.Mutex
	logPath   string
	appVersion string
	osExit    = os.Exit
)

// Install configures crash logging for the process. It redirects the Go runtime's
// fatal crash output to the crash log and should be called once at startup.
func Install(version string) error {
	appVersion = version
	path, err := resolvePath()
	if err != nil {
		return err
	}
	logPath = path

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("creating crash log dir: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600) //nolint:gosec // user state dir
	if err != nil {
		return fmt.Errorf("opening crash log: %w", err)
	}

	if err := debug.SetCrashOutput(f, debug.CrashOptions{}); err != nil {
		_ = f.Close()
		return fmt.Errorf("setting crash output: %w", err)
	}
	return nil
}

// Path returns the crash log file path after Install, or the resolved default path.
func Path() string {
	if logPath != "" {
		return logPath
	}
	path, err := resolvePath()
	if err != nil {
		return ""
	}
	return path
}

// Recover catches a panic, writes a crash report, prints a short message to stderr,
// and exits with status 2. Use as: defer crash.Recover("component").
func Recover(component string) {
	if r := recover(); r != nil {
		writeReport("panic", component, fmt.Sprint(r), string(debug.Stack()))
		osExit(2)
	}
}

// ReportError records a non-fatal exit error (for example from tea.Program.Run).
func ReportError(component string, err error) {
	if err == nil {
		return
	}
	writeReport("error", component, err.Error(), "")
}

func resolvePath() (string, error) {
	if override := os.Getenv(envCrashLog); override != "" {
		return override, nil
	}
	stateHome, err := stateHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(stateHome, "vibez", "crash.log"), nil
}

func stateHomeDir() (string, error) {
	if dir := os.Getenv("XDG_STATE_HOME"); dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home directory: %w", err)
	}
	return filepath.Join(home, ".local", "state"), nil
}

func writeReport(kind, component, message, stack string) {
	mu.Lock()
	defer mu.Unlock()

	path := Path()
	if path == "" {
		fmt.Fprintf(os.Stderr, "vibez %s (%s): %s\n", kind, component, message)
		if stack != "" {
			fmt.Fprint(os.Stderr, stack)
		}
		return
	}

	ts := time.Now().UTC().Format(time.RFC3339)
	body := fmt.Sprintf(
		"=== vibez %s report ===\n"+
			"time: %s\n"+
			"version: %s\n"+
			"component: %s\n"+
			"message: %s\n",
		kind, ts, appVersion, component, message,
	)
	if stack != "" {
		body += "\nstack:\n" + stack + "\n"
	}
	body += "=== end ===\n\n"

	if err := appendFile(path, body); err != nil {
		fmt.Fprintf(os.Stderr, "vibez %s (%s): %s (also failed writing crash log: %v)\n", kind, component, message, err)
		if stack != "" {
			fmt.Fprint(os.Stderr, stack)
		}
		return
	}

	fmt.Fprintf(os.Stderr, "vibez %s in %s; details written to %s\n", kind, component, path)
}

func appendFile(path, body string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600) //nolint:gosec // user state dir
	if err != nil {
		return err
	}
	_, err = f.WriteString(body)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	return closeErr
}
