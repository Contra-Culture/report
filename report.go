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
func (n *node) Allow(kinds ...Kind) {
	for _, k := range kinds {
		n.allowMap |= k
	}
}
func (n *node) Message() string {
	return ""
}
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
func ToString(n Node) string {
	var sb strings.Builder
	n.Traverse(
		func(path []int, k Kind, m string) (err error) {
			for range path {
				sb.WriteRune('\t')
			}
			switch k {
			case Structure:
				sb.WriteString("| ")
			case Error:
				sb.WriteString("[ error ] ")
			case Info:
				sb.WriteString("[ info ] ")
			case Debug:
				sb.WriteString("[ debug ] ")
			case Warn:
				sb.WriteString("[ warning ] ")
			case Deprecation:
				sb.WriteString("[ deprecated ] ")
			default:
				panic(fmt.Sprintf("wrong node kind - %#v", k))
			}
			sb.WriteString(m)
			sb.WriteRune('\n')
			return
		})
	return sb.String()
}
