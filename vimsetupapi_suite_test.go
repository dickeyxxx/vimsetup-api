package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestVimsetupapi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vimsetupapi Suite")
}
