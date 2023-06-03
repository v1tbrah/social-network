package papi

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

	"gitlab.com/pet-pr-social-network/api-gateway/internal/api/post-api/mocks"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
)

func TestPostAPI_CreatePost(t *testing.T) {
	tests := []struct {
		name              string
		payload           string
		userServiceClient func(t *testing.T) *mocks.PostServiceClient

		expectedResp CreatePostResp
		expectedCode int
		wantErr      bool
		expectedErr  error
	}{
		{
			name:    "invalid payload",
			payload: "{__+{}",
			userServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				return mocks.NewPostServiceClient(t)
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errors.New("deserialize req"),
		},
		{
			name:    "invalid arg on 'userServiceClient.CreatePost'",
			payload: "{\"user_id\": -1}",
			userServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("CreatePost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.CreatePostRequest{UserID: -1}).
					Return(nil, status.Error(codes.InvalidArgument, "invalid user id")).Once()
				return c
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errors.New("invalid user id"),
		},
		{
			name:    "already exists on 'userServiceClient.CreatePost'",
			payload: "{\"user_id\": 1}",
			userServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("CreatePost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.CreatePostRequest{UserID: 1}).
					Return(nil, status.Error(codes.AlreadyExists, "post already exists")).Once()
				return c
			},
			expectedCode: http.StatusConflict,
			wantErr:      true,
			expectedErr:  errors.New("post already exists"),
		},
		{
			name:    "unexpected err on 'userServiceClient.CreatePost'",
			payload: "{\"user_id\": 1}",
			userServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("CreatePost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.CreatePostRequest{UserID: 1}).
					Return(nil, status.Error(codes.Internal, "unexpected err")).Once()
				return c
			},
			expectedCode: http.StatusInternalServerError,
			wantErr:      true,
			expectedErr:  errors.New("unexpected err"),
		},
		{
			name:    "OK",
			payload: "{\"user_id\": 1}",
			userServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("CreatePost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.CreatePostRequest{UserID: 1}).
					Return(&ppbapi.CreatePostResponse{Id: int64(1)}, nil).Once()
				return c
			},
			expectedResp: CreatePostResp{ID: int64(1)},
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &PostAPI{
				postServiceClient: tt.userServiceClient(t),
			}

			respWriter := httptest.NewRecorder()

			var body bytes.Buffer
			body.WriteString(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/test", &body)

			handler := http.HandlerFunc(a.CreatePost)
			handler.ServeHTTP(respWriter, req)

			require.Equal(t, tt.expectedCode, respWriter.Code)

			if tt.wantErr {
				assert.True(t, strings.Contains(respWriter.Body.String(), tt.expectedErr.Error()))
			} else {
				var actualResp CreatePostResp
				if errDecoding := json.NewDecoder(respWriter.Body).Decode(&actualResp); errDecoding != nil {
					t.Fatal(errDecoding)
				}
				assert.Equal(t, tt.expectedResp, actualResp)
			}
		})
	}
}
