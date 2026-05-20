package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/brandongalang/cursor-goal/internal/hooks"
	"github.com/brandongalang/cursor-goal/internal/store"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "hook":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: goal hook <pre-tool-use|stop>")
			os.Exit(2)
		}
		var err error
		switch os.Args[2] {
		case "pre-tool-use":
			err = hooks.PreToolUse()
		case "stop":
			err = hooks.Stop()
		default:
			fmt.Fprintf(os.Stderr, "unknown hook: %s\n", os.Args[2])
			os.Exit(2)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "hook error: %v\n", err)
			os.Exit(1)
		}
	case "set", "edit", "complete", "clear", "status":
		if err := runCommand(os.Args[1], os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `cursor-goal — session goal enforcement for Cursor agents

Usage:
  goal set    --session-id <id> <objective text>
  goal edit   --session-id <id> <objective text>
  goal complete --session-id <id>
  goal clear  --session-id <id>
  goal status --session-id <id>
  goal hook pre-tool-use   (Cursor preToolUse hook; reads JSON from stdin)
  goal hook stop           (Cursor stop hook; reads JSON from stdin)

Install: ./scripts/install.sh
`)
}

type parsedArgs struct {
	sessionID string
	text      string
}

func parseArgs(args []string) (parsedArgs, error) {
	var p parsedArgs
	var textParts []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--session-id":
			if i+1 >= len(args) {
				return p, fmt.Errorf("--session-id requires a value")
			}
			i++
			p.sessionID = args[i]
		default:
			textParts = append(textParts, args[i])
		}
	}
	p.text = strings.TrimSpace(strings.Join(textParts, " "))
	return p, nil
}

func runCommand(cmd string, args []string) error {
	p, err := parseArgs(args)
	if err != nil {
		return err
	}
	if p.sessionID == "" {
		return fmt.Errorf("--session-id is required")
	}
	switch cmd {
	case "set":
		if p.text == "" {
			return fmt.Errorf("objective text is required")
		}
		g, err := store.Set(p.sessionID, p.text)
		if err != nil {
			return err
		}
		fmt.Printf("goal set (active)\nsession=%s\n", g.SessionID)
		return nil
	case "edit":
		if p.text == "" {
			return fmt.Errorf("objective text is required")
		}
		g, err := store.Edit(p.sessionID, p.text)
		if err != nil {
			return err
		}
		fmt.Printf("goal updated (active)\nsession=%s\n", g.SessionID)
		return nil
	case "complete":
		g, err := store.Complete(p.sessionID)
		if err != nil {
			return err
		}
		fmt.Printf("goal complete\nsession=%s\n", g.SessionID)
		return nil
	case "clear":
		if err := store.Clear(p.sessionID); err != nil {
			return err
		}
		fmt.Printf("goal cleared\nsession=%s\n", p.sessionID)
		return nil
	case "status":
		g, err := store.Load(p.sessionID)
		if err != nil {
			return err
		}
		fmt.Print(store.FormatStatus(g))
		return nil
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}
