package query

type GetVehicleQuery struct {
	VehicleID string
}

func (q *GetVehicleQuery) QueryName() string {
	return "GetVehicle"
}
