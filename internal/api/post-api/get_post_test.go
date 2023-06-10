package papi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/v1tbrah/api-gateway/internal/api/post-api/mocks"

	"github.com/v1tbrah/post-service/ppbapi"
)

func TestPostAPI_GetPost(t *testing.T) {
	tests := []struct {
		name              string
		id                string
		postServiceClient func(t *testing.T) *mocks.PostServiceClient

		expectedResp GetPostResp
		expectedCode int
		wantErr      bool
		expectedErr  error
	}{
		{
			name: "invalid id",
			id:   "invalid_id_string",
			postServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				return mocks.NewPostServiceClient(t)
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
			expectedErr:  errInvalidID,
		},
		{
			name: "not found by id",
			id:   "1",
			postServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("GetPost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.GetPostRequest{Id: int64(1)}).
					Return(nil, status.Error(codes.NotFound, "not found")).Once()
				return c
			},
			expectedCode: http.StatusNotFound,
			wantErr:      true,
			expectedErr:  errors.New("not found"),
		},
		{
			name: "unexpected err on 'postServiceClient.GetPost'",
			id:   "1",
			postServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("GetPost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.GetPostRequest{Id: int64(1)}).
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
			postServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("GetPost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.GetPostRequest{Id: int64(1)}).
					Return(&ppbapi.GetPostResponse{
						UserID: 1, Description: "TestDescription", HashtagsID: []int64{1, 2, 3}, CreatedAt: timestamppb.New(time.Unix(100, 0).UTC()),
					}, nil).Once()
				return c
			},
			expectedResp: GetPostResp{UserID: 1, Description: "TestDescription", HashtagsID: []int64{1, 2, 3}, CreatedAt: time.Unix(100, 0).UTC()},
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &PostAPI{
				postServiceClient: tt.postServiceClient(t),
			}

			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			reqCtx := chi.NewRouteContext()
			reqCtx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqCtx))

			respWriter := httptest.NewRecorder()
			handler := http.HandlerFunc(a.GetPost)
			handler.ServeHTTP(respWriter, req)

			require.Equal(t, tt.expectedCode, respWriter.Code)

			if tt.wantErr {
				assert.True(t, strings.Contains(respWriter.Body.String(), tt.expectedErr.Error()))
			} else {
				var actualResp GetPostResp
				if errDecoding := json.NewDecoder(respWriter.Body).Decode(&actualResp); errDecoding != nil {
					t.Fatal(errDecoding)
				}
				assert.Equal(t, tt.expectedResp, actualResp)
			}
		})
	}
}
