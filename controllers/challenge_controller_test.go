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
		It("Basic & Unique Task Name/Location", helpers.TestWithMongo("ChallengeController-Create-Unique", func() {

			// Create a sample request for a challenge to be created.
			context := helpers.CreateTestContext("", `{
				"latitude": 34.0701449,
				"longitude": -118.4422936,
				"address": "1234 Main Street",
				"taskName": "Task 1",
				"description": "Description",
				"points": 10
			}`)

			// Test the first challenge, and that the response is OK.
			challengeController.Create(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusOK))

			// Test that the challenge was written to mongo.
			challenge := &models.Challenge{}
			_ = mgm.Coll(challenge).First(bson.M{}, challenge)
			Expect(challenge).NotTo(BeNil())

			// Attempt creating a user with just the same location.
			context2 := helpers.CreateTestContext("", `{
				"latitude": 34.0701449,
				"longitude": -118.4422936,
				"address": "1234 Main Street",
				"taskName": "Task 2",
				"description": "Description",
				"points": 20
			}`)
			challengeController.Create(context2)
			Expect(context2.Writer.Status()).To(Equal(http.StatusOK))

			// Test that the challenge was written to mongo.
			_ = mgm.Coll(challenge).First(bson.M{}, challenge)
			Expect(challenge).NotTo(BeNil())

			// Attempt creating a user with just the same task name.
			context3 := helpers.CreateTestContext("", `{
				"latitude": 20.02,
				"longitude": 20.02,
				"address": "1234 Main Street",
				"taskName": "Task 1",
				"description": "Description",
				"points": 20
			}`)
			challengeController.Create(context3)
			Expect(context3.Writer.Status()).To(Equal(http.StatusOK))

			// Test that the challenge was written to mongo.
			_ = mgm.Coll(challenge).First(bson.M{}, challenge)
			Expect(challenge).NotTo(BeNil())

			// Attempt creating a challenge with the same location & taskName.
			context4 := helpers.CreateTestContext("", `{
				"latitude": 34.0701449,
				"longitude": -118.4422936,
				"address": "1234 Main Street",
				"taskName": "Task 1",
				"description": "Description",
				"points": 20
			}`)
			challengeController.Create(context4)
			Expect(context4.Writer.Status()).To(Equal(http.StatusBadRequest))

			// Attempt creating challenge with no location/address info, but same task name
			context5 := helpers.CreateTestContext("", `{
				"taskName": "Task 1",
				"description": "Description",
				"points": 20
			}`)
			challengeController.Create(context5)
			Expect(context5.Writer.Status()).To(Equal(http.StatusBadRequest))

		}))

		It("Test No Location Info", helpers.TestWithMongo("ChallengeController-Create-No-Location", func() {

			// Attempt creating a challenge with only the minimum required fields.
			context := helpers.CreateTestContext("", `{
				"taskName": "Task 1",
				"description": "Description",
				"points": 20
			}`)
			challengeController.Create(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusOK))

			// Test that the challenge was written to mongo.
			challenge := &models.Challenge{}
			_ = mgm.Coll(challenge).First(bson.M{}, challenge)
			Expect(challenge).NotTo(BeNil())

			// Attempt creating a challenge with the same taskName.
			duplicateContext := helpers.CreateTestContext("", `{
				"taskName": "Task 1",
				"description": "Description",
				"points": 20
			}`)
			challengeController.Create(duplicateContext)
			Expect(duplicateContext.Writer.Status()).To(Equal(http.StatusBadRequest))

		}))

		It("Test Invalid Arguments", helpers.TestWithMongo("ChallengeController-Create-Invalid-Args", func() {

			// Create a sample request for a challenge to be created, with an invalid (negative) point value.
			context := helpers.CreateTestContext("", `{
				"latitude": 10.01,
				"longitude": 10.01,
				"address": "1234 Main Street",
				"taskName": "Task 1",
				"description": "Description",
				"points": -1
			}`)

			// Test the challenge, and that the response is bad.
			challengeController.Create(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusBadRequest))

			// Create a sample request for a challenge to be created, with an invalid (decimal) point value.
			context2 := helpers.CreateTestContext("", `{
				"latitude": 10.01,
				"longitude": 10.01,
				"address": "1234 Main Street",
				"taskName": "Task 1",
				"description": "Description",
				"points": 10.5
			}`)

			// Test the challenge, and that the response is bad.
			challengeController.Create(context2)
			Expect(context2.Writer.Status()).To(Equal(http.StatusBadRequest))

			// Attempt creating a challenge with invalid lat/long coordinates (out of bounds).
			context3 := helpers.CreateTestContext("", `{
				"latitude": 500,
				"longitude": -500,
				"address": "1234 Main Street",
				"taskName": "Task 1",
				"description": "Description",
				"points": 20
			}`)

			// Test the challenge, and that the response is bad.
			challengeController.Create(context3)
			Expect(context3.Writer.Status()).To(Equal(http.StatusBadRequest))

		}))

		It("Autopopulates Arguments", helpers.TestWithMongo("ChallengeController-Create-Autopopulate", func() {

			// Create a sample request for a challenge to be created, without lat/long.
			context := helpers.CreateTestContext("", `{
				"address": "10740 Dickson Ct, Los Angeles, CA 90095",
				"taskName": "Task 1",
				"description": "Description",
				"points": 10
			}`)

			// Test the first challenge, and that the response is OK.
			challengeController.Create(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusOK))

			// Test that the challenge was written to mongo.
			challenge := &models.Challenge{}
			_ = mgm.Coll(challenge).First(bson.M{}, challenge)
			Expect(challenge).NotTo(BeNil())

			// Test that lat/long have been created.
			Expect(challenge.Location).NotTo(BeNil())

			// Attempt to create a challenge with no lat/long but same address.
			// Should autopopulate same coordinates & fail.
			context2 := helpers.CreateTestContext("", `{
				"address": "10740 Dickson Ct, Los Angeles, CA 90095",
				"taskName": "Task 1",
				"description": "Description",
				"points": 20
			}`)
			challengeController.Create(context2)
			Expect(context2.Writer.Status()).To(Equal(http.StatusBadRequest))

		}))

	})

})
