package models

import ()

type Mission struct {
	Filename string `json:"filename"`
}

type ChangeMissionBody struct {
	MissionName string `json:"mission_name"`
}
