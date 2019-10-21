package transport

import (
	"github.com/moorara/microservices-demo/services/asset/internal/model"
)

const (
	createAlarm  = "createAlarm"
	allAlarm     = "allAlarm"
	getAlarm     = "getAlarm"
	updateAlarm  = "updateAlarm"
	deleteAlarm  = "deleteAlarm"
	createCamera = "createCamera"
	allCamera    = "allCamera"
	getCamera    = "getCamera"
	updateCamera = "updateCamera"
	deleteCamera = "deleteCamera"
)

type (
	request struct {
		Kind string `json:"kind"`
		Span string `json:"span,omitempty"`
	}
	response struct {
		Kind  string `json:"kind"`
		Error error  `json:"error,omitempty"`
	}

	createAlarmRequest struct {
		request
		Input model.AlarmInput `json:"input"`
	}
	createAlarmResponse struct {
		response
		Alarm *model.Alarm `json:"alarm"`
	}

	allAlarmRequest struct {
		request
		SiteID string `json:"siteId"`
	}
	allAlarmResponse struct {
		response
		Alarms []model.Alarm `json:"alarms"`
	}

	getAlarmRequest struct {
		request
		ID string `json:"id"`
	}
	getAlarmResponse struct {
		response
		Alarm *model.Alarm `json:"alarm"`
	}

	updateAlarmRequest struct {
		request
		ID    string           `json:"id"`
		Input model.AlarmInput `json:"input"`
	}
	updateAlarmResponse struct {
		response
		Updated bool `json:"updated"`
	}

	deleteAlarmRequest struct {
		request
		ID string `json:"id"`
	}
	deleteAlarmResponse struct {
		response
		Deleted bool `json:"deleted"`
	}

	createCameraRequest struct {
		request
		Input model.CameraInput `json:"input"`
	}
	createCameraResponse struct {
		response
		Camera *model.Camera `json:"camera"`
	}

	allCameraRequest struct {
		request
		SiteID string `json:"siteId"`
	}
	allCameraResponse struct {
		response
		Cameras []model.Camera `json:"cameras"`
	}

	getCameraRequest struct {
		request
		ID string `json:"id"`
	}
	getCameraResponse struct {
		response
		Camera *model.Camera `json:"camera"`
	}

	updateCameraRequest struct {
		request
		ID    string            `json:"id"`
		Input model.CameraInput `json:"input"`
	}
	updateCameraResponse struct {
		response
		Updated bool `json:"updated"`
	}

	deleteCameraRequest struct {
		request
		ID string `json:"id"`
	}
	deleteCameraResponse struct {
		response
		Deleted bool `json:"deleted"`
	}
)
