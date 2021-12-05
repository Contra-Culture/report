package report_test

import (
	. "github.com/Contra-Culture/report"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("report", func() {
	Describe("Context", func() {
		Describe("creation", func() {
			Describe("New()", func() {
				It("creates root context", func() {
					r := New("test")
					Expect(r).NotTo(BeNil())
					Expect(r.String()).To(Equal("test\n"))
				})
			})
			Describe("Newf()", func() {
				It("creates root context", func() {
					r := Newf("test %s", "app")
					Expect(r).NotTo(BeNil())
					Expect(r.String()).To(Equal("test app\n"))
				})
			})
			Describe(".Context()", func() {
				It("creates child context", func() {
					r := New("test")
					child := r.Context("child")
					Expect(child).NotTo(BeNil())
					Expect(child.String()).To(Equal("child\n"))
					Expect(r.String()).To(Equal("test\n\tchild\n"))
				})
			})
			Describe(".Contextf()", func() {
				It("creates child context", func() {
					r := New("test")
					child := r.Contextf("child: %s", "someContext")
					Expect(child).NotTo(BeNil())
					Expect(child.String()).To(Equal("child: someContext\n"))
					Expect(r.String()).To(Equal("test\n\tchild: someContext\n"))
				})
			})
			Describe("records", func() {
				Describe(".Error()", func() {
					It("adds error record", func() {
						r := New("test")
						r.Error("some error")
						Expect(r.String()).To(Equal("test\n\t\t[ error ] some error\n"))
					})
				})
				Describe(".Errorf()", func() {
					It("adds error record", func() {
						r := New("test")
						r.Errorf("some error: %s", "bad error")
						Expect(r.String()).To(Equal("test\n\t\t[ error ] some error: bad error\n"))
					})
				})
				Describe(".Warn()", func() {
					It("adds warn record", func() {
						r := New("test")
						r.Warn("some warn")
						Expect(r.String()).To(Equal("test\n\t\t[ warning ] some warn\n"))
					})
				})
				Describe(".Warnf()", func() {
					It("adds warn record", func() {
						r := New("test")
						r.Warnf("some warn: %s", "better don't do that")
						Expect(r.String()).To(Equal("test\n\t\t[ warning ] some warn: better don't do that\n"))
					})
				})
				Describe(".Info()", func() {
					It("adds info record", func() {
						r := New("test")
						r.Info("some info")
						Expect(r.String()).To(Equal("test\n\t\t[ info ] some info\n"))
					})
				})
				Describe(".Infof()", func() {
					It("adds info record", func() {
						r := New("test")
						r.Infof("some info: %s", "useful info")
						Expect(r.String()).To(Equal("test\n\t\t[ info ] some info: useful info\n"))
					})
				})
				Describe(".Debug()", func() {
					It("adds debug record", func() {
						r := New("test")
						r.Debug("some debug info")
						Expect(r.String()).To(Equal("test\n\t\t[ debug ] some debug info\n"))
					})
				})
				Describe(".Debugf()", func() {
					It("adds debug record", func() {
						r := New("test")
						r.Debugf("some debug info: %s", "here is the bug")
						Expect(r.String()).To(Equal("test\n\t\t[ debug ] some debug info: here is the bug\n"))
					})
				})
				Describe(".Deprecation()", func() {
					It("adds deprecation record", func() {
						r := New("test")
						r.Deprecation("some deprecation")
						Expect(r.String()).To(Equal("test\n\t\t[ deprecated ] some deprecation\n"))
					})
				})
				Describe(".Deprecationf()", func() {
					It("adds deprecation record", func() {
						r := New("test")
						r.Deprecationf("some deprecation: %s", "will be removed soon")
						Expect(r.String()).To(Equal("test\n\t\t[ deprecated ] some deprecation: will be removed soon\n"))
					})
				})
			})
		})
		Describe("presentation", func() {
			Describe(".String()", func() {
				It("returns string presentation", func() {
					r := New("test")
					r.Info("root-info")
					r.Error("root-error")
					r1 := r.Context("child1")
					r1.Error("child1-error")
					r2 := r.Context("child2")
					r2.Info("child2-info")
					Expect(r.String()).To(Equal("test\n\t\t[ info ] root-info\n\t\t[ error ] root-error\n\tchild1\n\t\t\t[ error ] child1-error\n\tchild2\n\t\t\t[ info ] child2-info\n"))
				})
			})
		})
	})
})
