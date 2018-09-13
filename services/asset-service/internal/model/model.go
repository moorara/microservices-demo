package model

type (
	// Asset is the supertype for all assets
	Asset struct {
		ID       string `json:"id,omitempty" gorm:"primary_key"`
		SiteID   string `json:"siteId" gorm:"not null"`
		SerialNo string `json:"serialNo" gorm:"not null"`
	}

	// AssetInput is used for creating/updating an asset
	AssetInput struct {
		SiteID   string
		SerialNo string
	}

	// Alarm is an imaginary alarm!
	Alarm struct {
		Asset
		Material string `json:"material"`
	}

	// AlarmInput is used for creating/updating an alarm
	AlarmInput struct {
		AssetInput
		Material string
	}

	// Camera is an imaginary camera!
	Camera struct {
		Asset
		Resolution int `json:"resolution"`
	}

	// CameraInput is used for creating/updating a camera
	CameraInput struct {
		AssetInput
		Resolution int
	}
)
