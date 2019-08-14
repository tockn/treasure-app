package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/voyagegroup/treasure-app/httputil"
	"github.com/voyagegroup/treasure-app/model"
	"github.com/voyagegroup/treasure-app/service"
)

type Article struct {
	articleService service.Article
}

func NewArticle(as service.Article) *Article {
	return &Article{
		articleService: as,
	}
}

func (a *Article) Index(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	ctx := r.Context()
	articles, err := a.articleService.GetAll(ctx)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, articles, nil
}

func (a *Article) Show(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return http.StatusBadRequest, nil, &httputil.HTTPError{Message: "invalid path parameter"}
	}

	aid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	article, err := a.articleService.Get(ctx, aid)
	if err != nil && err == sql.ErrNoRows {
		return http.StatusNotFound, nil, err
	} else if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusCreated, article, nil
}

func (a *Article) Create(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	ctx := r.Context()
	newArticle := &model.Article{}
	if err := json.NewDecoder(r.Body).Decode(&newArticle); err != nil {
		return http.StatusBadRequest, nil, err
	}

	id, err := a.articleService.Create(ctx, newArticle)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	newArticle.ID = id

	return http.StatusCreated, newArticle, nil
}

func (a *Article) Update(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return http.StatusBadRequest, nil, &httputil.HTTPError{Message: "invalid path parameter"}
	}

	aid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	reqArticle := &model.Article{}
	if err := json.NewDecoder(r.Body).Decode(&reqArticle); err != nil {
		return http.StatusBadRequest, nil, err
	}

	err = a.articleService.Update(ctx, aid, reqArticle)
	if err != nil && errors.Cause(err) == sql.ErrNoRows {
		return http.StatusNotFound, nil, err
	} else if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusNoContent, nil, nil
}

func (a *Article) Destroy(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return http.StatusBadRequest, nil, &httputil.HTTPError{Message: "invalid path parameter"}
	}

	aid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	err = a.articleService.Destroy(ctx, aid)
	if err != nil && errors.Cause(err) == sql.ErrNoRows {
		return http.StatusNotFound, nil, err
	} else if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusNoContent, nil, nil
}
