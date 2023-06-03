package papi

import (
	"bytes"
	"context"
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

func TestPostAPI_AddHashtagToPost(t *testing.T) {
	tests := []struct {
		name              string
		payload           string
		userServiceClient func(t *testing.T) *mocks.PostServiceClient

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
			name:    "invalid arg on 'postServiceClient.AddHashtagToPost'",
			payload: "{\"post_id\": 1, \"hashtag_id\": 2}",
			userServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("AddHashtagToPost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.AddHashtagToPostRequest{PostID: 1, HashtagID: 2}).
					Return(nil, status.Error(codes.InvalidArgument, "invalid id")).Once()
				return c
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errors.New("invalid id"),
		},
		{
			name:    "unexpected err on 'postServiceClient.AddHashtagToPost'",
			payload: "{\"post_id\": 1, \"hashtag_id\": 2}",
			userServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("AddHashtagToPost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.AddHashtagToPostRequest{PostID: 1, HashtagID: 2}).
					Return(nil, status.Error(codes.Internal, "unexpected err")).Once()
				return c
			},
			expectedCode: http.StatusInternalServerError,
			wantErr:      true,
			expectedErr:  errors.New("unexpected err"),
		},
		{
			name:    "OK",
			payload: "{\"post_id\": 1, \"hashtag_id\": 2}",
			userServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("AddHashtagToPost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.AddHashtagToPostRequest{PostID: 1, HashtagID: 2}).
					Return(&ppbapi.Empty{}, nil).Once()
				return c
			},
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

			handler := http.HandlerFunc(a.AddHashtagToPost)
			handler.ServeHTTP(respWriter, req)

			require.Equal(t, tt.expectedCode, respWriter.Code)

			if tt.wantErr {
				assert.True(t, strings.Contains(respWriter.Body.String(), tt.expectedErr.Error()))
			}
		})
	}
}
