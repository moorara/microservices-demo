package transport

import (
	"github.com/moorara/microservices-demo/services/asset-service/internal/model"
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
		Kind string
		Span string
	}

	createAlarmRequest struct {
		request
		input model.AlarmInput
	}

	allAlarmRequest struct {
		request
		SiteID string
	}

	getAlarmRequest struct {
		request
		ID string
	}

	updateAlarmRequest struct {
		request
		ID    string
		input model.AlarmInput
	}

	deleteAlarmRequest struct {
		request
		ID string
	}

	createCameraRequest struct {
		request
		input model.CameraInput
	}

	allCameraRequest struct {
		request
		SiteID string
	}

	getCameraRequest struct {
		request
		ID string
	}

	updateCameraRequest struct {
		request
		ID    string
		input model.CameraInput
	}

	deleteCameraRequest struct {
		request
		ID string
	}
)
