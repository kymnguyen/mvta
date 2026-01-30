package command

type CreateVehicleCommand struct {
	RefID         string
	VIN           string
	VehicleName   string
	VehicleModel  string
	LicenseNumber string
	Status        string
	Latitude      float64
	Longitude     float64
	Altitude      float64
	Mileage       float64
	FuelLevel     float64
}

func (c *CreateVehicleCommand) CommandName() string {
	return "CreateVehicle"
}
