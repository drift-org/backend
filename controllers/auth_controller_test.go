package controllers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"strconv"

	"github.com/drift-org/backend/controllers"
	"github.com/drift-org/backend/helpers"
)

var _ = Describe("AuthController", func() {

	var (
		authController controllers.AuthController = controllers.NewAuthController()
	)

	Describe("Register", func() {
		Context("Basic", func() {
			It("First User", func() {
				c := helpers.CreateTestContext("", `{
        	"username": "TestUsername",
	        "name": "TestName",
			    "age": 10,
			    "emailAddress": "Test@Test.com",
			    "password": "TestPassword"
        }`)

				authController.Register(c)
				Expect(c.Writer.Status()).To(Equal(http.StatusOK))

				helpers.LogTest(strconv.Itoa(c.Writer.Status()))
			})
		})
	})

})
