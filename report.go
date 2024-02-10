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

func New(name string) *Reporter {
	return &Reporter{
		createdAt: time.Now(),
		name:      name,
		messages:  []message{},
		reporters: []*Reporter{},
	}
}
func (r *Reporter) Sub(name string) *Reporter {
	sub := New(name)
	r.reporters = append(r.reporters, sub)
	return sub
}
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
