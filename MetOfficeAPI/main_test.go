// package main

// import (
// 	"context"
// 	"testing"
// 	"time"
// )

// // Simplified forecast structure
// type SimplifiedForecast struct {
// 	Type     string `json:"type"`
// 	Features []struct {
// 		Properties struct {
// 			ModelRunDate string `json:"modelRunDate"`
// 			TimeSeries   []struct {
// 				Time        string  `json:"time"`
// 				Temperature float64 `json:"screenTemperature"`
// 			} `json:"timeSeries"`
// 		} `json:"properties"`
// 	} `json:"features"`
// }

// // Mock MetOfficeAPI
// type MockMetOfficeAPI struct {
// 	GetForecastFunc func(url string) (*SimplifiedForecast, error)
// }

// func (m *MockMetOfficeAPI) GetForecast(url string) (*SimplifiedForecast, error) {
// 	return m.GetForecastFunc(url)
// }

// func (m *MockMetOfficeAPI) CreateURL(coords []float64) string
// func (m *MockMetOfficeAPI) GetHeaders() http.Header
// func (m *MockMetOfficeAPI) CalculateTotals(data []TimeSeriesData) map[string]*ForecastTotals
// func (m *MockMetOfficeAPI) FindWindows(data []TimeSeriesData) [][]time.Time
// func (m *MockMetOfficeAPI) GetRateLimit() rate.Limit

// // Function to generate test data
// func generateTestForecast(modelRunDate string, timeSeriesCount int) *SimplifiedForecast {
// 	forecast := &SimplifiedForecast{
// 		Type: "FeatureCollection",
// 		Features: []struct {
// 			Properties struct {
// 				ModelRunDate string `json:"modelRunDate"`
// 				TimeSeries   []struct {
// 					Time        string  `json:"time"`
// 					Temperature float64 `json:"screenTemperature"`
// 				} `json:"timeSeries"`
// 			} `json:"properties"`
// 		}{
// 			{
// 				Properties: struct {
// 					ModelRunDate string `json:"modelRunDate"`
// 					TimeSeries   []struct {
// 						Time        string  `json:"time"`
// 						Temperature float64 `json:"screenTemperature"`
// 					} `json:"timeSeries"`
// 				}{
// 					ModelRunDate: modelRunDate,
// 					TimeSeries: make([]struct {
// 						Time        string  `json:"time"`
// 						Temperature float64 `json:"screenTemperature"`
// 					}, timeSeriesCount),
// 				},
// 			},
// 		},
// 	}

// 	baseTime, _ := time.Parse(time.RFC3339, modelRunDate)
// 	for i := 0; i < timeSeriesCount; i++ {
// 		forecast.Features[0].Properties.TimeSeries[i] = struct {
// 			Time        string  `json:"time"`
// 			Temperature float64 `json:"screenTemperature"`
// 		}{
// 			Time:        baseTime.Add(time.Duration(i) * time.Hour).Format(time.RFC3339),
// 			Temperature: float64(15 + i%5), // Simple temperature variation
// 		}
// 	}

// 	return forecast
// }

// // Test function
// func TestUpdateForecasts(t *testing.T) {
// 	// Generate test data
// 	testForecast := generateTestForecast("2024-06-15T15:00Z", 24)

// 	// Create mock API
// 	mockAPI := &MockMetOfficeAPI{
// 		GetForecastFunc: func(url string) (*SimplifiedForecast, error) {
// 			return testForecast, nil
// 		},
// 	}

// 	// Create mock store
// 	mockStore := &MockMetStore{
// 		AddFunc: func(ctx context.Context, cragName string, payload ForecastPayload) error {
// 			// Verify the payload
// 			if payload.LastModelRunTime != testForecast.Features[0].Properties.ModelRunDate {
// 				t.Errorf("Expected LastModelRunTime %v, got %v", testForecast.Features[0].Properties.ModelRunDate, payload.LastModelRunTime)
// 			}
// 			// Add more assertions as needed
// 			return nil
// 		},
// 	}

// 	// Create error channel
// 	errChan := make(chan error, 1)

// 	// Call the function
// 	if err := UpdateForecasts(context.Background(), time.Now().Add(-2*time.Hour), mockAPI, mockStore, errChan); err != nil {
// 		t.Fatal("failed updating", err)
// 	}

// 	// Check for errors
// 	select {
// 	case err := <-errChan:
// 		t.Errorf("Unexpected error: %v", err)
// 	default:
// 		// No error
// 	}

// 	// Add more assertions as needed
// }

// // MockMetStore implementation
// type MockMetStore struct {
// 	AddFunc func(ctx context.Context, cragName string, payload ForecastPayload) error
// }

// func (m *MockMetStore) Add(ctx context.Context, cragName string, payload ForecastPayload) error {
// 	return m.AddFunc(ctx, cragName, payload)
// }

package main
