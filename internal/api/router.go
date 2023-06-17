package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/v1tbrah/api-gateway/docs"
	fapi "github.com/v1tbrah/api-gateway/internal/api/feed-api"
	mapi "github.com/v1tbrah/api-gateway/internal/api/media-api"
	papi "github.com/v1tbrah/api-gateway/internal/api/post-api"
	rapi "github.com/v1tbrah/api-gateway/internal/api/relation-api"
	uapi "github.com/v1tbrah/api-gateway/internal/api/user-api"
)

//	@title			Social-network API
//	@version		1.0
//	@description	This is a simple social-network server.

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	/
func (a *API) newRouter() (r *chi.Mux) {
	r = chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/ping", a.ping)

	r.Handle("/metrics", promhttp.Handler())

	r.Route("/user", func(r chi.Router) {

		newUAPI := uapi.New(a.userServiceClient)

		r.Post("/city", newUAPI.CreateCity)
		r.Get("/city/{id}", newUAPI.GetCity)
		r.Get("/city", newUAPI.GetAllCities)

		r.Post("/interest", newUAPI.CreateInterest)
		r.Get("/interest/{id}", newUAPI.GetInterest)
		r.Get("/interest", newUAPI.GetAllInterests)

		r.Post("/user", newUAPI.CreateUser)
		r.Get("/user/{id}", newUAPI.GetUser)

	})

	r.Route("/post", func(r chi.Router) {

		newPAPI := papi.New(a.postServiceClient)

		r.Post("/hashtag", newPAPI.CreateHashtag)
		r.Get("/hashtag/{id}", newPAPI.GetHashtag)

		r.Get("/post/{id}", newPAPI.GetPost)
		r.Post("/post", newPAPI.CreatePost)
		r.Post("/post/get_by_hashtag", newPAPI.GetPostsByHashtag)
		r.Post("/post/hashtag", newPAPI.AddHashtagToPost)
	})

	r.Route("/relation", func(r chi.Router) {

		newRAPI := rapi.New(a.relationServiceClient)

		r.Post("/friend", newRAPI.AddFriend)
		r.Delete("/friend", newRAPI.RemoveFriend)
		r.Post("/friend/get_friends_by_user", newRAPI.GetFriendsByUser)
	})

	r.Route("/feed", func(r chi.Router) {

		newFAPI := fapi.New(a.feedServiceClient)

		r.Get("/{id}", newFAPI.GetFeed)
	})

	r.Route("/media", func(r chi.Router) {

		newMAPI := mapi.New(a.mediaServiceClient)

		r.Get("/post/{guid}", newMAPI.GetPost)
		r.Post("/post", newMAPI.SavePost)
	})

	r.Mount("/swagger", httpSwagger.WrapHandler)

	return r
}
