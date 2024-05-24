package aggregate

import (
	"math"
	"sync"
)

type Hex struct {
	HexId int    `bson:"hex_id,omitempty"`
	Hex   string `bson:"hex,omitempty"`
}

type Segment struct {
	PageId      string `json:"page_id,omitempty"`
	HashSegment string `json:"hash_segment,omitempty"`
}

type Candidates struct {
	Mu          sync.Mutex `json:"mu"`
	QueryResult []Segment  `json:"results,omitempty"`
}

type Result struct {
	Hash     string
	Distance int
}

const (
	BucketCount            = 8
	DefaultRadius          = 5
	DefaultSubstringRadius = 3
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SplitLength(s string, buckets int) int {
	length := math.Ceil(float64(len(s)) / float64(buckets))
	return int(length)
}

func SplitSegments(s string, segment int) []string {
	i := 0
	var out []string
	for i < len(s) {
		out = append(out, s[i:Min(len(s), i+segment)])
		i += segment
	}

	return out
}

func Compare(s1, s2 string) int {
	count := 0
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			count += 1
		}
	}
	return count
}

// GroupNeighbourSegments combines segments of the same page from a list of segments
func GroupNeighbourSegments(candidates *[]Segment) *map[string][]string {
	groups := make(map[string][]string)
	for _, i := range *candidates {
		groups[i.PageId] = append(groups[i.PageId], i.HashSegment)
	}

	return &groups
}
