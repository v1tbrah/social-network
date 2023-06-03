package papi

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

	"gitlab.com/pet-pr-social-network/api-gateway/internal/api/post-api/mocks"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
)

func TestPostAPI_GetHashtag(t *testing.T) {
	tests := []struct {
		name              string
		id                string
		postServiceClient func(t *testing.T) *mocks.PostServiceClient

		expectedResp GetHashtagResp
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
				c.On("GetHashtag",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.GetHashtagRequest{Id: int64(1)}).
					Return(nil, status.Error(codes.NotFound, "not found")).Once()
				return c
			},
			expectedCode: http.StatusNotFound,
			wantErr:      true,
			expectedErr:  errors.New("not found"),
		},
		{
			name: "unexpected err on 'postServiceClient.GetHashtag'",
			id:   "1",
			postServiceClient: func(t *testing.T) *mocks.PostServiceClient {
				c := mocks.NewPostServiceClient(t)
				c.On("GetHashtag",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.GetHashtagRequest{Id: int64(1)}).
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
				c.On("GetHashtag",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&ppbapi.GetHashtagRequest{Id: int64(1)}).
					Return(&ppbapi.GetHashtagResponse{Hashtag: &ppbapi.Hashtag{Id: 1, Name: "TestHashtag"}}, nil).Once()
				return c
			},
			expectedResp: GetHashtagResp{ID: int64(1), Name: "TestHashtag"},
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
			handler := http.HandlerFunc(a.GetHashtag)
			handler.ServeHTTP(respWriter, req)

			require.Equal(t, tt.expectedCode, respWriter.Code)

			if tt.wantErr {
				assert.True(t, strings.Contains(respWriter.Body.String(), tt.expectedErr.Error()))
			} else {
				var actualResp GetHashtagResp
				if errDecoding := json.NewDecoder(respWriter.Body).Decode(&actualResp); errDecoding != nil {
					t.Fatal(errDecoding)
				}
				assert.Equal(t, tt.expectedResp, actualResp)
			}
		})
	}
}
