package uapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.com/pet-pr-social-network/api-gateway/internal/api/user-api/mocks"
	"gitlab.com/pet-pr-social-network/user-service/pbapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserAPI_CreateCity(t *testing.T) {
	tests := []struct {
		name              string
		payload           string
		userServiceClient func(t *testing.T) *mocks.UserServiceClient

		expectedResp CreateCityResp
		expectedCode int
		wantErr      bool
		expectedErr  error
	}{
		{
			name:    "invalid payload",
			payload: "{__+{}",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				return mocks.NewUserServiceClient(t)
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errors.New("deserialize req"),
		},
		{
			name:    "empty name",
			payload: "{\"name\": \"\"}",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				return mocks.NewUserServiceClient(t)
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errEmptyName,
		},
		{
			name:    "empty name with spaces",
			payload: "{\"name\": \"   \"}",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				return mocks.NewUserServiceClient(t)
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errEmptyName,
		},
		{
			name:    "invalid arg on 'userServiceClient.CreateCity'",
			payload: "{\"name\": \"TestName\"}",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				c := mocks.NewUserServiceClient(t)
				c.On("CreateCity",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&pbapi.CreateCityRequest{Name: "TestName"}).
					Return(nil, status.Error(codes.InvalidArgument, "invalid name")).Once()
				return c
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errors.New("invalid name"),
		},
		{
			name:    "already exists on 'userServiceClient.CreateCity'",
			payload: "{\"name\": \"TestName\"}",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				c := mocks.NewUserServiceClient(t)
				c.On("CreateCity",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&pbapi.CreateCityRequest{Name: "TestName"}).
					Return(nil, status.Error(codes.AlreadyExists, "city already exists")).Once()
				return c
			},
			expectedCode: http.StatusConflict,
			wantErr:      true,
			expectedErr:  errors.New("city already exists"),
		},
		{
			name:    "unexpected err on 'userServiceClient.CreateCity'",
			payload: "{\"name\": \"TestName\"}",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				c := mocks.NewUserServiceClient(t)
				c.On("CreateCity",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&pbapi.CreateCityRequest{Name: "TestName"}).
					Return(nil, status.Error(codes.Internal, "unexpected err")).Once()
				return c
			},
			expectedCode: http.StatusInternalServerError,
			wantErr:      true,
			expectedErr:  errors.New("unexpected err"),
		},
		{
			name:    "OK",
			payload: "{\"name\": \"TestName\"}",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				c := mocks.NewUserServiceClient(t)
				c.On("CreateCity",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&pbapi.CreateCityRequest{Name: "TestName"}).
					Return(&pbapi.CreateCityResponse{Id: int64(1)}, nil).Once()
				return c
			},
			expectedResp: CreateCityResp{ID: int64(1)},
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UserAPI{
				userServiceClient: tt.userServiceClient(t),
			}

			respWriter := httptest.NewRecorder()

			var body bytes.Buffer
			body.WriteString(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/test", &body)

			handler := http.HandlerFunc(a.CreateCity)
			handler.ServeHTTP(respWriter, req)

			require.Equal(t, tt.expectedCode, respWriter.Code)

			if tt.wantErr {
				assert.True(t, strings.Contains(respWriter.Body.String(), tt.expectedErr.Error()))
			} else {
				var actualResp CreateCityResp
				if errDecoding := json.NewDecoder(respWriter.Body).Decode(&actualResp); errDecoding != nil {
					t.Fatal(errDecoding)
				}
				assert.Equal(t, tt.expectedResp, actualResp)
			}
		})
	}
}
