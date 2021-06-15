package helpers

import (
	"context"
	"time"

	"github.com/drift-org/backend/models"
)

func SetupIndexes() {
	ctx, canc := context.WithTimeout(context.Background(), 5*time.Second)
	defer canc()

	models.ChallengeIndex(ctx)
}
