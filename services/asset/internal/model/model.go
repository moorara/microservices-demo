package model

type (
	// Asset is the supertype for all assets
	Asset struct {
		ID       string `json:"id" gorm:"primary_key"`
		SiteID   string `json:"siteId" gorm:"not null"`
		SerialNo string `json:"serialNo" gorm:"not null"`
	}

	// AssetInput is used for creating/updating an asset
	AssetInput struct {
		SiteID   string `json:"siteId"`
		SerialNo string `json:"serialNo"`
	}

	// Alarm is an imaginary alarm!
	Alarm struct {
		Asset
		Material string `json:"material"`
	}

	// AlarmInput is used for creating/updating an alarm
	AlarmInput struct {
		AssetInput
		Material string `json:"material"`
	}

	// Camera is an imaginary camera!
	Camera struct {
		Asset
		Resolution int `json:"resolution"`
	}

	// CameraInput is used for creating/updating a camera
	CameraInput struct {
		AssetInput
		Resolution int `json:"resolution"`
	}
)
