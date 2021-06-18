package controllers_test

import (
	"github.com/kamva/mgm/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/drift-org/backend/controllers"
	"github.com/drift-org/backend/helpers"
	"github.com/drift-org/backend/models"
)

var _ = Describe("GroupController", func() {

	var (
		groupController controllers.GroupController = controllers.NewGroupController()
	)

	//----------------------------------------------------------------------------------------

	Describe("Create", func() {
		It("Basic", helpers.TestWithMongo("GroupController-Create-Basic", func() {

			// Set up basic user accounts.
			helpers.CreateDummyUser("username1")
			helpers.CreateDummyUser("username2")
			helpers.CreateDummyUser("username3")

			// Create the first request to Create the group.
			context := helpers.CreateTestContext("", `{
				"usernames": ["username1", "username2", "username3"]
			}`)

			// Login as username1
			helpers.LoginAuthUsername(context, "username1")

			// Test the first group, and that the response is OK.
			groupController.Create(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusOK))

			// Test that the group was written to mongo.
			group := &models.Group{}
			_ = mgm.Coll(group).First(bson.M{}, group)
			Expect(group).NotTo(BeNil())

			// Attempt creating a group with the same usernames, but in a different order.
			duplicateContext := helpers.CreateTestContext("", `{
				"usernames": ["username3", "username1", "username2"]
			}`)
			helpers.LoginAuthUsername(duplicateContext, "username1")
			groupController.Create(duplicateContext)
			Expect(duplicateContext.Writer.Status()).To(Equal(http.StatusBadRequest))
		}))

		It("Test Unique Validation", helpers.TestWithMongo("GroupController-Create-Unique", func() {

			// Set up basic user accounts.
			helpers.CreateDummyUser("username1")
			helpers.CreateDummyUser("username3")

			// Create the request to Create the group, with a duplicate username
			context := helpers.CreateTestContext("", `{
				"usernames": ["username1", "username3", "username3"]
			}`)

			// Login as username1
			helpers.LoginAuthUsername(context, "username1")

			// Test the group, and that the response is StatusBadRequest.
			groupController.Create(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusBadRequest))
		}))

		It("Test Bearer Inclusion", helpers.TestWithMongo("GroupController-Create-Bearer", func() {

			// Set up basic user accounts.
			helpers.CreateDummyUser("username1")
			helpers.CreateDummyUser("username2")
			helpers.CreateDummyUser("username3")

			// Create the request to Create the group, but leave out username1.
			context := helpers.CreateTestContext("", `{
				"usernames": ["username2", "username3"]
			}`)

			// Login as username1
			helpers.LoginAuthUsername(context, "username1")

			// Test the group, and that the response is StatusBadRequest.
			groupController.Create(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusBadRequest))
		}))

		It("Test Non Existing User", helpers.TestWithMongo("GroupController-Create-Existing", func() {

			// Set up basic user accounts, but leave out username3.
			helpers.CreateDummyUser("username1")
			helpers.CreateDummyUser("username2")

			// Create the request to Create the group.
			context := helpers.CreateTestContext("", `{
				"usernames": ["username1", "username2", "username3"]
			}`)
			helpers.LoginAuthUsername(context, "username1")

			// Test the group, and that the response is StatusBadRequest.
			groupController.Create(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusBadRequest))
		}))
	})
})
