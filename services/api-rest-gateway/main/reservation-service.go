package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"context"
)

func (service GatewayServiceImpl) CreateReservation(ctx context.Context, req *proto.CreateReservationRequest) (*proto.CreateReservationResponse, error) {
	return service.resService.CreateReservation(ctx, req)
}

func (service GatewayServiceImpl) GetReservation(ctx context.Context, req *proto.GetReservationRequest) (*proto.GetReservationResponse, error) {
	return service.resService.GetReservation(ctx, req)
}

func (service GatewayServiceImpl) ListReservations(ctx context.Context, req *proto.ListReservationsRequest) (*proto.ListReservationsResponse, error) {
	return service.resService.ListReservations(ctx, req)
}

func (service GatewayServiceImpl) UpdateReservation(ctx context.Context, req *proto.UpdateReservationRequest) (*proto.UpdateReservationResponse, error) {
	return service.resService.UpdateReservation(ctx, req)
}

func (service GatewayServiceImpl) DeleteReservation(ctx context.Context, req *proto.DeleteReservationRequest) (*proto.DeleteReservationResponse, error) {
	return service.resService.DeleteReservation(ctx, req)
}
