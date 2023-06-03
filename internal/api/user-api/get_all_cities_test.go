package uapi

import (
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
	"gitlab.com/pet-pr-social-network/user-service/pbapi"
)

func TestUserAPI_GetAllCities(t *testing.T) {
	tests := []struct {
		name              string
		userServiceClient func(t *testing.T) *mocks.UserServiceClient

		expectedResp []GetCityResp
		expectedCode int
		wantErr      bool
		expectedErr  error
	}{
		{
			name: "unexpected err on 'userServiceClient.GetAllCities'",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				c := mocks.NewUserServiceClient(t)
				c.On("GetAllCities",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), &pbapi.Empty{}).
					Return(nil, status.Error(codes.Internal, "unexpected err")).Once()
				return c
			},
			expectedCode: http.StatusInternalServerError,
			wantErr:      true,
			expectedErr:  errors.New("unexpected err"),
		},
		{
			name: "OK",
			userServiceClient: func(t *testing.T) *mocks.UserServiceClient {
				c := mocks.NewUserServiceClient(t)
				c.On("GetAllCities",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), &pbapi.Empty{}).
					Return(&pbapi.GetAllCitiesResponse{Cities: []*pbapi.City{{Id: 1, Name: "TestCity"}}}, nil).Once()
				return c
			},
			expectedResp: []GetCityResp{{ID: int64(1), Name: "TestCity"}},
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UserAPI{
				userServiceClient: tt.userServiceClient(t),
			}

			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			respWriter := httptest.NewRecorder()
			handler := http.HandlerFunc(a.GetAllCities)
			handler.ServeHTTP(respWriter, req)

			require.Equal(t, tt.expectedCode, respWriter.Code)

			if tt.wantErr {
				assert.True(t, strings.Contains(respWriter.Body.String(), tt.expectedErr.Error()))
			} else {
				var actualResp []GetCityResp
				if errDecoding := json.NewDecoder(respWriter.Body).Decode(&actualResp); errDecoding != nil {
					t.Fatal(errDecoding)
				}
				assert.Equal(t, tt.expectedResp, actualResp)
			}
		})
	}
}
