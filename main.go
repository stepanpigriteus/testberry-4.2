package main

import "disgreps/utils"

func main() {
	cfg := utils.FlagParser()
	lines := utils.ReadInput(cfg)
	resultLines := utils.Proccessor(cfg, lines)
	utils.OutRes(resultLines, cfg)

	
}
