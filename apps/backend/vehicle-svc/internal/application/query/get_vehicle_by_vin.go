package query

type GetVehicleByVINQuery struct {
	VIN string
}

func (q *GetVehicleByVINQuery) QueryName() string {
	return "GetVehicleByVIN"
}
