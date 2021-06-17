package helpers

import (
	"context"
	"time"

	"github.com/drift-org/backend/models"
)

/*
Setup the indexes of all the collections

As we create indexes for models, make sure to add them to this function
*/
func SetupIndexes() {
	ctx, canc := context.WithTimeout(context.Background(), 5*time.Second)
	defer canc()

	models.ChallengeIndex(ctx)
}
