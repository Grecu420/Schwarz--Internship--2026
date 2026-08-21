package main_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"Schwarz--Internship--2026/services/user-base/main"
	"Schwarz--Internship--2026/services/user-base/main/proto"
)

func TestListUsers(t *testing.T) {

	baseRequest := &proto.ListUsersRequest{
		PageSize: 2,
	}

	emptyHash := main.HashFilter("", "")
	validToken, _ := main.BuildNextPageToken(10, emptyHash)
	mismatchedToken, _ := main.BuildNextPageToken(10, "mismatched_hash")

	firstNameHash := main.HashFilter("Andrei", "")
	filteredValidToken, _ := main.BuildNextPageToken(5, firstNameHash)

	tests := []struct {
		name         string
		request      *proto.ListUsersRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.ListUsersRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.ListUsersResponse, req *proto.ListUsersRequest)
	}{
		{
			name:    "Success_SinglePage",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListUsersRequest) {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "user_name", "email"}).
					AddRow(1, "Andrei", "Ivancu", "aivancu", "andrei@test.com").
					AddRow(2, "Ion", "Popescu", "ipopescu", "ion@test.com")

				query := `SELECT id, first_name, last_name, user_name, email FROM users ORDER BY id ASC LIMIT $1;`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(3)). 
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListUsersResponse, req *proto.ListUsersRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Users) != 2 {
					t.Errorf("expected 2 users, got %d", len(res.Users))
				}
				if res.NextPageToken != "" {
					t.Errorf("expected empty nextPageToken, got %s", res.NextPageToken)
				}
				if res.Users[0].Id != 1 || res.Users[0].FirstName != "Andrei" {
					t.Errorf("unexpected first user content: %v", res.Users[0])
				}
			},
		},
		{
			name:    "Success_HasNextPageToken",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListUsersRequest) {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "user_name", "email"}).
					AddRow(1, "Andrei", "Ivancu", "aivancu", "andrei@test.com").
					AddRow(2, "Ion", "Popescu", "ipopescu", "ion@test.com").
					AddRow(3, "Maria", "Ionescu", "mionescu", "maria@test.com")

				query := `SELECT id, first_name, last_name, user_name, email FROM users ORDER BY id ASC LIMIT $1;`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListUsersResponse, req *proto.ListUsersRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Users) != 2 { 
					t.Errorf("expected 2 sliced users, got %d", len(res.Users))
				}
				if res.NextPageToken == "" {
					t.Error("expected non-empty nextPageToken, got empty")
				}
			},
		},
		{
			name: "Success_WithNextPageToken",
			request: &proto.ListUsersRequest{
				PageSize:      2,
				NextPageToken: validToken,
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListUsersRequest) {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "user_name", "email"}).
					AddRow(15, "George", "Enescu", "genescu", "george@test.com")

				query := `SELECT id, first_name, last_name, user_name, email FROM users WHERE id > $1 ORDER BY id ASC LIMIT $2;`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(10), int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListUsersResponse, req *proto.ListUsersRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Users) != 1 {
					t.Errorf("expected 1 user, got %d", len(res.Users))
				}
				if res.NextPageToken != "" {
					t.Errorf("expected empty nextPageToken, got %s", res.NextPageToken)
				}
			},
		},
		{
			name: "Success_FilterByFirstName",
			request: &proto.ListUsersRequest{
				PageSize: 2,
				Filters: []*proto.ListUsersFiltersOneOf{
					{
						Filter: &proto.ListUsersFiltersOneOf_FirstName{
							FirstName: &proto.FilterByFirstName{Value: "Andrei"},
						},
					},
				},
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListUsersRequest) {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "user_name", "email"}).
					AddRow(1, "Andrei", "Ivancu", "aivancu", "andrei@test.com").
					AddRow(4, "Andrei", "Tarkovsky", "atarkovsky", "tarkovsky@test.com")

				query := `SELECT id, first_name, last_name, user_name, email FROM users WHERE first_name = $1 ORDER BY id ASC LIMIT $2;`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("Andrei", int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListUsersResponse, req *proto.ListUsersRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Users) != 2 {
					t.Fatalf("expected 2 users, got %d", len(res.Users))
				}
				for _, u := range res.Users {
					if u.FirstName != "Andrei" {
						t.Errorf("expected FirstName Andrei, got %s", u.FirstName)
					}
				}
			},
		},
		{
			name: "Success_CombinedFiltersAndPaginationToken",
			request: &proto.ListUsersRequest{
				PageSize: 2,
				Filters: []*proto.ListUsersFiltersOneOf{
					{
						Filter: &proto.ListUsersFiltersOneOf_FirstName{
							FirstName: &proto.FilterByFirstName{Value: "Andrei"},
						},
					},
				},
				NextPageToken: filteredValidToken,
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListUsersRequest) {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "user_name", "email"}).
					AddRow(6, "Andrei", "Popa", "apopa", "apopa@test.com")

				query := `SELECT id, first_name, last_name, user_name, email FROM users WHERE first_name = $1 AND id > $2 ORDER BY id ASC LIMIT $3;`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("Andrei", int64(5), int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListUsersResponse, req *proto.ListUsersRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Users) != 1 {
					t.Fatalf("expected 1 user, got %d", len(res.Users))
				}
				if res.Users[0].Id != 6 {
					t.Errorf("expected user ID 6, got %d", res.Users[0].Id)
				}
			},
		},
		{
			name: "FilterMismatch_TokenHashDoesNotMatchCurrentFilters",
			request: &proto.ListUsersRequest{
				PageSize: 2,
				Filters: []*proto.ListUsersFiltersOneOf{
					{
						Filter: &proto.ListUsersFiltersOneOf_FirstName{
							FirstName: &proto.FilterByFirstName{Value: "Andrei"},
						},
					},
				},
				NextPageToken: validToken,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListUsersRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "FilterMismatch_DuplicateFilter",
			request: &proto.ListUsersRequest{
				PageSize: 2,
				Filters: []*proto.ListUsersFiltersOneOf{
					{
						Filter: &proto.ListUsersFiltersOneOf_FirstName{
							FirstName: &proto.FilterByFirstName{Value: "Andrei"},
						},
					},
					{
						Filter: &proto.ListUsersFiltersOneOf_FirstName{
							FirstName: &proto.FilterByFirstName{Value: "Ion"},
						},
					},
				},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListUsersRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:    "DatabaseError",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListUsersRequest) {
				query := `SELECT id, first_name, last_name, user_name, email FROM users ORDER BY id ASC LIMIT $1;`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(3)).
					WillReturnError(errors.New("db query execution failed"))
			},
			expectedCode: codes.Internal,
			validate: func(t *testing.T, res *proto.ListUsersResponse, req *proto.ListUsersRequest) {
				if res != nil {
					t.Errorf("expected nil response on error, got %v", res)
				}
			},
		},
		{
			name: "InvalidPageSize_Zero",
			request: &proto.ListUsersRequest{
				PageSize: 0,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListUsersRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MalformedNextPageToken",
			request: &proto.ListUsersRequest{
				PageSize:      2,
				NextPageToken: "invalid_base64_token",
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListUsersRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListUsersRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "NextPageToken_FilterMismatch",
			request: &proto.ListUsersRequest{
				PageSize:      2,
				NextPageToken: mismatchedToken,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListUsersRequest) {},
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

			
			svc := &main.UserServiceImpl{DB: db}
			res, err := svc.ListUsers(context.Background(), tt.request)

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
