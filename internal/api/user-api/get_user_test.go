package uapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/pet-pr-social-network/api-gateway/internal/api/user-api/mocks"
	"gitlab.com/pet-pr-social-network/user-service/upbapi"
)

func TestUserAPI_GetUser(t *testing.T) {
	tests := []struct {
		name              string
		id                string
		userServiceClient func(t *testing.T) *mocks.UserServiceClient

		expectedResp GetUserResp
		expectedCode int
		wantErr      bool
		expectedErr  error
	}{
		{
			name: "invalid id",
			id:   "invalid_id_string",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				return mocks.NewUserServiceClient(t)
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errInvalidID,
		},
		{
			name: "not found by id",
			id:   "1",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				c := mocks.NewUserServiceClient(t)
				c.On("GetUser",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&upbapi.GetUserRequest{Id: int64(1)}).
					Return(nil, status.Error(codes.NotFound, "not found")).Once()
				return c
			},
			expectedCode: http.StatusNotFound,
			wantErr:      true,
			expectedErr:  errors.New("not found"),
		},
		{
			name: "unexpected err on 'userServiceClient.GetUser'",
			id:   "1",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				c := mocks.NewUserServiceClient(t)
				c.On("GetUser",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&upbapi.GetUserRequest{Id: int64(1)}).
					Return(nil, status.Error(codes.Internal, "unexpected err")).Once()
				return c
			},
			expectedCode: http.StatusInternalServerError,
			wantErr:      true,
			expectedErr:  errors.New("unexpected err"),
		},
		{
			name: "OK",
			id:   "1",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				c := mocks.NewUserServiceClient(t)
				c.On("GetUser",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&upbapi.GetUserRequest{Id: int64(1)}).
					Return(&upbapi.GetUserResponse{
						Name: "TestName", Surname: "TestSurname", InterestsID: []int64{1, 2, 3}, CityID: int64(1),
					}, nil).Once()
				return c
			},
			expectedResp: GetUserResp{Name: "TestName", Surname: "TestSurname", InterestsID: []int64{1, 2, 3}, CityID: int64(1)},
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UserAPI{
				userServiceClient: tt.userServiceClient(t),
			}

			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			reqCtx := chi.NewRouteContext()
			reqCtx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqCtx))

			respWriter := httptest.NewRecorder()
			handler := http.HandlerFunc(a.GetUser)
			handler.ServeHTTP(respWriter, req)

			require.Equal(t, tt.expectedCode, respWriter.Code)

			if tt.wantErr {
				assert.True(t, strings.Contains(respWriter.Body.String(), tt.expectedErr.Error()))
			} else {
				var actualResp GetUserResp
				if errDecoding := json.NewDecoder(respWriter.Body).Decode(&actualResp); errDecoding != nil {
					t.Fatal(errDecoding)
				}
				assert.Equal(t, tt.expectedResp, actualResp)
			}
		})
	}
}
