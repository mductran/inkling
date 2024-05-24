package search

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"math"

	"github.com/schollz/progressbar/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

type Hex struct {
	hash string
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func insert(client *mongo.Client) {
	collection := client.Database("fnn").Collection("hash")
	upper := math.Pow10(8)
	bar := progressbar.Default(int64(upper))

	for i := 0; i < int(upper); i++ {
		h, _ := randomHex(64)
		collection.InsertOne(context.TODO(), Hex{
			hash: h,
		})
		// fmt.Println(result, err)
		bar.Add(1)
	}
}
