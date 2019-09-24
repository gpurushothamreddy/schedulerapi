package acceptance

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("Post Schedule", func() {
	var (
		err error
	)

	Describe("with valid request", func() {
		var response *http.Response
		JustBeforeEach(
			func() {
				req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:%v/%s", hostName, serverPort, "schedule"), bytes.NewBuffer([]byte(`{
					"command": "ls -lah",
					"scheduleDateTime": "9999-09-22 15:27:30 PST"
					}`)))
				response, err = client.Do(req)
				Expect(err).NotTo(HaveOccurred())

			})
		It("returns 200 ok status code along with success response", func() {
			Expect(response.Header.Get("Content-Type")).To(Equal(""))
			Expect(response.StatusCode).To(Equal(http.StatusCreated))
		})
		It("should match contract", func() {
			body := response.Body
			defer body.Close()
			responseInBytes, err := ioutil.ReadAll(body)
			Expect(err).NotTo(HaveOccurred())
			subject := string(responseInBytes)
			Expect(subject).To(Equal(""))
		})
	})

})
