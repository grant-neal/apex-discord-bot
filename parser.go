package main

import "github.com/tidwall/gjson"

func getMapRotation(mapRotation string, gameMode string) map[string]string {
	mapInfo := make(map[string]string)

	mapInfo["currentMap"] = gjson.Get(mapRotation, gameMode+".current.map").String()
	if gameMode != "ranked" {
		mapInfo["remainingTimer"] = gjson.Get(mapRotation, gameMode+".current.remainingTimer").String()
	}
	mapInfo["nextMap"] = gjson.Get(mapRotation, gameMode+".next.map").String()

	return mapInfo
}

func getCharacterInfo(characterInfo string) map[string]string {
	character := make(map[string]string)

	character["Name"] = gjson.Get(characterInfo, "global.name").String()
	character["RankScore"] = gjson.Get(characterInfo, "global.rank.rankScore").String()
	character["rankName"] = gjson.Get(characterInfo, "global.rank.rankName").String()
	character["rankDiv"] = gjson.Get(characterInfo, "global.rank.rankDiv").String()

	// TODO: We need to look into downloading this file so we can manipulate it and send it back as part of the embedded response.
	character["rankImage"] = gjson.Get(characterInfo, "global.rank.rankImg").String()

	return character
}
