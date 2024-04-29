package main_test

import (
	"io"
	"net/http"
	"time"

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

	Context("run", func() {
		It("start and terminate", func() {
			api := &App{}
			portNum := 9010
			api.Initialize(&portNum)
			go api.Run()

			Expect(api).To(Not(BeNil()))
			Expect(api.Port).To(Equal(":9010"))

			respT, errT := http.Get("http://localhost:9010/api/v1/terminate")
			if errT != nil {
				Expect(errT).To(BeNil())
			}
			defer respT.Body.Close()
			bodyT, _ := io.ReadAll(respT.Body)

			time.Sleep(1 * time.Second)

			respM, errM := http.Get("http://localhost:8081")
			if errM != nil {
				Expect(errM).To(BeNil())
			}
			defer respM.Body.Close()
			bodyM, _ := io.ReadAll(respM.Body)

			Expect(string(bodyT)).To(Equal("{\"termination\":started}"))
			Expect(string(bodyM)).To(Equal("ongoing hashes=[]\n"))
		})
	})
})
