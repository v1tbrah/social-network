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

	"github.com/v1tbrah/api-gateway/internal/api/post-api/mocks"

	"github.com/v1tbrah/post-service/ppbapi"
)

func TestUserAPI_CreateHashtag(t *testing.T) {
	tests := []struct {
		name              string
		payload           string
		postServiceClient func(t *testing.T) *mocks.PostServiceClient

		expectedResp CreateHashtagResp
		expectedCode int
		wantErr      bool
		expectedErr  error
	}{
		{
			name:    "invalid payload",
			payload: "{__+{}",
			postServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				return mocks.NewPostServiceClient(t)
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errors.New("deserialize req"),
		},
		{
			name:    "empty name",
			payload: "{\"name\": \"\"}",
			postServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				return mocks.NewPostServiceClient(t)
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errEmptyName,
		},
		{
			name:    "empty name with spaces",
			payload: "{\"name\": \"   \"}",
			postServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				return mocks.NewPostServiceClient(t)
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errEmptyName,
		},
		{
			name:    "invalid arg on 'postServiceClient.CreateHashtag'",
			payload: "{\"name\": \"TestName\"}",
			postServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("CreateHashtag",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.CreateHashtagRequest{Name: "TestName"}).
					Return(nil, status.Error(codes.InvalidArgument, "invalid name")).Once()
				return c
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errors.New("invalid name"),
		},
		{
			name:    "already exists on 'postServiceClient.CreateHashtag'",
			payload: "{\"name\": \"TestName\"}",
			postServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("CreateHashtag",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.CreateHashtagRequest{Name: "TestName"}).
					Return(nil, status.Error(codes.AlreadyExists, "hashtag already exists")).Once()
				return c
			},
			expectedCode: http.StatusConflict,
			wantErr:      true,
			expectedErr:  errors.New("hashtag already exists"),
		},
		{
			name:    "unexpected err on 'postServiceClient.CreatePost'",
			payload: "{\"name\": \"TestName\"}",
			postServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("CreateHashtag",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.CreateHashtagRequest{Name: "TestName"}).
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
			postServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("CreateHashtag",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.CreateHashtagRequest{Name: "TestName"}).
					Return(&ppbapi.CreateHashtagResponse{Id: int64(1)}, nil).Once()
				return c
			},
			expectedResp: CreateHashtagResp{ID: int64(1)},
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &PostAPI{
				postServiceClient: tt.postServiceClient(t),
			}

			respWriter := httptest.NewRecorder()

			var body bytes.Buffer
			body.WriteString(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/test", &body)

			handler := http.HandlerFunc(a.CreateHashtag)
			handler.ServeHTTP(respWriter, req)

			require.Equal(t, tt.expectedCode, respWriter.Code)

			if tt.wantErr {
				assert.True(t, strings.Contains(respWriter.Body.String(), tt.expectedErr.Error()))
			} else {
				var actualResp CreateHashtagResp
				if errDecoding := json.NewDecoder(respWriter.Body).Decode(&actualResp); errDecoding != nil {
					t.Fatal(errDecoding)
				}
				assert.Equal(t, tt.expectedResp, actualResp)
			}
		})
	}
}
