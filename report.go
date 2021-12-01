package report

import (
	"fmt"
	"strings"
)

type (
	recordKind int
	RContext   struct {
		depth    int
		children []interface{} // interface{} is *Context or *Record
		title    string
	}
	Record struct {
		kind    recordKind
		message string
	}
)

const (
	_ recordKind = iota
	Error
	Info
	Debug
	Warn
	Deprecation
)

func New(t string) (c *RContext) {
	return &RContext{
		depth:    0,
		title:    t,
		children: []interface{}{},
	}
}
func Newf(t string, injections ...interface{}) (c *RContext) {
	return &RContext{
		depth:    0,
		title:    fmt.Sprintf(t, injections...),
		children: []interface{}{},
	}
}
func (c *RContext) String() string {
	acc := []string{}
	acc = append(acc, c.title)
	acc = append(acc, "\n")
	for _, rawChild := range c.children {
		for i := 0; i <= c.depth; i++ {
			acc = append(acc, "\t")
		}
		switch child := rawChild.(type) {
		case *RContext:
			acc = append(acc, child.String())
		case *Record:
			switch child.kind {
			case Error:
				acc = append(acc, "\t[ error ] ")
			case Info:
				acc = append(acc, "\t[ info ] ")
			case Debug:
				acc = append(acc, "\t[ debug ] ")
			case Warn:
				acc = append(acc, "\t[ warning ] ")
			case Deprecation:
				acc = append(acc, "\t[ deprecated ] ")
			default:
				panic("wrong record kind")
			}
			acc = append(acc, child.message)
			acc = append(acc, "\n")
		default:
			panic("wrong children type")
		}
	}
	return strings.Join(acc, "")
}
func (c *RContext) Error(m string) {
	c.children = append(
		c.children,
		&Record{
			kind:    Error,
			message: m,
		})
}
func (c *RContext) Warn(m string) {
	c.children = append(
		c.children,
		&Record{
			kind:    Warn,
			message: m,
		})
}
func (c *RContext) Deprecation(m string) {
	c.children = append(
		c.children,
		&Record{
			kind:    Deprecation,
			message: m,
		})
}
func (c *RContext) Info(m string) {
	c.children = append(
		c.children,
		&Record{
			kind:    Info,
			message: m,
		})
}
func (c *RContext) Debug(m string) {
	c.children = append(
		c.children,
		&Record{
			kind:    Debug,
			message: m,
		})
}
func (c *RContext) Context(t string) (child *RContext) {
	child = &RContext{
		depth: c.depth + 1,
		title: t,
	}
	c.children = append(c.children, child)
	return child
}
func (c *RContext) Errorf(t string, injections ...interface{}) {
	c.children = append(
		c.children,
		&Record{
			kind:    Error,
			message: fmt.Sprintf(t, injections...),
		})
}
func (c *RContext) Warnf(m string) {
	c.children = append(
		c.children,
		&Record{
			kind:    Warn,
			message: m,
		})
}
func (c *RContext) Deprecationf(t string, injections ...interface{}) {
	c.children = append(
		c.children,
		&Record{
			kind:    Deprecation,
			message: fmt.Sprintf(t, injections...),
		})
}
func (c *RContext) Infof(t string, injections ...interface{}) {
	c.children = append(
		c.children,
		&Record{
			kind:    Info,
			message: fmt.Sprintf(t, injections...),
		})
}
func (c *RContext) Debugf(t string, injections ...interface{}) {
	c.children = append(
		c.children,
		&Record{
			kind:    Debug,
			message: fmt.Sprintf(t, injections...),
		})
}
func (c *RContext) Contextf(t string, injections ...interface{}) (child *RContext) {
	child = &RContext{
		depth: c.depth + 1,
		title: fmt.Sprintf(t, injections...),
	}
	c.children = append(c.children, child)
	return child
}
