package report

import (
	"fmt"
	"strings"
)

type (
	Kind uint8
	node struct {
		kind     Kind
		allowMap Kind
		nodes    []*node
		msg      string
	}
	Node interface {
		Message() string
		Traverse(fn func([]int, Kind, string) error) error
		Structure(string, ...interface{}) Node
		Error(string, ...interface{}) Node
		Info(string, ...interface{}) Node
		Debug(string, ...interface{}) Node
		Warn(string, ...interface{}) Node
		Deprecation(string, ...interface{}) Node
		Allow(...Kind)
		HasErrors() bool
		HasWarns() bool
		HasDeprecations() bool
	}
)

const (
	Structure Kind = 1 << iota
	Error
	Info
	Debug
	Warn
	Deprecation
)

const allKinds = Structure | Error | Info | Debug | Warn | Deprecation

// New() creates top level reporting node.
func New(m string, injs ...interface{}) Node {
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	return &node{
		kind:     Structure,
		allowMap: allKinds,
		nodes:    []*node{},
		msg:      m,
	}
}

// Traverse() allows to go through all the report nodes starting from the current one.
func (n *node) Traverse(fn func([]int, Kind, string) error) error {
	return n.traverse([]int{}, fn)
}
func (n *node) traverse(path []int, fn func([]int, Kind, string) error) (err error) {
	err = fn(path, n.kind, n.msg)
	if err != nil {
		return
	}
	for i, ch := range n.nodes {
		err = ch.traverse(append(path, i), fn)
		if err != nil {
			return
		}
	}
	return
}

// Error() allows to add error message node to report/log.
func (n *node) Error(m string, injs ...interface{}) Node {
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	child := &node{
		kind:     Error,
		allowMap: Error,
		msg:      m,
	}
	if (n.allowMap & Error) > 0 {
		n.nodes = append(n.nodes, child)
	}
	return child
}

// Warn() allows to add warning message node to report/log.
func (n *node) Warn(m string, injs ...interface{}) Node {
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	child := &node{
		kind:     Warn,
		allowMap: Warn,
		msg:      m,
	}
	if (n.allowMap & Warn) > 0 {
		n.nodes = append(n.nodes, child)
	}
	return child
}

// Deprecation() allows to add deprecation message node to report/log.
func (n *node) Deprecation(m string, injs ...interface{}) Node {
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	child := &node{
		kind:     Deprecation,
		allowMap: Deprecation,
		msg:      m,
	}
	if (n.allowMap & Deprecation) > 0 {
		n.nodes = append(n.nodes, child)
	}
	return child
}

// Info() allows to add informational message node to report/log.
func (n *node) Info(m string, injs ...interface{}) Node {
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	child := &node{
		kind:     Info,
		allowMap: Info,
		msg:      m,
	}
	if (n.allowMap & Info) > 0 {
		n.nodes = append(n.nodes, child)
	}
	return child
}

// Debug() allows to add debug message node to report/log.
func (n *node) Debug(m string, injs ...interface{}) Node {
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	child := &node{
		kind:     Debug,
		allowMap: Debug,
		msg:      m,
	}
	if (n.allowMap & Debug) > 0 {
		n.nodes = append(n.nodes, child)
	}
	return child
}

// Structure() allows to add structure message node to report/log.
// Structure nodes allows to label units or architecture levels.
func (n *node) Structure(m string, injs ...interface{}) Node {
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	child := &node{
		kind:     Structure,
		allowMap: allKinds,
		msg:      m,
	}
	if (n.allowMap & Structure) > 0 {
		n.nodes = append(n.nodes, child)
	}
	return child
}

// Allow() allows to manage what message node kinds are acceptable for the parent node which is a method receiver.
func (n *node) Allow(kinds ...Kind) {
	for _, k := range kinds {
		n.allowMap |= k
	}
}

// Message() returns current node message with its type.
func (n *node) Message() string {
	return fmt.Sprintf("%s %s", kindString(n.kind), n.msg)
}

// HasErrors() returns true if current node or one of its children has/is an error message.
func (n *node) HasErrors() (hasErrors bool) {
	hasErrors = false
	n.Traverse(
		func(_ []int, k Kind, _ string) (err error) {
			if k == Error {
				hasErrors = true
				return
			}
			return
		})
	return
}

// HasWarns() returns true if current node or one of its children has/is a warning message.
func (n *node) HasWarns() (hasWarns bool) {
	hasWarns = false
	n.Traverse(
		func(_ []int, k Kind, _ string) (err error) {
			if k == Warn {
				hasWarns = true
				return
			}
			return
		})
	return
}

// HasDeprecations() returns true if current node or one of its children has/is a deprecation message.
func (n *node) HasDeprecations() (hasDeprecations bool) {
	hasDeprecations = false
	n.Traverse(
		func(_ []int, k Kind, _ string) (err error) {
			if k == Deprecation {
				hasDeprecations = true
				return
			}
			return
		})
	return
}

// ToString() returns string representation of the given node and all its children.
func ToString(n Node) string {
	var sb strings.Builder
	n.Traverse(
		func(path []int, k Kind, m string) (err error) {
			for range path {
				sb.WriteRune('\t')
			}
			sb.WriteString(kindString(k))
			sb.WriteRune(' ')
			sb.WriteString(m)
			sb.WriteRune('\n')
			return
		})
	return sb.String()
}
func kindString(k Kind) string {
	switch k {
	case Structure:
		return "|"
	case Error:
		return "[ error ]"
	case Info:
		return "[ info ]"
	case Debug:
		return "[ debug ]"
	case Warn:
		return "[ warning ]"
	case Deprecation:
		return "[ deprecated ]"
	default:
		panic(fmt.Sprintf("wrong node kind - %#v", k)) // should not occure
	}
}
