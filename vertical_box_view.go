package main

import "io"

// VerticalBoxView renders child views in the provided order.
type VerticalBoxView struct {
	children []View
}

// NewVerticalBoxView creates a new vertical box view for rendering
// children in the provided order.
func NewVerticalBoxView(children ...View) *VerticalBoxView {
	return &VerticalBoxView{
		children: children,
	}
}

// Render implements View by calling View for each of its children.
//
// Rendering is aborted on the first error, which is returned.
func (v *VerticalBoxView) Render(out io.Writer) error {
	for _, child := range v.children {
		err := child.Render(out)
		if err != nil {
			return err
		}
	}
	return nil
}

// Erase implements View by calling Erase for each child view in reverse order.
func (v *VerticalBoxView) Erase(out io.Writer) error {
	for i := len(v.children) - 1; i >= 0; i-- {
		err := v.children[i].Erase(out)
		if err != nil {
			return err
		}
	}

	return nil
}
