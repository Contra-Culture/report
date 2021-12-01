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
					Expect(r.String()).To(Equal("root: test\n"))
				})
			})
			Describe(".Context()", func() {
				It("creates child context", func() {
					r := New("test")
					child := r.Context("child")
					Expect(child).NotTo(BeNil())
					Expect(child.String()).To(Equal("child\n"))
					Expect(r.String()).To(Equal("root: test\n\tchild\n"))
				})
			})
			Describe("records", func() {
				Describe(".Error()", func() {
					It("adds error record", func() {
						r := New("test")
						r.Error("some error")
						Expect(r.String()).To(Equal("root: test\n\t\t[ error ] some error\n"))
					})
				})
				Describe(".Warn()", func() {
					It("adds warn record", func() {
						r := New("test")
						r.Warn("some warn")
						Expect(r.String()).To(Equal("root: test\n\t\t[ warning ] some warn\n"))
					})
				})
				Describe(".Info()", func() {
					It("adds info record", func() {
						r := New("test")
						r.Info("some info")
						Expect(r.String()).To(Equal("root: test\n\t\t[ info ] some info\n"))
					})
				})
				Describe(".Deprecation()", func() {
					It("adds deprecation record", func() {
						r := New("test")
						r.Deprecation("some deprecation")
						Expect(r.String()).To(Equal("root: test\n\t\t[ deprecated ] some deprecation\n"))
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
					Expect(r.String()).To(Equal("root: test\n\t\t[ info ] root-info\n\t\t[ error ] root-error\n\tchild1\n\t\t\t[ error ] child1-error\n\tchild2\n\t\t\t[ info ] child2-info\n"))
				})
			})
		})
	})
})
