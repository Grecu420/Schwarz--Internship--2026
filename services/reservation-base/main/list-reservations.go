package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"Schwarz--Internship--2026/services/reservation-base/main/proto"
)

type TokenPayload struct {
	ID         int64  `json:"id"`
	FilterHash string `json:"hash"`
}

func HashReservationFilter(propertyIds []int64, userId int64, ownerId int64, resStatus proto.ReservationStatus) string {
	h := fnv.New32a()
	str := fmt.Sprintf("%v-%d-%d-%d", propertyIds, userId, ownerId, resStatus)
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum32())
}

func getPayload(nextPageToken string) (*TokenPayload, error) {
	decoded, err := base64.StdEncoding.DecodeString(nextPageToken)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 page token: %w", err)
	}

	var payload TokenPayload
	if err := json.Unmarshal(decoded, &payload); err != nil {
		return nil, fmt.Errorf("malformed page token: %w", err)
	}
	return &payload, nil
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

func (service *ReservationServiceImpl) ListReservations(ctx context.Context, req *proto.ListReservationsRequest) (*proto.ListReservationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request can't be nil")
	}

	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page size: %v", pageSize)
	}

	filters := req.GetFilters()
	var propertyIds []int64
	var userId int64
	var ownerId int64
	var resStatus proto.ReservationStatus

	for _, f := range filters {
		switch v := f.GetFilter().(type) {
		case *proto.ListReservationsFiltersOneOf_PropertyId:
			propertyIds = append(propertyIds, v.PropertyId.GetValue())
		case *proto.ListReservationsFiltersOneOf_UserId:
			if userId != 0 {
				return nil, status.Errorf(codes.InvalidArgument, "duplicate user_id filter")
			}
			userId = v.UserId.GetValue()
		case *proto.ListReservationsFiltersOneOf_OwnerId:
			if ownerId != 0 {
				return nil, status.Errorf(codes.InvalidArgument, "duplicate owner_id filter")
			}
			ownerId = v.OwnerId.GetValue()
		case *proto.ListReservationsFiltersOneOf_Status:
			if resStatus != proto.ReservationStatus_RESERVATION_STATUS_UNSPECIFIED {
				return nil, status.Errorf(codes.InvalidArgument, "duplicate status filter")
			}
			resStatus = v.Status.GetValue()
		}
	}

	filterHash := HashReservationFilter(propertyIds, userId, ownerId, resStatus)

	var offsetId int64
	if req.GetNextPageToken() != "" {
		payload, err := getPayload(req.GetNextPageToken())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "next page token failure: %v", err)
		}
		if payload.FilterHash != filterHash {
			return nil, status.Error(codes.InvalidArgument, "filter criteria changed during pagination")
		}
		offsetId = payload.ID
	}

	reservations, err := SelectReservationListInDB(ctx, service.DB, offsetId, pageSize, propertyIds, userId, ownerId, resStatus)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed db query: %v", err)
	}

	var nextPageToken = ""
	
	if len(reservations) == int(pageSize)+1 {
		last := reservations[pageSize-1]
		newOffsetId := last.GetId()

		var err error
		nextPageToken, err = BuildNextPageToken(newOffsetId, filterHash)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed nextPageToken creation: %v", err)
		}

		reservations = reservations[:pageSize]
	}

	return &proto.ListReservationsResponse{
		NextPageToken: nextPageToken,
		Reservations:  reservations,
	}, nil
}