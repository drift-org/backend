package controllers_test

import (
	"net/http"

	"github.com/kamva/mgm/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/drift-org/backend/controllers"
	"github.com/drift-org/backend/helpers"
	"github.com/drift-org/backend/models"
)

var _ = Describe("UserController", func() {

	var (
		userController controllers.UserController = controllers.NewUserController()
	)

	Describe("Get", func() {
		It("Test Specificity 0", helpers.TestWithMongo("UserController-Get-0", func() {

			// Register and retrieve a sample user to test on.
			helpers.CreateDummyUser("username1")
			user := &models.User{}
			_ = mgm.Coll(user).First(bson.M{"username": "username1"}, user)

			// Create a sample request with the sample user ID and valid specificity.
			var input string = `{
				"id": "` + user.DefaultModel.ID.Hex() + `",
				"specificity": 0
			}`
			context := helpers.CreateTestContext("", input)

			// Test that the response is OK.
			helpers.LoginAuthUsername(context, "username1")
			userController.Get(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusOK))

			// TODO: Test that response is same as sample user

			// Test that request with no specificity has same result.
			var input2 string = `{
				"id": "` + user.DefaultModel.ID.Hex() + `"
			}`
			context2 := helpers.CreateTestContext("", input2)

			// Test that the response is OK.
			helpers.LoginAuthUsername(context, "username1")
			userController.Get(context2)
			Expect(context2.Writer.Status()).To(Equal(http.StatusOK))

			// TODO: Test that response is same as sample user
		}))

		It("Test Specificity 1", helpers.TestWithMongo("UserController-Get-1", func() {

			// Register and retrieve a sample user to test on.
			helpers.CreateDummyUser("username1")
			user := &models.User{}
			_ = mgm.Coll(user).First(bson.M{"username": "username1"}, user)

			// Create a sample request with the sample user ID and valid specificity.
			var input string = `{
				"id": "` + user.DefaultModel.ID.Hex() + `",
				"specificity": 1
			}`
			context := helpers.CreateTestContext("", input)

			// Test that the response is OK.
			helpers.LoginAuthUsername(context, "username1")
			userController.Get(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusOK))

			// TODO: Test that response is same as sample user + friend objects
		}))

		It("Test Specificity Validation", helpers.TestWithMongo("UserController-Get-Validation", func() {

			// Register and retrieve a sample user to test on.
			helpers.CreateDummyUser("username1")
			user := &models.User{}
			_ = mgm.Coll(user).First(bson.M{"username": "username1"}, user)

			// Create a sample request with the sample user ID and invalid specificity.
			var input string = `{
				"id": "` + user.DefaultModel.ID.Hex() + `",
				"specificity": -1
			}`
			context := helpers.CreateTestContext("", input)

			// Test that the response is an error.
			helpers.LoginAuthUsername(context, "username1")
			userController.Get(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusBadRequest))
		}))

		It("Test Invalid ID", helpers.TestWithMongo("UserController-Get-ID", func() {

			// Register and retrieve a sample user to test on.
			helpers.CreateDummyUser("username1")
			user := &models.User{}
			_ = mgm.Coll(user).First(bson.M{"username": "username1"}, user)

			// Create a sample request with the sample user's ID.
			context := helpers.CreateTestContext("", `{
				"id": ""
			}`)

			// Test that the response is an error.
			helpers.LoginAuthUsername(context, "username1")
			userController.Get(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusBadRequest))
		}))
	})
})
