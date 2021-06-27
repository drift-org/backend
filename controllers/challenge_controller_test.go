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

var _ = Describe("ChallengeController", func() {

	var (
		challengeController controllers.ChallengeController = controllers.NewChallengeController()
	)

	Describe("Register", func() {
		It("Basic", helpers.TestWithMongo("ChallengeController-Create", func() {

			// Create a sample request for a challenge to be created.
			context := helpers.CreateTestContext("", `{
				"latitude": 10.01,
				"longitude": 10.01,
				"address": "1234 Main Street",
				"taskName": "Task 1",
				"description": "Description 1",
				"points": 10
			}`)

			// Test the first challenge, and that the response is OK.
			challengeController.Create(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusOK))

			// Test that the challenge was written to mongo.
			challenge := &models.Challenge{}
			_ = mgm.Coll(challenge).First(bson.M{}, challenge)
			Expect(challenge).NotTo(BeNil())

			// Attempt creating a challenge with the same address & taskName.
			duplicateContext1 := helpers.CreateTestContext("", `{
				"latitude": 10.01,
				"longitude": 10.01,
				"address": "1234 Main Street",
				"taskName": "Task 1",
				"description": "Description 2",
				"points": 20
			}`)
			challengeController.Create(duplicateContext1)
			Expect(duplicateContext1.Writer.Status()).To(Equal(http.StatusBadRequest))

			// Attempt creating a user with just the same address.
			duplicateContext2 := helpers.CreateTestContext("", `{
				"latitude": 10.01,
				"longitude": 10.01,
				"address": "1234 Main Street",
				"taskName": "Task 2",
				"description": "Description 2",
				"points": 20
			}`)
			challengeController.Create(duplicateContext2)
			Expect(duplicateContext2.Writer.Status()).To(Equal(http.StatusOK))

			// Attempt creating a user with just the same task name.
			duplicateContext3 := helpers.CreateTestContext("", `{
				"latitude": 20.02,
				"longitude": 20.02,
				"address": "9999 Oak Street",
				"taskName": "Task 1",
				"description": "Description 2",
				"points": 20
			}`)
			challengeController.Create(duplicateContext3)
			Expect(duplicateContext3.Writer.Status()).To(Equal(http.StatusOK))

		}))
	})
})
