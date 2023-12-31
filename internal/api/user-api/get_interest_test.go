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

	"github.com/v1tbrah/api-gateway/internal/api/user-api/mocks"

	"github.com/v1tbrah/user-service/upbapi"
)

func TestUserAPI_GetInterest(t *testing.T) {
	tests := []struct {
		name              string
		id                string
		userServiceClient func(t *testing.T) *mocks.UserServiceClient

		expectedResp GetInterestResp
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
				c.On("GetInterest",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&upbapi.GetInterestRequest{Id: int64(1)}).
					Return(nil, status.Error(codes.NotFound, "not found")).Once()
				return c
			},
			expectedCode: http.StatusNotFound,
			wantErr:      true,
			expectedErr:  errors.New("not found"),
		},
		{
			name: "unexpected err on 'userServiceClient.GetInterest'",
			id:   "1",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				c := mocks.NewUserServiceClient(t)
				c.On("GetInterest",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&upbapi.GetInterestRequest{Id: int64(1)}).
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
				c.On("GetInterest",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&upbapi.GetInterestRequest{Id: int64(1)}).
					Return(&upbapi.GetInterestResponse{Interest: &upbapi.Interest{Id: 1, Name: "TestInterest"}}, nil).Once()
				return c
			},
			expectedResp: GetInterestResp{ID: int64(1), Name: "TestInterest"},
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
			handler := http.HandlerFunc(a.GetInterest)
			handler.ServeHTTP(respWriter, req)

			require.Equal(t, tt.expectedCode, respWriter.Code)

			if tt.wantErr {
				assert.True(t, strings.Contains(respWriter.Body.String(), tt.expectedErr.Error()))
			} else {
				var actualResp GetInterestResp
				if errDecoding := json.NewDecoder(respWriter.Body).Decode(&actualResp); errDecoding != nil {
					t.Fatal(errDecoding)
				}
				assert.Equal(t, tt.expectedResp, actualResp)
			}
		})
	}
}
