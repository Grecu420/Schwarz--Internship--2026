package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenPayload struct {
	ID         int64  `json:"id"`
	FilterHash string `json:"hash"`
}

type Identifiable interface {
	GetId() int64
}

func HashFilters(values ...any) string {
	h := fnv.New32a()
	for _, val := range values {
		fmt.Fprint(h, val)
	}
	return fmt.Sprintf("%x", h.Sum32())
}

func BuildNextPageToken(id int64, filterHash string) (string, error) {

	newToken := TokenPayload{
		ID:         id,
		FilterHash: filterHash,
	}
	tokenBytes, err := json.Marshal(newToken)
	if err != nil {
		return "", err
	}
	nextPageToken := base64.StdEncoding.EncodeToString(tokenBytes)
	return nextPageToken, nil
}

// Paginate abstracts token parsing, offset lookup, query execution, and response token generation.
func Paginate[T Identifiable](
	pageSize int64,
	nextPageToken string,
	filterHash string,
	fetchFunc func(offsetId int64, limit int64) ([]T, error),
) ([]T, string, error) {
	if pageSize <= 0 {
		return nil, "", status.Errorf(codes.InvalidArgument, "invalid page size: %v", pageSize)
	}

	// Parse token to obtain offset
	var offsetId int64
	if nextPageToken != "" {
		decoded, err := base64.StdEncoding.DecodeString(nextPageToken)
		if err != nil {
			return nil, "", status.Errorf(codes.InvalidArgument, "invalid base64 page token: %v", err)
		}

		var payload TokenPayload
		if err := json.Unmarshal(decoded, &payload); err != nil {
			return nil, "", status.Errorf(codes.InvalidArgument, "malformed page token: %v", err)
		}

		// Verify hash for alteration
		if payload.FilterHash != filterHash {
			return nil, "", status.Error(codes.InvalidArgument, "filter criteria changed during pagination")
		}
		offsetId = payload.ID
	}

	// Execute query to obtain
	items, err := fetchFunc(offsetId, pageSize)
	if err != nil {
		return nil, "", status.Errorf(codes.Internal, "failed db query: %v", err)
	}

	// Build next token if there are items remaining
	if len(items) == int(pageSize)+1 {
		lastItem := items[len(items)-1]
		newOffsetId := lastItem.GetId()
		nextPageToken, err = BuildNextPageToken(newOffsetId, filterHash)
		if err != nil {
			return nil, "", status.Errorf(codes.Internal, "failed nextPageToken creation: %v", err)
		}
		items = items[:pageSize]
	} else {
		nextPageToken = ""
	}

	return items, nextPageToken, nil
}
