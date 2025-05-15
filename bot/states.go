package bot

type UserState struct {
	CurrentStep string
	Readings    MeterReadings
}

type MeterReadings struct {
	ColdWater        int
	HotWater         int
	ElectricityDay   int
	ElectricityNight int
}
