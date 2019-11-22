package infrastructure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/michaldziurowski/tech-challenge-time/server/timetracking/usecases"
)

const USERID string = "user@domain.com"

func HttpHandler() http.Handler {
	inMemoryStorage := NewInMemoryStorage()
	dateProvider := NewDateProvider()
	service := usecases.NewService(inMemoryStorage, inMemoryStorage, dateProvider)

	router := httprouter.New()

	router.GET("/ping", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		_, _ = w.Write([]byte("pong"))
	})

	router.POST("/api/v1/sessions", addSession(service))
	router.PATCH("/api/v1/sessions/:sessionId", setSessionName(service))
	router.POST("/api/v1/sessions/:sessionId/stop", stopSession(service))
	router.POST("/api/v1/sessions/:sessionId/resume", resumeSession(service))
	router.GET("/api/v1/sessions", getSessions(service))

	return router
}

type addSessionRequest struct {
	Name string
	Time time.Time
}

func addSession(s usecases.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		decoder := json.NewDecoder(r.Body)

		var sessionRequest addSessionRequest
		err := decoder.Decode(&sessionRequest)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sessionId, _ := s.StartSession(USERID, sessionRequest.Name, sessionRequest.Time)

		w.Header().Set("Location", fmt.Sprintf("/sessions/%v", sessionId))
		w.WriteHeader(http.StatusCreated)
	}
}

type setSessionNameRequest struct {
	NewName string
}

func setSessionName(s usecases.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		decoder := json.NewDecoder(r.Body)
		sessionId, err := strconv.Atoi(p.ByName("sessionId"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var setNameRequest setSessionNameRequest
		err = decoder.Decode(&setNameRequest)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_ = s.SetSessionName(USERID, int64(sessionId), setNameRequest.NewName)
		w.WriteHeader(http.StatusOK)
	}
}

type sessionEventRequest struct {
	Time time.Time
}

func stopSession(s usecases.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		decoder := json.NewDecoder(r.Body)
		sessionId, err := strconv.Atoi(p.ByName("sessionId"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var stopRequest sessionEventRequest
		err = decoder.Decode(&stopRequest)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_ = s.StopSession(USERID, int64(sessionId), stopRequest.Time)
		w.WriteHeader(http.StatusOK)
	}
}

func resumeSession(s usecases.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		decoder := json.NewDecoder(r.Body)
		sessionId, err := strconv.Atoi(p.ByName("sessionId"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var resumeRequest sessionEventRequest
		err = decoder.Decode(&resumeRequest)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_ = s.ResumeSession(USERID, int64(sessionId), resumeRequest.Time)
		w.WriteHeader(http.StatusOK)
	}
}

func getSessions(s usecases.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		queryValues := r.URL.Query()
		from, err := time.Parse(time.RFC3339, queryValues.Get("from"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		to, err := time.Parse(time.RFC3339, queryValues.Get("to"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sessions, err := s.GetSessionsByRange(USERID, from, to)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		sessionsJson, _ := json.Marshal(sessions)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(sessionsJson)
	}
}
