package main

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"Schwarz--Internship--2026/services/common/rabbitmq"
	"Schwarz--Internship--2026/services/reservation-base/main/proto"
)

func (service *ReservationServiceImpl) CreateReservation(ctx context.Context, req *proto.CreateReservationRequest) (*proto.CreateReservationResponse, error) {
	if req == nil || req.Reservation == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request or missing reservation data")
	}

	res := req.GetReservation()

	if res.GetPropertyId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "property_id is missing")
	}
	if res.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is missing")
	}
	if res.GetOwnerId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "owner_id is missing")
	}
	if res.GetCheckInDate() == "" {
		return nil, status.Error(codes.InvalidArgument, "check_in_date is missing")
	}
	if res.GetCheckOutDate() == "" {
		return nil, status.Error(codes.InvalidArgument, "check_out_date is missing")
	}

	res.Status = proto.ReservationStatus_RESERVATION_STATUS_PENDING

	id, err := InsertReservation(ctx, service.DB, res)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert reservation into database: %v", err)
	}

	res.Id = id
	now := timestamppb.Now()
	res.CreatedAt = now
	res.UpdatedAt = now

	if service.EmailProd != nil {
		service.sendReservationNotification(ctx, res)
	}

	return &proto.CreateReservationResponse{
		Reservation: res,
	}, nil
}

func (service *ReservationServiceImpl) sendReservationNotification(ctx context.Context, res *proto.Reservation) {
	emailResp, err := service.UserService.GetEmail(ctx, &proto.GetEmailRequest{Id: res.GetOwnerId()})
	if err != nil {
		slog.ErrorContext(ctx, "failed to get owner email for notification", "reservation_id", res.GetId(), "error", err)
		return
	}
	userResp, err := service.UserService.GetUserProfile(ctx, &proto.GetUserProfileRequest{Id: res.GetUserId()})
	if err != nil {
		slog.ErrorContext(ctx, "failed to get user profile for notification", "reservation_id", res.GetId(), "error", err)
		return
	}
	propResp, err := service.PropertyService.GetProperty(ctx, &proto.GetPropertyRequest{Id: res.GetPropertyId()})
	if err != nil {
		slog.ErrorContext(ctx, "failed to get property details for notification", "reservation_id", res.GetId(), "error", err)
		return
	}
	user := userResp.GetUser()
	prop := propResp.GetProperty()

	email := rabbitmq.EmailMessage{
		To:      emailResp.GetEmail(),
		Subject: fmt.Sprintf("Reservation created for %s", prop.GetName()),
		Body: fmt.Sprintf(
			"Created by %s %s\nFrom %s Until %s\n",
			user.GetFirstName(), user.GetLastName(), res.GetCheckInDate(), res.GetCheckOutDate(),
		),
	}

	if err := service.EmailProd.Publish(ctx, email); err != nil {
		slog.ErrorContext(ctx, "failed to publish email message", "reservation_id", res.GetId(), "error", err)
	}
}
