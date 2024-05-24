package aggregate

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SearchBucket searches in one bucket for segment within edit distance equals to radius
// By default, there are 8 buckets labeled segments_1 through 8
// TODO: implement different bucket count for different substring size
func SearchBucket(bucket *mongo.Collection, query string, radius int, ctx context.Context) (*[]Segment, error) {
	cursor, err := bucket.Find(ctx, bson.D{})
	if err != nil {
		return &[]Segment{}, err
	}
	var results []Segment
	if err = cursor.All(ctx, &results); err != nil {
		return &[]Segment{}, err
	}

	var neighbours []Segment
	for _, line := range results {
		if Compare(line.HashSegment, query) <= radius {
			neighbours = append(neighbours, line)
		}
	}

	return &neighbours, nil
}

// WriteHash breaks down hash string into segments and write them to different buckets.
// Each bucket is a MongoDB collection.
func WriteHash(db *mongo.Database, pageId, phash string, buckets int, ctx context.Context) (*[]string, error) {
	segments := SplitSegments(phash, SplitLength(phash, buckets))
	var segmentIds []string
	for i, j := range segments {
		result, err := db.Collection(fmt.Sprintf("segment_%d", i+1)).InsertOne(ctx, Segment{PageId: pageId, HashSegment: j})
		if err != nil {
			return &[]string{}, err
		}
		segmentIds = append(segmentIds, result.InsertedID.(primitive.ObjectID).Hex())
	}

	return &segmentIds, nil
}

// Query searches for neighbour within certain edit radius.
// Returns the list of candidates ordered by edit distance in ascending order.
func Query(database *mongo.Database, qString string) (*map[string][]string, error) {
	segments := SplitSegments(qString, BucketCount)
	var candidates = Candidates{}
	var radius = DefaultRadius

	var limit = DefaultRadius - DefaultSubstringRadius*BucketCount
	wg := &sync.WaitGroup{}

	for i := 0; i < BucketCount; i++ {
		wg.Add(1)
		go func(i int, candidates *Candidates, wg *sync.WaitGroup) {
			defer wg.Done()
			candidates.Mu.Lock()

			bucket := database.Collection(fmt.Sprintf("segment_%d", i))
			if i >= limit {
				radius = DefaultRadius - 1
			}
			neighbour, err := SearchBucket(bucket, segments[i], radius, context.TODO())
			if err != nil {
				return
			}
			candidates.QueryResult = append(candidates.QueryResult, *neighbour...)

			candidates.Mu.Unlock()
		}(i, &candidates, wg)
	}
	wg.Wait()

	groups := GroupNeighbourSegments(&candidates.QueryResult)

	rwLock := sync.RWMutex{}
	for key, val := range *groups {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			rwLock.RLock()
			r := strings.Join(val, "")
			if len(val) != BucketCount || Compare(qString, r) > DefaultRadius {
				rwLock.Lock()
				delete(*groups, key)
				rwLock.Unlock()
			}
			rwLock.Unlock()
		}(wg)
	}
	wg.Wait()

	return groups, nil
}
