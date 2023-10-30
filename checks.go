package main

import "github.com/fatih/color"

type Check struct {
	Message string
	Points  int
	// points were left unspecified
	PointsEmpty bool

	Condition
	// separate root hint from condition tree
	Hint string
}

func (c *Check) Debug() string {
	cl := color.New(color.Bold)
	ret := cl.Sprintf("%s (%d Points)%s\n", c.Message, c.Points, formatHint(c.Hint))
	return ret + DebugCondition(c.Condition)
}
