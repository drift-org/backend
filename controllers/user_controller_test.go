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
		authController controllers.AuthController = controllers.NewAuthController()
	)

	Describe("Get", func() {
		It("Basic", helpers.TestWithMongo("UserController-Get", func() {

			// Register a sample user to test on.
			authContext := helpers.CreateTestContext("", `{
				"username": "TestUsername",
				"name": "TestName",
				"age": 10,
				"emailAddress": "Test@Test.com",
				"password": "TestPassword"
			}`)

			authController.Register(authContext)
			sampleUser := &models.User{}
			_ = mgm.Coll(sampleUser).First(bson.M{}, sampleUser)
			var sampleUserID string = sampleUser.DefaultModel.IDField.ID.Hex()

			// Create a sample request with the sample user's ID.
			var input string = `{
				"id": "` + sampleUserID + `",
				"specificity": 1
			}`
			context := helpers.CreateTestContext("", input)

			// Test that the response is OK.
			userController.Get(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusOK))

			// TODO: Test that response is same as sample user
		}))

		It("Test Specificity Validation", helpers.TestWithMongo("UserController-Get", func() {

			// Register a sample user to test on.
			authContext := helpers.CreateTestContext("", `{
				"username": "TestUsername",
				"name": "TestName",
				"age": 10,
				"emailAddress": "Test@Test.com",
				"password": "TestPassword"
			}`)

			authController.Register(authContext)
			sampleUser := &models.User{}
			_ = mgm.Coll(sampleUser).First(bson.M{}, sampleUser)
			var sampleUserID string = sampleUser.DefaultModel.IDField.ID.Hex()

			// Create a sample request with the sample user's ID.
			var input string = `{
				"id": "` + sampleUserID + `",
				"specificity": -1
			}`
			context := helpers.CreateTestContext("", input)

			// Test that the response is OK.
			userController.Get(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusBadRequest))

			// TODO: Test that response is same as sample user
		}))
	})
})
