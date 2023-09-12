package models

import ()

// {"missionlist":{"listLoop":false,"listShuffle":false,"listStartIndex":1,"missionList":["C:\\users\\dcs\\Saved Games\\DCS.openbeta_server\\Missions\\Training-Marianas GMT aligned-Rain.miz"],"missionTheatres":["MarianaIslands"]}}
type MissionListResult struct {
  MissionList MissionListContent `json:"missionlist"` 
}

type MissionListContent struct {
  Loop bool `json:"listLoop"`
  Shuffle bool `json:"listShuffle"`
  StartIndex int `json:"listStartIndex"`
  Missions []string `json:"missionList"`
  MissionTheatres []string `json:"missionTheatres"`
}


type Mission struct {
  Index int `json:"index"`
	Filename string `json:"filename"`
}

type ChangeMissionBody struct {
	MissionIndex int `json:"mission_index"`
}
