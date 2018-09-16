package transport

import (
	"github.com/moorara/microservices-demo/services/asset-service/internal/model"
)

type (
	response struct {
		Kind string
		Span string
	}

	createAlarmResponse struct {
		response
		Alarm model.Alarm
		Error error
	}

	allAlarmResponse struct {
		response
		Alarms []model.Alarm
		Error  error
	}

	getAlarmResponse struct {
		response
		Alarm model.Alarm
		Error error
	}

	updateAlarmResponse struct {
		response
		Updated bool
		Error   error
	}

	deleteAlarmResponse struct {
		response
		Deleted bool
		Error   error
	}

	createCameraResponse struct {
		response
		Camera model.Camera
		Error  error
	}

	allCameraResponse struct {
		response
		Cameras []model.Camera
		Error   error
	}

	getCameraResponse struct {
		response
		Camera model.Camera
		Error  error
	}

	updateCameraResponse struct {
		response
		Updated bool
		Error   error
	}

	deleteCameraResponse struct {
		response
		Deleted bool
		Error   error
	}
)
