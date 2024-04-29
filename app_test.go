package main_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/RafikMikhael/gct"
)

var _ = Describe("app.go", func() {
	Context("initialize", func() {
		It("positive", func() {
			api := &App{}
			portNum := 9000
			api.Initialize(&portNum)

			Expect(api).To(Not(BeNil()))
			Expect(api.Port).To(Equal(":9000"))
		})
	})
})
