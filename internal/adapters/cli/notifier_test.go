package cli_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/el-j/nix-config-collector/internal/adapters/cli"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestNotifierInfo(t *testing.T) {
	n := cli.New()
	output := captureStdout(func() {
		n.Info("test info message")
	})
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}
	fmt.Println("Info output:", output)
}

func TestNotifierSuccess(t *testing.T) {
	n := cli.New()
	output := captureStdout(func() {
		n.Success("test success")
	})
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestNotifierWarning(t *testing.T) {
	n := cli.New()
	output := captureStdout(func() {
		n.Warning("test warning")
	})
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestNotifierError(t *testing.T) {
	n := cli.New()
	output := captureStdout(func() {
		n.Error("test error")
	})
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestNotifierProgress(t *testing.T) {
	n := cli.New()
	output := captureStdout(func() {
		n.Progress(1, 10, "working")
	})
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestNotifierPrintScanSummary(t *testing.T) {
	n := cli.New()
	output := captureStdout(func() {
		n.PrintScanSummary([][2]string{
			{"Packages", "42"},
			{"Services", "7"},
			{"Hostname", "my-mac"},
		})
	})
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}
	if !strings.Contains(output, "Packages") {
		t.Error("expected output to contain 'Packages'")
	}
	if !strings.Contains(output, "42") {
		t.Error("expected output to contain '42'")
	}
}

func TestNotifierPrintScanSummaryEmpty(t *testing.T) {
	n := cli.New()
	output := captureStdout(func() {
		n.PrintScanSummary([][2]string{})
	})
	if output != "" {
		t.Error("expected empty output for empty rows")
	}
}
