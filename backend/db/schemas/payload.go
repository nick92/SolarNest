package schemas

type ReadingPayload struct {
	DeviceID string  `json:"device_id"`
	Voltage  float64 `json:"voltage"`
	Current  float64 `json:"current"`
}
