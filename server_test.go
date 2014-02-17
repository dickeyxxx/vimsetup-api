package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	Describe("notFound", func() {
		It("returns the not found code", func() {
			code, _ := notFound()
			Ω(code).Should(Equal(404))
		})
	})
})
