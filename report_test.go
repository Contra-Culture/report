package report_test

import (
	. "github.com/Contra-Culture/report"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("report", func() {
	Describe("creation", func() {
		Describe("New()", func() {
			It("creates root context", func() {
				r := New("test")
				Expect(r).NotTo(BeNil())
				Expect(ToString(r)).To(Equal("| test\n"))
				// with string template
				r = New("test %s", "app")
				Expect(r).NotTo(BeNil())
				Expect(ToString(r)).To(Equal("| test app\n"))
			})
		})
		Describe(".Structure()", func() {
			It("creates child context", func() {
				r := New("test")
				child := r.Structure("child")
				Expect(child).NotTo(BeNil())
				Expect(ToString(child)).To(Equal("| child\n"))
				Expect(ToString(r)).To(Equal("| test\n\t| child\n"))
				// with string template
				r = New("test")
				child = r.Structure("child: %s", "someContext")
				Expect(child).NotTo(BeNil())
				Expect(ToString(child)).To(Equal("| child: someContext\n"))
				Expect(ToString(r)).To(Equal("| test\n\t| child: someContext\n"))
			})
		})
		Describe(".Error()", func() {
			It("adds error record", func() {
				r := New("test")
				r.Error("some error")
				Expect(ToString(r)).To(Equal("| test\n\t[ error ] some error\n"))
				// with string template
				r = New("test")
				r.Error("some error: %s", "bad error")
				Expect(ToString(r)).To(Equal("| test\n\t[ error ] some error: bad error\n"))
			})
		})
		Describe(".Warn()", func() {
			It("adds warn record", func() {
				r := New("test")
				r.Warn("some warn")
				Expect(ToString(r)).To(Equal("| test\n\t[ warning ] some warn\n"))
				// with string template
				r = New("test")
				r.Warn("some warn: %s", "better don't do that")
				Expect(ToString(r)).To(Equal("| test\n\t[ warning ] some warn: better don't do that\n"))
			})
		})
		Describe(".Info()", func() {
			It("adds info record", func() {
				r := New("test")
				r.Info("some info")
				Expect(ToString(r)).To(Equal("| test\n\t[ info ] some info\n"))
				// with string template
				r = New("test")
				r.Info("some info: %s", "useful info")
				Expect(ToString(r)).To(Equal("| test\n\t[ info ] some info: useful info\n"))
			})
		})
		Describe(".Debug()", func() {
			It("adds debug record", func() {
				r := New("test")
				r.Debug("some debug info")
				Expect(ToString(r)).To(Equal("| test\n\t[ debug ] some debug info\n"))
				// with string template
				r = New("test")
				r.Debug("some debug info: %s", "here is the bug")
				Expect(ToString(r)).To(Equal("| test\n\t[ debug ] some debug info: here is the bug\n"))
			})
		})
		Describe(".Deprecation()", func() {
			It("adds deprecation record", func() {
				r := New("test")
				r.Deprecation("some deprecation")
				Expect(ToString(r)).To(Equal("| test\n\t[ deprecated ] some deprecation\n"))
				// with template string
				r = New("test")
				r.Deprecation("some deprecation: %s", "will be removed soon")
				Expect(ToString(r)).To(Equal("| test\n\t[ deprecated ] some deprecation: will be removed soon\n"))
			})
		})
	})
	Describe("presentation", func() {
		Describe(".ToString()", func() {
			It("returns string presentation", func() {
				r := New("test")
				r.Info("root-info")
				r.Error("root-error")
				r1 := r.Structure("child1")
				r1.Error("child1-error")
				r2 := r.Structure("child2")
				r2.Info("child2-info")
				Expect(ToString(r)).To(Equal("| test\n\t[ info ] root-info\n\t[ error ] root-error\n\t| child1\n\t\t[ error ] child1-error\n\t| child2\n\t\t[ info ] child2-info\n"))
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
