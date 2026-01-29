package query

type GetVehicleChangeHistoryQuery struct {
	VehicleID string
	Limit     int
	Offset    int
}

func (q *GetVehicleChangeHistoryQuery) QueryName() string {
	return "GetVehicleChangeHistory"
}
