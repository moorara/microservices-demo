package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/moorara/microservices-demo/services/sensor-service/service"
	"github.com/opentracing/opentracing-go"
)

type (
	// SensorHandler provides http handlers for Sensor resource
	SensorHandler interface {
		PostSensor(w http.ResponseWriter, r *http.Request)
		GetSensors(w http.ResponseWriter, r *http.Request)
		GetSensor(w http.ResponseWriter, r *http.Request)
		PutSensor(w http.ResponseWriter, r *http.Request)
		DeleteSensor(w http.ResponseWriter, r *http.Request)
	}

	postgresSensorHandler struct {
		manager service.SensorManager
		logger  log.Logger
	}
)

// NewSensorHandler creates a new sensor handler
func NewSensorHandler(db service.DB, logger log.Logger, tracer opentracing.Tracer) SensorHandler {
	return &postgresSensorHandler{
		manager: service.NewSensorManager(db, logger, tracer),
		logger:  logger,
	}
}

func (h *postgresSensorHandler) PostSensor(w http.ResponseWriter, r *http.Request) {
	s := struct {
		SiteID  string  `json:"siteId"`
		Name    string  `json:"name"`
		Unit    string  `json:"unit"`
		MinSafe float64 `json:"minSafe"`
		MaxSafe float64 `json:"maxSafe"`
	}{}

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil || s.SiteID == "" || s.Name == "" || s.Unit == "" || s.MinSafe > s.MaxSafe {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	sensor, err := h.manager.Create(r.Context(), s.SiteID, s.Name, s.Unit, s.MinSafe, s.MaxSafe)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(sensor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *postgresSensorHandler) GetSensors(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["siteId"]
	if siteID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	sensors, err := h.manager.All(r.Context(), siteID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(sensors)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *postgresSensorHandler) GetSensor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sensorID := vars["id"]
	if sensorID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	sensor, err := h.manager.Get(r.Context(), sensorID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sensor == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(sensor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *postgresSensorHandler) PutSensor(w http.ResponseWriter, r *http.Request) {
	s := service.Sensor{}

	vars := mux.Vars(r)
	s.ID = vars["id"]
	if s.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil || s.SiteID == "" || s.Name == "" || s.Unit == "" || s.MinSafe > s.MaxSafe {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	n, err := h.manager.Update(r.Context(), s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if n == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *postgresSensorHandler) DeleteSensor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := h.manager.Delete(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
