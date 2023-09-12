package models

import ()

type MissionListContent struct {
  Loop bool `json:"listLoop"`
  Shuffle bool `json:"listShuffle"`
  StartIndex int `json:"listStartIndex"`
  Missions []string `json:"missionList"`
}

type MissionListResult struct {
  MissionList MissionListContent `json:"missionlist"` 
}

type Mission struct {
  Index int `json:"index"`
	Filename string `json:"filename"`
}

type ChangeMissionBody struct {
	MissionIndex int `json:"mission_index"`
}
