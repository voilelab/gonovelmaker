package main

import (
	"testing"
)

func TestGenCurrCmd(t *testing.T) {
	t.Run("requires filepath flag", func(t *testing.T) {
		cmd := NewGenCurrCmd(dummyBackendMaker)
		cmd.cmd.SetArgs([]string{})
		err := cmd.cmd.Execute()
		if err == nil {
			t.Fatal("expected error when filepath flag is missing")
		}
	})

	t.Run("accepts valid filepath", func(t *testing.T) {
		cmd := NewGenCurrCmd(dummyBackendMaker)
		// This will fail at runtime but should pass flag validation
		cmd.cmd.SetArgs([]string{"--filepath", "Story/001_ch1.md"})
		// We can't actually execute this without a proper vault setup
		// so we just verify the flag parsing works
		if err := cmd.cmd.ParseFlags([]string{"--filepath", "Story/001_ch1.md"}); err != nil {
			t.Fatalf("failed to parse valid flags: %v", err)
		}
	})
}
