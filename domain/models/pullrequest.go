package models

import (
	"fmt"
	"time"
)

// PullRequest is the abstraction of a Pull Request from any provider
type PullRequest struct {
	Orgenization string
	ID           string
	Title        string
	Description  string
	Repository   Hierarchy
	Project      Hierarchy
	Status       string
	MergeStatus  string
	CreatedBy    string
	URL          string
	IsDraft      bool
	Created      time.Time
}

type CreatedPullRequest struct {
	ID           string
	Title        string
	Description  string
	Repository   Hierarchy
	Project      Hierarchy
	URL          string
	IsDraft      bool
	Organization string
}

// Hierarchy of a PR
type Hierarchy struct {
	ID   string
	Name string
	URL  string
}

// ShortenedTitle returns title shortened to maxlength...
func (p PullRequest) ShortenedTitle(maxLength int, isDraft bool) string {

	if len(p.Title) <= maxLength {
		return p.Title
	}

	draftText := ""
	if isDraft {
		draftText = " (Draft)"
	}

	shortenLenght := maxLength - 3 - len(draftText)

	title := fmt.Sprintf("%s...", p.Title[0:shortenLenght])
	return title
}
