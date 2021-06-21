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

/*
Simple utility function that lets a dummy prize in the database, with the given parameters..
*/
func CreateDummyPrize(name string, desc string) {
	mgm.Coll(&models.Prize{}).Create(&models.Prize{
		Name:        name,
		Description: desc,
	})
}

var _ = Describe("PrizeController", func() {

	var (
		prizeController controllers.PrizeController = controllers.NewPrizeController()
	)

	//----------------------------------------------------------------------------------------

	Describe("Create", func() {
		It("Basic", helpers.TestWithMongo("GroupController-Create-Basic", func() {

			// Create the first request to Create the prize.
			context := helpers.CreateTestContext("", `{
				"name": "50% Discount at Sweetheart Cafe",
				"description": "XY78-10DK"
			}`)

			// Test the creation of the prize, and that the response is OK.
			prizeController.Create(context)
			Expect(context.Writer.Status()).To(Equal(http.StatusOK))

			// Test that the prize was written to mongo.
			prize := &models.Prize{}
			_ = mgm.Coll(prize).First(bson.M{}, prize)
			Expect(prize).NotTo(BeNil())

		}))

	})
})
