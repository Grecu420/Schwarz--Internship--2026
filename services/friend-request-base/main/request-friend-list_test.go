package main_test

import (
	"Schwarz--Internship--2026/services/friend-request-base/main"
	"Schwarz--Internship--2026/services/friend-request-base/main/proto"
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestListFriendRequests(t *testing.T) {
	customTime := time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)

	baseRequest := &proto.ListFriendRequestsRequest{
		PageSize: 2,
	}

	emptyHash := main.HashFilter("", "", proto.RequestStatus_STATUS_UNKNOWN)
	validToken, _ := main.BuildNextPageToken(10, emptyHash)
	mismatchedToken, _ := main.BuildNextPageToken(10, "mismatched_hash")
	// Pre-computed hash for filtered token test
	senderAndStatusHash := main.HashFilter("user_a", "", proto.RequestStatus_STATUS_PENDING)
	filteredValidToken, _ := main.BuildNextPageToken(5, senderAndStatusHash)

	tests := []struct {
		name         string
		request      *proto.ListFriendRequestsRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.ListFriendRequestsResponse, req *proto.ListFriendRequestsRequest)
	}{
		{name: "Success_SinglePage",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest) {
				rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "status", "created_at"}).
					AddRow(1, "user_a", "user_b", proto.RequestStatus_STATUS_PENDING, customTime).
					AddRow(2, "user_c", "user_d", proto.RequestStatus_STATUS_ACCEPTED, customTime)

				query := `SELECT "id", "sender_id", "receiver_id", "status", "created_at" FROM "friend_requests" WHERE id >= $1 ORDER BY id LIMIT $2;`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(0), int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListFriendRequestsResponse, req *proto.ListFriendRequestsRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Requests) != 2 {
					t.Errorf("expected 2 requests, got %d", len(res.Requests))
				}
				if res.NextPageToken != "" {
					t.Errorf("expected empty nextPageToken, got %s", res.NextPageToken)
				}
				if res.Requests[0].Id != 1 || res.Requests[0].SenderId != "user_a" {
					t.Errorf("unexpected first request content: %v", res.Requests[0])
				}
			},
		},
		{name: "Success_HasNextPageToken",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest) {
				rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "status", "created_at"}).
					AddRow(1, "user_a", "user_b", proto.RequestStatus_STATUS_PENDING, customTime).
					AddRow(2, "user_c", "user_d", proto.RequestStatus_STATUS_ACCEPTED, customTime).
					AddRow(3, "user_e", "user_f", proto.RequestStatus_STATUS_PENDING, customTime)

				query := `SELECT "id", "sender_id", "receiver_id", "status", "created_at" FROM "friend_requests" WHERE id >= $1 ORDER BY id LIMIT $2;`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(0), int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListFriendRequestsResponse, req *proto.ListFriendRequestsRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Requests) != 2 {
					t.Errorf("expected 2 sliced requests, got %d", len(res.Requests))
				}
				if res.NextPageToken == "" {
					t.Error("expected non-empty nextPageToken, got empty")
				}
			},
		},
		{name: "Success_WithNextPageToken",
			request: &proto.ListFriendRequestsRequest{
				PageSize:      2,
				NextPageToken: validToken,
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest) {
				rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "status", "created_at"}).
					AddRow(10, "user_x", "user_y", proto.RequestStatus_STATUS_PENDING, customTime)

				query := `SELECT "id", "sender_id", "receiver_id", "status", "created_at" FROM "friend_requests" WHERE id >= $1 ORDER BY id LIMIT $2;`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(10), int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListFriendRequestsResponse, req *proto.ListFriendRequestsRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Requests) != 1 {
					t.Errorf("expected 1 request, got %d", len(res.Requests))
				}
				if res.NextPageToken != "" {
					t.Errorf("expected empty nextPageToken, got %s", res.NextPageToken)
				}
			},
		},
		{name: "Success_FilterBySenderId",
			request: &proto.ListFriendRequestsRequest{
				PageSize: 2,
				Filters: []*proto.ListFriendRequestsFiltersOneOf{
					{
						Filter: &proto.ListFriendRequestsFiltersOneOf_SenderId{
							SenderId: "user_a",
						},
					},
				},
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest) {
				rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "status", "created_at"}).
					AddRow(1, "user_a", "user_b", proto.RequestStatus_STATUS_PENDING, customTime).
					AddRow(4, "user_a", "user_z", proto.RequestStatus_STATUS_ACCEPTED, customTime)

				query := `SELECT "id", "sender_id", "receiver_id", "status", "created_at" FROM "friend_requests" WHERE sender_id = $1 AND id >= $2 ORDER BY id LIMIT $3;`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("user_a", int64(0), int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListFriendRequestsResponse, req *proto.ListFriendRequestsRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Requests) != 2 {
					t.Fatalf("expected 2 requests, got %d", len(res.Requests))
				}
				for _, r := range res.Requests {
					if r.SenderId != "user_a" {
						t.Errorf("expected SenderId user_a, got %s", r.SenderId)
					}
				}
			},
		},
		{name: "Success_FilterByStatus",
			request: &proto.ListFriendRequestsRequest{
				PageSize: 2,
				Filters: []*proto.ListFriendRequestsFiltersOneOf{
					{
						Filter: &proto.ListFriendRequestsFiltersOneOf_Status{
							Status: proto.RequestStatus_STATUS_PENDING,
						},
					},
				},
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest) {
				rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "status", "created_at"}).
					AddRow(1, "user_a", "user_b", proto.RequestStatus_STATUS_PENDING, customTime)

				query := `SELECT "id", "sender_id", "receiver_id", "status", "created_at" FROM "friend_requests" WHERE status = $1 AND id >= $2 ORDER BY id LIMIT $3;`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(proto.RequestStatus_STATUS_PENDING, int64(0), int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListFriendRequestsResponse, req *proto.ListFriendRequestsRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Requests) != 1 {
					t.Fatalf("expected 1 request, got %d", len(res.Requests))
				}
				if res.Requests[0].Status != proto.RequestStatus_STATUS_PENDING {
					t.Errorf("expected status STATUS_PENDING, got %v", res.Requests[0].Status)
				}
			},
		},
		{name: "Success_CombinedFiltersAndPaginationToken",
			request: &proto.ListFriendRequestsRequest{
				PageSize: 2,
				Filters: []*proto.ListFriendRequestsFiltersOneOf{
					{
						Filter: &proto.ListFriendRequestsFiltersOneOf_SenderId{
							SenderId: "user_a",
						},
					},
					{
						Filter: &proto.ListFriendRequestsFiltersOneOf_Status{
							Status: proto.RequestStatus_STATUS_PENDING,
						},
					},
				},
				NextPageToken: filteredValidToken,
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest) {
				rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "status", "created_at"}).
					AddRow(5, "user_a", "user_x", proto.RequestStatus_STATUS_PENDING, customTime)

				query := `SELECT "id", "sender_id", "receiver_id", "status", "created_at" FROM "friend_requests" WHERE sender_id = $1 AND status = $2 AND id >= $3 ORDER BY id LIMIT $4;`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("user_a", proto.RequestStatus_STATUS_PENDING, int64(5), int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListFriendRequestsResponse, req *proto.ListFriendRequestsRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Requests) != 1 {
					t.Fatalf("expected 1 request, got %d", len(res.Requests))
				}
				if res.Requests[0].Id != 5 {
					t.Errorf("expected starting offset ID 5, got %d", res.Requests[0].Id)
				}
			},
		},
		{name: "FilterMismatch_TokenHashDoesNotMatchCurrentFilters",
			request: &proto.ListFriendRequestsRequest{
				PageSize: 2,
				Filters: []*proto.ListFriendRequestsFiltersOneOf{
					{
						Filter: &proto.ListFriendRequestsFiltersOneOf_SenderId{
							SenderId: "user_a",
						},
					},
				},
				NextPageToken: validToken, // validToken was computed for empty filters
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "DatabaseError",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest) {
				query := `SELECT "id", "sender_id", "receiver_id", "status", "created_at" FROM "friend_requests" WHERE id >= $1 ORDER BY id LIMIT $2;`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(0), int64(3)).
					WillReturnError(errors.New("db query execution failed"))
			},
			expectedCode: codes.Internal,
			validate: func(t *testing.T, res *proto.ListFriendRequestsResponse, req *proto.ListFriendRequestsRequest) {
				if res != nil {
					t.Errorf("expected nil response on error, got %v", res)
				}
			},
		},
		{name: "InvalidPageSize_Zero",
			request: &proto.ListFriendRequestsRequest{
				PageSize: 0,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "InvalidPageSize_Negative",
			request: &proto.ListFriendRequestsRequest{
				PageSize: -5,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "MalformedNextPageToken",
			request: &proto.ListFriendRequestsRequest{
				PageSize:      2,
				NextPageToken: "invalid_base64_token",
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "NextPageToken_FilterMismatch",
			request: &proto.ListFriendRequestsRequest{
				PageSize:      2,
				NextPageToken: mismatchedToken,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListFriendRequestsRequest) {},
			expectedCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to initialize sqlmock: %v", err)
			}
			defer db.Close()

			if tt.setupMock != nil {
				tt.setupMock(mock, tt.request)
			}

			svc := main.FriendRequestServiceImpl{DB: db}
			res, err := svc.ListFriendRequests(context.Background(), tt.request)

			if tt.expectedCode != codes.OK {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected gRPC status error, got non-status error: %v", err)
				}
				if st.Code() != tt.expectedCode {
					t.Errorf("expected gRPC code %v, got %v", tt.expectedCode, st.Code())
				}
			} else if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, res, tt.request)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled sqlmock expectations: %s", err)
			}
		})
	}
}
