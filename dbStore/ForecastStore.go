package dbStore


//in case I forget this will be part of the Forecast service struct, that provides services to the forecast data mode
type SqlForecastStore struct{
	*SqlStore
}

func newSqlForecastStore(sqlStore *SqlStore) ForecastStore{
	fs := &SqlForecastStore{
		SqlStore: sqlStore,
	}
	return fs
}

func (fs SqlForecastStore) StoreForecast(){
	return
}
func (fs SqlForecastStore) GetForecast(){
	return
}
func (fs SqlForecastStore) UpdateForecast(){
	return
}
func (fs SqlForecastStore) DeleteForecast(){
	return
}
func (fs SqlForecastStore) GetForecastByDate(){
	return
}
func (fs SqlForecastStore) GetForecastByDryest(){
	return
}
func (fs SqlForecastStore) GetOldestForecast(){
	return
}

