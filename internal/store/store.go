package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	StatusActive   = "active"
	StatusComplete = "complete"
)

var (
	ErrNoGoal      = errors.New("no goal for this session")
	ErrNotActive   = errors.New("goal is not active")
	ErrAlreadyDone = errors.New("goal is already complete")
)

type Goal struct {
	SessionID   string     `json:"session_id"`
	Text        string     `json:"text"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".cursor", "goals")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func pathFor(sessionID string) (string, error) {
	if sessionID == "" {
		return "", errors.New("session id is required (use --session-id or run via Cursor goal hooks)")
	}
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, sessionID+".json"), nil
}

func Load(sessionID string) (*Goal, error) {
	p, err := pathFor(sessionID)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoGoal
		}
		return nil, err
	}
	var g Goal
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func Save(g *Goal) error {
	p, err := pathFor(g.SessionID)
	if err != nil {
		return err
	}
	g.UpdatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

func Set(sessionID, text string) (*Goal, error) {
	now := time.Now().UTC()
	g := &Goal{
		SessionID: sessionID,
		Text:      text,
		Status:    StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if existing, err := Load(sessionID); err == nil {
		g.CreatedAt = existing.CreatedAt
	}
	return g, Save(g)
}

func Edit(sessionID, text string) (*Goal, error) {
	g, err := Load(sessionID)
	if err != nil {
		return nil, err
	}
	if g.Status != StatusActive {
		return nil, ErrNotActive
	}
	g.Text = text
	return g, Save(g)
}

func Complete(sessionID string) (*Goal, error) {
	g, err := Load(sessionID)
	if err != nil {
		return nil, err
	}
	if g.Status == StatusComplete {
		return g, ErrAlreadyDone
	}
	now := time.Now().UTC()
	g.Status = StatusComplete
	g.CompletedAt = &now
	return g, Save(g)
}

func Clear(sessionID string) error {
	p, err := pathFor(sessionID)
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func FormatStatus(g *Goal) string {
	if g.CompletedAt != nil {
		return fmt.Sprintf("status=%s\nsession=%s\ntext=%s\ncompleted_at=%s\n",
			g.Status, g.SessionID, g.Text, g.CompletedAt.Format(time.RFC3339))
	}
	return fmt.Sprintf("status=%s\nsession=%s\ntext=%s\nupdated_at=%s\n",
		g.Status, g.SessionID, g.Text, g.UpdatedAt.Format(time.RFC3339))
}
