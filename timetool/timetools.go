package timetool

import "fmt"

type TimeTools struct {
}

func NewTimeTools() *TimeTools {
	return &TimeTools{}
}

func (t *TimeTools) Hello() {
	fmt.Println("hello world!")
}
