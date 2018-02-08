package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/moorara/microservices-demo/services/sensor-service/service"
)

type (
	// SensorHandler provides http handlers for Sensor resource
	SensorHandler interface {
		PostSensor(w http.ResponseWriter, r *http.Request)
		GetSensors(w http.ResponseWriter, r *http.Request)
		GetSensor(w http.ResponseWriter, r *http.Request)
		DeleteSensor(w http.ResponseWriter, r *http.Request)
	}

	postgresSensorHandler struct {
		manager service.SensorManager
		logger  log.Logger
	}
)

// NewSensorHandler creates a new sensor handler
func NewSensorHandler(db service.DB, logger log.Logger) SensorHandler {
	return &postgresSensorHandler{
		manager: service.NewSensorManager(db, logger),
		logger:  logger,
	}
}

func (h *postgresSensorHandler) PostSensor(w http.ResponseWriter, r *http.Request) {
	req := struct {
		SiteID  string  `json:"siteId"`
		Name    string  `json:"name"`
		Unit    string  `json:"unit"`
		MinSafe float64 `json:"minSafe"`
		MaxSafe float64 `json:"maxSafe"`
	}{}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.SiteID == "" || req.Name == "" || req.Unit == "" || req.MinSafe > req.MaxSafe {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sensor, err := h.manager.Create(context.Background(), req.SiteID, req.Name, req.Unit, req.MinSafe, req.MaxSafe)
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

	sensors, err := h.manager.All(context.Background(), siteID)
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

	sensor, err := h.manager.Get(context.Background(), sensorID)
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

func (h *postgresSensorHandler) DeleteSensor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sensorID := vars["id"]
	if sensorID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := h.manager.Delete(context.Background(), sensorID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
