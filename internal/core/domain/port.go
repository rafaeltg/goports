package domain

import (
	"encoding/json"
)

type (
	// Port is a struct representing the data structure of each port.
	Port struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		City        string    `json:"city"`
		Country     string    `json:"country"`
		Alias       []string  `json:"alias,omitempty"`
		Regions     []string  `json:"regions,omitempty"`
		Coordinates []float64 `json:"coordinates,omitempty"`
		Province    string    `json:"province"`
		Timezone    string    `json:"timezone"`
		Unlocs      []string  `json:"unlocs,omitempty"`
		Code        string    `json:"code"`
	}

	Ports []Port
)

func (p *Port) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Port) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}

	return nil
}
