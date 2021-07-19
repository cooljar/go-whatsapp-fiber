package domain

import "encoding/json"

type WaGroup struct {
	ID string `json:"id"`
	Owner string `json:"owner"`
	Subject string `json:"subject"`
	Creation int `json:"creation"`
	SubjectTime int `json:"subjectTime"`
	SubjectOwner string `json:"subjectOwner"`
	Desc string `json:"desc"`
	DescId string `json:"descId"`
	DescTime int `json:"descTime"`
	DescOwner string `json:"descOwner"`
	Participants []WaGroupParticipants `json:"participants"`
}

type WaGroupParticipants struct {
	ID string `json:"id"`
	IsAdmin bool `json:"isAdmin"`
	IsSuperAdmin bool `json:"isSuperAdmin"`
}

// FromJSON decode json to book struct
func (w *WaGroup) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, w)
}

// ToJSON encode book struct to json
func (w *WaGroup) ToJSON() []byte {
	str, _ := json.Marshal(w)
	return str
}
