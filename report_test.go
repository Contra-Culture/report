package report_test

import (
	"fmt"
	"time"

	. "github.com/Contra-Culture/report"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type dumbTimer struct {
	c         int64
	beginning time.Time
}

func DumbTimer(now time.Time) *dumbTimer {
	return &dumbTimer{beginning: now}
}
func (t *dumbTimer) Now() time.Time {
	t.c++
	d := time.Duration(t.c * 100)
	return t.beginning.Add(d)
}

var _ = Describe("report", func() {
	Describe("test timer", func() {
		Describe(".Now()", func() {
			It("returns current time", func() {
				now := time.Now()
				t := DumbTimer(now)
				n := t.Now()
				Expect(n.Sub(now)).To(Equal(time.Duration(100)))
				n = t.Now()
				Expect(n.Sub(now)).To(Equal(time.Duration(200)))
				n = t.Now()
				Expect(n.Sub(now)).To(Equal(time.Duration(300)))
			})
		})
	})
	Describe("creation", func() {
		Describe("New()", func() {
			It("creates root context", func() {
				now := time.Now()
				r := NewWithTimer(DumbTimer(now), "test")
				Expect(r).NotTo(BeNil())
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n", now.Add(100).Format(time.RFC3339Nano))))
				// with string template
				now = time.Now()
				r = NewWithTimer(DumbTimer(now), "test %s", "app")
				Expect(r).NotTo(BeNil())
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test app\n", now.Add(100).Format(time.RFC3339Nano))))
			})
		})
		Describe(".Structure()", func() {
			It("creates child context", func() {
				now := time.Now()
				r := NewWithTimer(DumbTimer(now), "test")
				child := r.Structure("child")
				Expect(child).NotTo(BeNil())
				Expect(ToString(child)).To(Equal(fmt.Sprintf("#[%s] child\n", now.Add(200).Format(time.RFC3339Nano))))
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n\t#[%s] child\n", now.Add(100).Format(time.RFC3339Nano), now.Add(200).Format(time.RFC3339Nano))))
				// with string template
				now = time.Now()
				r = NewWithTimer(DumbTimer(now), "test")
				child = r.Structure("child: %s", "someContext")
				Expect(child).NotTo(BeNil())
				Expect(ToString(child)).To(Equal(fmt.Sprintf("#[%s] child: someContext\n", now.Add(200).Format(time.RFC3339Nano))))
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n\t#[%s] child: someContext\n", now.Add(100).Format(time.RFC3339Nano), now.Add(200).Format(time.RFC3339Nano))))
			})
		})
		Describe(".Error()", func() {
			It("adds error record", func() {
				now := time.Now()
				r := NewWithTimer(DumbTimer(now), "test")
				r.Error("some error")
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n\t<error>[%s] some error\n", now.Add(100).Format(time.RFC3339Nano), now.Add(200).Format(time.RFC3339Nano))))
				// with string template
				now = time.Now()
				r = NewWithTimer(DumbTimer(now), "test")
				r.Error("some error: %s", "bad error")
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n\t<error>[%s] some error: bad error\n", now.Add(100).Format(time.RFC3339Nano), now.Add(200).Format(time.RFC3339Nano))))
			})
		})
		Describe(".Warn()", func() {
			It("adds warn record", func() {
				now := time.Now()
				r := NewWithTimer(DumbTimer(now), "test")
				r.Warn("some warn")
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n\t<warning>[%s] some warn\n", now.Add(100).Format(time.RFC3339Nano), now.Add(200).Format(time.RFC3339Nano))))
				// with string template
				now = time.Now()
				r = NewWithTimer(DumbTimer(now), "test")
				r.Warn("some warn: %s", "better don't do that")
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n\t<warning>[%s] some warn: better don't do that\n", now.Add(100).Format(time.RFC3339Nano), now.Add(200).Format(time.RFC3339Nano))))
			})
		})
		Describe(".Info()", func() {
			It("adds info record", func() {
				now := time.Now()
				r := NewWithTimer(DumbTimer(now), "test")
				r.Info("some info")
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n\t<info>[%s] some info\n", now.Add(100).Format(time.RFC3339Nano), now.Add(200).Format(time.RFC3339Nano))))
				// with string template
				now = time.Now()
				r = NewWithTimer(DumbTimer(now), "test")
				r.Info("some info: %s", "useful info")
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n\t<info>[%s] some info: useful info\n", now.Add(100).Format(time.RFC3339Nano), now.Add(200).Format(time.RFC3339Nano))))
			})
		})
		Describe(".Debug()", func() {
			It("adds debug record", func() {
				now := time.Now()
				r := NewWithTimer(DumbTimer(now), "test")
				r.Debug("some debug info")
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n\t<debug>[%s] some debug info\n", now.Add(100).Format(time.RFC3339Nano), now.Add(200).Format(time.RFC3339Nano))))
				// with string template
				now = time.Now()
				r = NewWithTimer(DumbTimer(now), "test")
				r.Debug("some debug info: %s", "here is the bug")
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n\t<debug>[%s] some debug info: here is the bug\n", now.Add(100).Format(time.RFC3339Nano), now.Add(200).Format(time.RFC3339Nano))))
			})
		})
		Describe(".Deprecation()", func() {
			It("adds deprecation record", func() {
				now := time.Now()
				r := NewWithTimer(DumbTimer(now), "test")
				r.Deprecation("some deprecation")
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n\t<deprecated>[%s] some deprecation\n", now.Add(100).Format(time.RFC3339Nano), now.Add(200).Format(time.RFC3339Nano))))
				// with template string
				now = time.Now()
				r = NewWithTimer(DumbTimer(now), "test")
				r.Deprecation("some deprecation: %s", "will be removed soon")
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n\t<deprecated>[%s] some deprecation: will be removed soon\n", now.Add(100).Format(time.RFC3339Nano), now.Add(200).Format(time.RFC3339Nano))))
			})
		})
	})
	Describe("presentation", func() {
		Describe(".ToString()", func() {
			It("returns string presentation", func() {
				now := time.Now()
				r := NewWithTimer(DumbTimer(now), "test")
				r.Info("root-info")
				r.Error("root-error")
				r1 := r.Structure("child1")
				r1.Error("child1-error")
				r2 := r.Structure("child2")
				r2.Info("child2-info")
				Expect(ToString(r)).To(Equal(fmt.Sprintf("#[%s] test\n\t<info>[%s] root-info\n\t<error>[%s] root-error\n\t#[%s] child1\n\t\t<error>[%s] child1-error\n\t#[%s] child2\n\t\t<info>[%s] child2-info\n", now.Add(100).Format(time.RFC3339Nano), now.Add(200).Format(time.RFC3339Nano), now.Add(300).Format(time.RFC3339Nano), now.Add(400).Format(time.RFC3339Nano), now.Add(500).Format(time.RFC3339Nano), now.Add(600).Format(time.RFC3339Nano), now.Add(700).Format(time.RFC3339Nano))))
			})
		})
	})
	Describe("predicates", func() {
		Describe(".HasErrors()", func() {
			Context("when has errors", func() {
				It("returns true", func() {
					r := New("test")
					r.Structure("nested").Error("some error")
					Expect(r.HasErrors()).To(BeTrue())
				})
			})
			Context("when has no errors", func() {
				It("returns true", func() {
					r := New("test")
					r.Structure("nested").Info("some info")
					Expect(r.HasErrors()).To(BeFalse())
				})
			})
		})
		Describe(".HasWarns()", func() {
			Context("when has warns", func() {
				It("returns true", func() {
					r := New("test")
					r.Structure("nested").Warn("some warn")
					Expect(r.HasWarns()).To(BeTrue())
				})
			})
			Context("when has no warns", func() {
				It("returns true", func() {
					r := New("test")
					r.Structure("nested").Info("some info")
					Expect(r.HasWarns()).To(BeFalse())
				})
			})
		})
		Describe(".HasDeprecations()", func() {
			Context("when has deprecations", func() {
				It("returns true", func() {
					r := New("test")
					r.Structure("nested").Deprecation("some deprecation")
					Expect(r.HasDeprecations()).To(BeTrue())
				})
			})
			Context("when has no deprecations", func() {
				It("returns true", func() {
					r := New("test")
					r.Structure("nested").Info("some info")
					Expect(r.HasDeprecations()).To(BeFalse())
				})
			})
		})
	})
})
