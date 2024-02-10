package report

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	Reporter struct {
		createdAt time.Time
		name      string
		messages  []message
		reporters []*Reporter
	}
	MessageKind uint8
	message     struct {
		createdAt time.Time
		kind      MessageKind
		message   string
		payload   map[string]interface{}
	}
)

// Unix Syslog + TRACE types
const (
	_ MessageKind = iota
	EMERG
	ALERT
	CRIT
	ERROR
	WARN
	NOTICE
	INFO
	DEBUG
	TRACE
)

var mtStrings = map[MessageKind]string{
	EMERG:  "EMERG",
	ALERT:  "ALERT",
	CRIT:   "CRIT",
	ERROR:  "ERROR",
	WARN:   "WARN",
	NOTICE: "NOTICE",
	INFO:   "INFO",
	DEBUG:  "DEBUG",
	TRACE:  "TRACE",
}

// New() - creates top level reporter.
// You may want to create reporters for HTTP request handlers and then create sub-repporters to debug,
// trace or just to gather structured logs from nested contexts, like HTML rendering, DB queries, 3rd-party API calls.
//
// name string - name of the root context. For example, for HTTP handlers it could be ReqestID - randomly generated string identifier.
func New(name string) *Reporter {
	return &Reporter{
		createdAt: time.Now(),
		name:      name,
		messages:  []message{},
		reporters: []*Reporter{},
	}
}

// *Report.Sub() creates sub (nested) reporter which you can put into some child context.
//
// + name string - name of the child context which will use sub(nested) report.
func (r *Reporter) Sub(name string) *Reporter {
	sub := New(name)
	r.reporters = append(r.reporters, sub)
	return sub
}

// *Reporter.Msg() Writes messages within current report.
//
// + kind MessageKind - enum kind of message (Unix Syslog + Trace)
// + msg - text message
// + payload map[string]interface{} - interface{} could be string, int, float64 or bool.
func (r *Reporter) Msg(kind MessageKind, msg string, payload map[string]interface{}) {
	r.messages = append(
		r.messages,
		message{
			createdAt: time.Now(),
			kind:      kind,
			message:   msg,
			payload:   payload,
		})
}

// *Reporter.JSON() generates JSON report out of current reporter and all its children, grandchildren and etc.
func (r *Reporter) JSON() string {
	var sb strings.Builder
	sb.WriteString("{\"t\":\"")
	sb.WriteString(strconv.Itoa(int(r.createdAt.Unix())))
	sb.WriteString("\",\"s\":\"")
	sb.WriteString(r.name)
	sb.WriteString("\",\"")
	sb.WriteString("m\":[")
	for i, m := range r.messages {
		if i != 0 {
			sb.WriteByte(',')
		}
		sb.WriteString("{\"t\":\"")
		sb.WriteString(strconv.Itoa(int(m.createdAt.Unix())))
		sb.WriteString("\",\"k\":\"")
		sb.WriteString(mtStrings[m.kind])
		sb.WriteString("\",\"m\":\"")
		sb.WriteString(m.message)
		sb.WriteString("\",\"p\":{")
		i := 0
		for k, v := range m.payload {
			if i > 0 {
				sb.WriteRune(',')
			}
			sb.WriteRune('"')
			sb.WriteString(k)
			sb.WriteString("\":")
			switch val := v.(type) {
			case int:
				sb.WriteString(strconv.Itoa(val))
			case float64:
				sb.WriteString(strconv.FormatFloat(val, 'e', 10, 64))
			case string:
				sb.WriteRune('"')
				sb.WriteString(val)
				sb.WriteRune('"')
			case bool:
				sb.WriteString(strconv.FormatBool(val))
			}
			i++
		}
		sb.WriteString("}}")
	}
	if len(r.reporters) > 0 {
		sb.WriteString("],\"n\":[")
		for i, r := range r.reporters {
			if i > 0 {
				sb.WriteRune(',')
			}
			sb.WriteString(r.JSON())
		}
	}
	sb.WriteString("]}")
	return sb.String()
}
func (r *Reporter) Print() {
	_, _ = fmt.Print("\n\n", r.JSON(), "\n\n")
}
