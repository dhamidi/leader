package main

import (
	"fmt"
	"io"
	"strings"
)

// BreadcrumbsView displays the path to the currently selected key map.
type BreadcrumbsView struct {
	// Entries on the path to the currently selected key map.
	Breadcrumbs []string
}

// NewBreadcrumbsView creates a new view for the given breadcrumbs
func NewBreadcrumbsView(breadcrumbs []string) *BreadcrumbsView {
	return &BreadcrumbsView{
		Breadcrumbs: breadcrumbs,
	}
}

// Render displays all breadcrumbs separated with a greater-than sign.
func (v *BreadcrumbsView) Render(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s\n\r", strings.Join(v.Breadcrumbs, " > "))
	return err
}

// Erase erases this view by deleting one line
func (v *BreadcrumbsView) Erase(out io.Writer) error {
	_, err := fmt.Fprintf(out, "\033[A\033[2K")
	return err
}
