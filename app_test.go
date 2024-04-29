package main_test

import (
	"io"
	"net/http"
	"strings"
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
		It("run, terminate, then monitor", func() {
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

			respM, errM := http.Get("http://localhost:9010/api/v1/monitor")
			if errM != nil {
				Expect(errM).To(BeNil())
			}
			defer respM.Body.Close()
			bodyM, _ := io.ReadAll(respM.Body)

			Expect(string(bodyT)).To(Equal("{\"termination\":started}"))
			Expect(string(bodyM)).To(Equal("ongoing hashes=[]\n"))
		})

		It("post job, monitor, terminate then probe", func() {
			api := &App{}
			portNum := 9020
			api.Initialize(&portNum)
			go api.Run()

			Expect(api).To(Not(BeNil()))
			Expect(api.Port).To(Equal(":9020"))

			respP, errP := http.Post("http://localhost:9020/api/v1/job/high?inputpath=%2Ftmp%2Fsrc&outputpath=%2Ftmp%2Fdst&w=1920&h=1080", "application/json; charset=utf-8", nil)
			if errP != nil {
				Expect(errP).To(BeNil())
			}
			defer respP.Body.Close()
			bodyP, _ := io.ReadAll(respP.Body)
			time.Sleep(1 * time.Second)

			Expect(string(bodyP)).To(ContainSubstring("{\"id\":"))
			time.Sleep(1 * time.Second)

			respM, errM := http.Get("http://localhost:9020/api/v1/monitor")
			if errM != nil {
				Expect(errM).To(BeNil())
			}
			defer respM.Body.Close()
			bodyM, _ := io.ReadAll(respM.Body)
			Expect(string(bodyM)).To(ContainSubstring("ongoing hashes=["))

			hashBracket := strings.Split(string(bodyM), "ongoing hashes=[")[1]
			hash := strings.Split(hashBracket, "]")[0]

			respT, errT := http.Get("http://localhost:9020/api/v1/terminate")
			if errT != nil {
				Expect(errT).To(BeNil())
			}
			defer respT.Body.Close()
			bodyT, _ := io.ReadAll(respT.Body)

			Expect(string(bodyT)).To(Equal("{\"termination\":started}"))

			respG, errG := http.Get("http://localhost:9020/api/v1/probe/" + hash)
			if errG != nil {
				Expect(errG).To(BeNil())
			}
			defer respG.Body.Close()
			bodyG, _ := io.ReadAll(respG.Body)
			Expect(string(bodyG)).To(ContainSubstring("{\"done\":}"))
		})
	})
})
