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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/pet-pr-social-network/api-gateway/internal/api/user-api/mocks"
	"gitlab.com/pet-pr-social-network/user-service/upbapi"
)

func TestUserAPI_CreateInterest(t *testing.T) {
	tests := []struct {
		name              string
		payload           string
		userServiceClient func(t *testing.T) *mocks.UserServiceClient

		expectedResp CreateInterestResp
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
			name:    "invalid arg on 'userServiceClient.CreateInterest'",
			payload: "{\"name\": \"TestName\"}",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				c := mocks.NewUserServiceClient(t)
				c.On("CreateInterest",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&upbapi.CreateInterestRequest{Name: "TestName"}).
					Return(nil, status.Error(codes.InvalidArgument, "invalid name")).Once()
				return c
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errors.New("invalid name"),
		},
		{
			name:    "already exists on 'userServiceClient.CreateInterest'",
			payload: "{\"name\": \"TestName\"}",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				c := mocks.NewUserServiceClient(t)
				c.On("CreateInterest",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&upbapi.CreateInterestRequest{Name: "TestName"}).
					Return(nil, status.Error(codes.AlreadyExists, "interest already exists")).Once()
				return c
			},
			expectedCode: http.StatusConflict,
			wantErr:      true,
			expectedErr:  errors.New("interest already exists"),
		},
		{
			name:    "unexpected err on 'userServiceClient.CreateInterest'",
			payload: "{\"name\": \"TestName\"}",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				c := mocks.NewUserServiceClient(t)
				c.On("CreateInterest",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&upbapi.CreateInterestRequest{Name: "TestName"}).
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
				c.On("CreateInterest",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&upbapi.CreateInterestRequest{Name: "TestName"}).
					Return(&upbapi.CreateInterestResponse{Id: int64(1)}, nil).Once()
				return c
			},
			expectedResp: CreateInterestResp{ID: int64(1)},
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

			handler := http.HandlerFunc(a.CreateInterest)
			handler.ServeHTTP(respWriter, req)

			require.Equal(t, tt.expectedCode, respWriter.Code)

			if tt.wantErr {
				assert.True(t, strings.Contains(respWriter.Body.String(), tt.expectedErr.Error()))
			} else {
				var actualResp CreateInterestResp
				if errDecoding := json.NewDecoder(respWriter.Body).Decode(&actualResp); errDecoding != nil {
					t.Fatal(errDecoding)
				}
				assert.Equal(t, tt.expectedResp, actualResp)
			}
		})
	}
}
