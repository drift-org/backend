/*
Helper functions that are used for testing purposes in test files.
*/

package helpers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/http/httptest"
	"os"
)

/*
Creates a gin Context object for testing purposes.

Parameters:
- query are the query parameters of the request
- body is the verbatim json body of the request

Example:
	CreateTestContext('?foo=bar&page=10', `{
		"hello": "foo"
	}`)
*/
func CreateTestContext(query string, body string) *gin.Context {
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Add the http request to the context. Technically, not all our requests are POST requests, but
	// it doesn't actually matter for testing since we don't reference the context type in any of our
	// controllers. The only thing that the context should be used for is extracting the body and/or
	// query parameters.
	context.Request, _ = http.NewRequest("POST", "/"+query, bytes.NewBufferString(body))

	// Add the json header.
	context.Request.Header.Add("Content-Type", "application/json")
	return context
}

/*
Defines a ginkgo test that tests a function that uses mongodb. This helper
function automatically adds a connection to your local mondgb host.

Parameters:
- databaseName is the name of the database, which should be **unique**.
- tester is a function that runs the test.

Example:
It("Test Description", helpers.TestWithMongo("controller-method", func() {
	// test a controller that uses mongo/mgm
	...
}))
*/
func TestWithMongo(databaseName string, tester func()) func() {
	return func() {
		// First connect to the database.
		database := ConnectTestDB(databaseName)

		// Immediately reset the database. Even though we reset the database after
		// the test runs, we still reset it before as well. This is because if a test
		// fails halfway through, the reset afterwards is never hit.
		ResetTestDB(database)

		// Run the test and reset the database.
		tester()
		ResetTestDB(database)
	}
}

/*
Connects to a passed-in database on your LOCAL mongodb host, for testing purposes.

To allow for concurrent testing specs, this function should be called
**separately** on each test (that requires mongodb), with a **different** database name.
For more information, see https://github.com/drift-org/backend/wiki/Testing

Parameters:
- databaseName is the name of the database.
*/
func ConnectTestDB(databaseName string) *mongo.Database {

	// Ginkgo runs concurrent tests in different processes, so load the .env file each time.
	err := godotenv.Load(os.ExpandEnv("../.env"))
	Expect(err).NotTo(HaveOccurred())

	// Parse the local database url.
	databaseURL := os.Getenv("TESTING_MONGO_URL")
	Expect(databaseURL).NotTo(BeEmpty(), "TESTING_MONGO_URL is missing from .env")

	// Make sure that the word "localhost" is inside the databaseURL. We want to double check that this
	// is a locally ran mongo client and make sure we aren't writing to our production database.
	Expect(databaseURL).To(ContainSubstring("localhost"), "TESTING_MONGO_URL must contain 'localhost'")

	// Create the client options. Since this test is in a different process, set the default config.
	clientOptions := options.Client().ApplyURI(databaseURL)
	err = mgm.SetDefaultConfig(nil, databaseName, clientOptions)
	Expect(err).NotTo(HaveOccurred())

	// Connect to the client.
	client, err := mongo.Connect(context.TODO(), clientOptions)
	Expect(err).NotTo(HaveOccurred())
	Expect(client).NotTo(BeNil())

	return client.Database(databaseName)
}

/*
Resets to a passed-in database on your LOCAL mongodb host, for testing purposes.

This function should be called after each test that used mongodb has finished.
We do this to make sure our tests are completely deterministic with no state
prior to testing.
*/
func ResetTestDB(database *mongo.Database) {

	// Get a slice of all collections.
	collections, err := database.ListCollectionNames(mgm.Ctx(), bson.D{})
	Expect(err).NotTo(HaveOccurred())
	Expect(collections).NotTo(BeNil())

	// Drop each collection.
	for _, coll := range collections {
		err := database.Collection(coll).Drop(mgm.Ctx())
		Expect(err).NotTo(HaveOccurred())
	}
}

/*
Logs std output. The standard fmt.Println doesn't show up when testing with ginkgo.
For debugging purposes ONLY.
*/
func GinkgoLog(message string) {
	_, _ = fmt.Fprintf(GinkgoWriter, "[DEBUG]"+message)
}
