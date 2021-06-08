package controllers_test

import (
	"github.com/kamva/mgm/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"go.mongodb.org/mongo-driver/bson"
	"net/http"

	"github.com/drift-org/backend/controllers"
	"github.com/drift-org/backend/helpers"
	"github.com/drift-org/backend/models"
)

var _ = Describe("AuthController", func() {

	var (
		authController controllers.AuthController = controllers.NewAuthController()
	)

	Describe("Register", func() {
		It("Basic", helpers.TestWithMongo("AuthController-Register", func() {

			// Create a sample request to register.
			context := helpers.CreateTestContext("", `{
        "username": "TestUsername",
        "name": "TestName",
        "age": 10,
        "emailAddress": "Test@Test.com",
        "password": "TestPassword"
      }`)

			// Test the first user, and that the response is OK.
			authController.Register(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusOK))

			// Test that the user was written to mongo.
			user := &models.User{}
			_ = mgm.Coll(user).First(bson.M{}, user)
			Expect(user).NotTo(BeNil())

			// Test that the password is encrypted.
			Expect(user.Password).NotTo(Equal("TestPassword"))

			// Attempt creating a user with the same username.
			duplicateContext1 := helpers.CreateTestContext("", `{
        "username": "TestUsername",
        "name": "differentName",
        "age": 2,
        "emailAddress": "differentEmail@gmail.com",
        "password": "TestPassword"
      }`)
			authController.Register(duplicateContext1)
			Expect(duplicateContext1.Writer.Status()).To(Equal(http.StatusBadRequest))

			// Attempt creating a user with the same email.
			duplicateContext2 := helpers.CreateTestContext("", `{
        "username": "differentUsername",
        "name": "differentName",
        "age": 2,
        "emailAddress": "Test@Test.com",,
        "password": "TestPassword"
      }`)
			authController.Register(duplicateContext2)
			Expect(duplicateContext2.Writer.Status()).To(Equal(http.StatusBadRequest))
		}))
	})
})
