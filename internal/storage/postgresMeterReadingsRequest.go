package storage

import "context"

func (s *PostgresStorage) SaveMeterReadings(ctx context.Context, userID int64,
	coldWater, hotWater, electricityDay, electricityNight int) error {

	_, err := s.db.ExecContext(ctx, `
        INSERT INTO meter_readings 
        (user_id, cold_water, hot_water, electricity_day, electricity_night)
        VALUES ($1, $2, $3, $4, $5)
    `, userID, coldWater, hotWater, electricityDay, electricityNight)

	return err
}

func (s *PostgresStorage) GetLastMeterReadings(ctx context.Context, userID int64) (map[string]int, error) {
	var readings struct {
		ColdWater        int
		HotWater         int
		ElectricityDay   int
		ElectricityNight int
	}

	err := s.db.QueryRowContext(ctx, `
        SELECT cold_water, hot_water, electricity_day, electricity_night
        FROM meter_readings
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT 1
    `, userID).Scan(&readings.ColdWater, &readings.HotWater,
		&readings.ElectricityDay, &readings.ElectricityNight)

	if err != nil {
		return nil, err
	}

	return map[string]int{
		"cold_water":        readings.ColdWater,
		"hot_water":         readings.HotWater,
		"electricity_day":   readings.ElectricityDay,
		"electricity_night": readings.ElectricityNight,
	}, nil
}
