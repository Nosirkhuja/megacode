package api

import (
	"MegaCode/internal/pkg/model"
	"MegaCode/internal/pkg/pg"
	"encoding/json"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
)

type API struct {
	Logger zerolog.Logger
	Db     pg.Database
}

func NewAPI(logger zerolog.Logger, db pg.Database) *API {
	return &API{Logger: logger, Db: db}
}

func (a *API) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "bad method", http.StatusBadRequest)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.Logger.Error().Msgf("read body error %v", err)
	}
	keyVal := make(map[string]string)
	err = json.Unmarshal(body, &keyVal)
	userinfo := new(model.User)
	err = json.Unmarshal(body, &userinfo)
	if err != nil {
		a.Logger.Error().Msgf("unmarshalling error %v", err)
	}
	err = a.Db.Insert(*userinfo)
	if err != nil {
		a.Logger.Error().Msgf("inserting error %v", err)
	}
}
