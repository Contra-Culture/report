package report

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	Kind uint8
	node struct {
		kind     Kind
		allowMap Kind
		nodes    []*node
		msg      string
		timer    Timer
		time     time.Time
		duration *time.Duration
	}
	standardTimer struct {
		beginning time.Time
	}
	dumbTimer struct {
		c         int64
		beginning time.Time
	}
	Node interface {
		Finalize()
		Message() string
		Traverse(fn func([]int, Kind, time.Time, *time.Duration, string) error) error
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
	Timer interface {
		Now() time.Time
		Finalize() time.Duration
		New() (Timer, time.Time)
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

func DumbTimer(now time.Time) *dumbTimer {
	return &dumbTimer{beginning: now}
}
func (t *dumbTimer) Now() time.Time {
	t.c++
	d := time.Duration(t.c * 100)
	return t.beginning.Add(d)
}
func (t *dumbTimer) Finalize() time.Duration {
	t.c++
	return time.Duration(t.c * 100)
}
func (t *dumbTimer) New() (Timer, time.Time) {
	t.c++
	now := t.beginning.Add(time.Duration(t.c * 100))
	return &dumbTimer{
		beginning: now,
	}, now
}
func (t *standardTimer) Now() time.Time {
	return time.Now()
}
func (t *standardTimer) Finalize() time.Duration {
	return time.Since(t.beginning)
}
func (t *standardTimer) New() (Timer, time.Time) {
	n := t.Now()
	return &standardTimer{beginning: n}, n
}
func NewWithTimer(t Timer, m string, injs ...interface{}) Node {
	node := new(m, injs...)
	node.timer = t
	node.time = t.Now()
	return node
}

func ReportCreator(t Timer) func(string, ...interface{}) Node {
	return func(m string, injs ...interface{}) Node {
		return NewWithTimer(t, m, injs...)
	}
}

// New() creates top level reporting node.
func New(m string, injs ...interface{}) Node {
	return new(m, injs...)
}
func new(m string, injs ...interface{}) *node {
	now := time.Now()
	t := &standardTimer{beginning: now}
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	return &node{
		kind:     Structure,
		allowMap: allKinds,
		nodes:    []*node{},
		msg:      m,
		timer:    t,
		time:     now,
	}
}
func (n *node) Finalize() {
	d := n.timer.Finalize()
	n.duration = &d
}

// Traverse() allows to go through all the report nodes starting from the current one.
func (n *node) Traverse(fn func([]int, Kind, time.Time, *time.Duration, string) error) error {
	return n.traverse([]int{}, fn)
}
func (n *node) traverse(path []int, fn func([]int, Kind, time.Time, *time.Duration, string) error) (err error) {
	err = fn(path, n.kind, n.time, n.duration, n.msg)
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
	t, now := n.timer.New()
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	child := &node{
		kind:     Error,
		allowMap: Error,
		msg:      m,
		timer:    t,
		time:     now,
	}
	if (n.allowMap & Error) > 0 {
		n.nodes = append(n.nodes, child)
	}
	return child
}

// Warn() allows to add warning message node to report/log.
func (n *node) Warn(m string, injs ...interface{}) Node {
	t, now := n.timer.New()
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	child := &node{
		kind:     Warn,
		allowMap: Warn,
		msg:      m,
		timer:    t,
		time:     now,
	}
	if (n.allowMap & Warn) > 0 {
		n.nodes = append(n.nodes, child)
	}
	return child
}

// Deprecation() allows to add deprecation message node to report/log.
func (n *node) Deprecation(m string, injs ...interface{}) Node {
	t, now := n.timer.New()
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	child := &node{
		kind:     Deprecation,
		allowMap: Deprecation,
		msg:      m,
		timer:    t,
		time:     now,
	}
	if (n.allowMap & Deprecation) > 0 {
		n.nodes = append(n.nodes, child)
	}
	return child
}

// Info() allows to add informational message node to report/log.
func (n *node) Info(m string, injs ...interface{}) Node {
	t, now := n.timer.New()
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	child := &node{
		kind:     Info,
		allowMap: Info,
		msg:      m,
		timer:    t,
		time:     now,
	}
	if (n.allowMap & Info) > 0 {
		n.nodes = append(n.nodes, child)
	}
	return child
}

// Debug() allows to add debug message node to report/log.
func (n *node) Debug(m string, injs ...interface{}) Node {
	t, now := n.timer.New()
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	child := &node{
		kind:     Debug,
		allowMap: Debug,
		msg:      m,
		timer:    t,
		time:     now,
	}
	if (n.allowMap & Debug) > 0 {
		n.nodes = append(n.nodes, child)
	}
	return child
}

// Structure() allows to add structure message node to report/log.
// Structure nodes allows to label units or architecture levels.
func (n *node) Structure(m string, injs ...interface{}) Node {
	t, now := n.timer.New()
	if len(injs) > 0 {
		m = fmt.Sprintf(m, injs...)
	}
	child := &node{
		kind:     Structure,
		allowMap: allKinds,
		msg:      m,
		timer:    t,
		time:     now,
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
	if n.duration != nil {
		return fmt.Sprintf("[%s - %sns] %s %s", n.time.Format(time.RFC3339Nano), n.duration, kindString(n.kind), n.msg)
	}
	return fmt.Sprintf("[%s] %s %s", n.time.Format(time.RFC3339Nano), kindString(n.kind), n.msg)
}

// HasErrors() returns true if current node or one of its children has/is an error message.
func (n *node) HasErrors() (hasErrors bool) {
	hasErrors = false
	n.Traverse(
		func(_ []int, k Kind, _ time.Time, _ *time.Duration, _ string) (err error) {
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
		func(_ []int, k Kind, _ time.Time, _ *time.Duration, _ string) (err error) {
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
		func(_ []int, k Kind, _ time.Time, _ *time.Duration, _ string) (err error) {
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
		func(path []int, k Kind, t time.Time, d *time.Duration, m string) (err error) {
			for range path {
				sb.WriteRune('\t')
			}
			sb.WriteString(kindString(k))
			sb.WriteRune('[')
			sb.WriteString(t.Format(time.RFC3339Nano))
			if d != nil {
				sb.WriteRune(' ')
				dur := int64(*d)
				sb.WriteString(strconv.FormatInt(dur, 10))
				sb.WriteString("ns")
			}
			sb.WriteRune(']')
			sb.WriteRune(' ')
			sb.WriteString(m)
			sb.WriteRune('\n')
			return
		})
	return sb.String()
}
func ToError(n Node) error {
	var sb strings.Builder
	sb.WriteString("multiple errors:\n")
	n.Traverse(
		func(path []int, k Kind, t time.Time, d *time.Duration, m string) (err error) {
			for range path {
				sb.WriteRune('\t')
			}
			switch k {
			case Error:
				sb.WriteString("\nerror: ")
				sb.WriteString(m)
				sb.WriteRune('\n')
			case Structure:
				sb.WriteRune('\n')
				sb.WriteString(m)
				sb.WriteRune('\n')
			default:
				// do nothing
			}
			return
		})
	errMsg := sb.String()
	if len(errMsg) == 0 {
		return nil
	}
	return errors.New(errMsg)
}
func kindString(k Kind) string {
	switch k {
	case Structure:
		return "#"
	case Error:
		return "<error>"
	case Info:
		return "<info>"
	case Debug:
		return "<debug>"
	case Warn:
		return "<warning>"
	case Deprecation:
		return "<deprecated>"
	default:
		panic(fmt.Sprintf("wrong node kind - %#v", k)) // should not occure
	}
}
