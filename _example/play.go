package main

import (
	"./sw"
	"bufio"
	"fmt"
	"github.com/onesuper/hobbit"
	"os"
	// "reflect"
	"strings"
)

/////////////////////////////////////////////////////////// Main
func main() {
	// Create terrains
	farmland := hobbit.NewTerrain(sw.Farmland, '.')
	sea := hobbit.NewTerrain(sw.Sea, '~')
	swamp := hobbit.NewTerrain(sw.Swamp, '+')
	forest := hobbit.NewTerrain(sw.Forest, 'o')
	mountain := hobbit.NewTerrain(sw.Mountain, '^')
	plain := hobbit.NewTerrain(sw.Plain, ' ')

	// Create atlas
	atlas, _ := hobbit.NewAtlas(5, 4)
	atlas.SetRegion(0, 0, hobbit.NewRegion(sea))
	atlas.SetRegion(0, 1, hobbit.NewRegion(swamp))
	atlas.SetRegion(0, 2, hobbit.NewRegion(forest))
	atlas.SetRegion(0, 3, hobbit.NewRegion(farmland))

	atlas.SetRegion(1, 0, hobbit.NewRegion(mountain))
	atlas.SetRegion(1, 1, hobbit.NewRegion(plain))
	atlas.SetRegion(1, 2, hobbit.NewRegion(sea))
	atlas.SetRegion(1, 3, hobbit.NewRegion(plain))

	atlas.SetRegion(2, 0, hobbit.NewRegion(forest))
	atlas.SetRegion(2, 1, hobbit.NewRegion(mountain))
	atlas.SetRegion(2, 2, hobbit.NewRegion(swamp))
	atlas.SetRegion(2, 3, hobbit.NewRegion(farmland))

	atlas.SetRegion(3, 0, hobbit.NewRegion(swamp))
	atlas.SetRegion(3, 1, hobbit.NewRegion(farmland))
	atlas.SetRegion(3, 2, hobbit.NewRegion(mountain))

	atlas.SetRegion(4, 0, hobbit.NewRegion(plain))
	atlas.SetRegion(4, 1, hobbit.NewRegion(forest))
	atlas.SetRegion(4, 2, hobbit.NewRegion(sea))

	// Create races
	humans := sw.NewHumans()
	orcs := sw.NewOrcs()

	// Create skills
	foresting := &sw.Foresting{}

	races := []hobbit.RaceI{}
	races = append(races, humans)
	races = append(races, orcs)

	skills := []hobbit.RaceI{}
	skills = append(skills, foresting)

	// Dealing the interaction part
	cmd := hobbit.NewCmd()
	fmt.Println(bannerLarge)
	for round := 1; round <= 10; round++ {
		for _, race := range races {
			race.GatherSoldiers(atlas)
			Screen(atlas, races)
			/////////////////////////////////////////////////////////// Recalling Stage
			if race.OccupiedRegions(atlas) > 0 {
				cmd.Banner(fmt.Sprintf(" Round %d | %s | %s ", round, race.GetName(), "Recalling Stage"), '.')
				cmd.Promptln("Choose a region to recall soldier from. " +
					"Insert a coordinate like: `1-2`.\n" +
					"Or insert `f` to finish this stage.")
				for {
					command, _ := cmd.ReadCommand()
					if command == "f" {
						break
					}
					if row, col, err := hobbit.ParseCoord(command); err == nil {
						if region, err1 := atlas.GetRegion(row, col); err1 == nil {
							if foul := race.RecallFromOccupied(region); foul == nil {
								fmt.Printf("%s recall the soldier from region %d-%d.\n", race.GetName(), row, col)
								Screen(atlas, races)
							} else {
								cmd.Promptln(foul.Error())
							}
						} else {
							cmd.Promptln(err1.Error())
						}
					} else {
						cmd.Promptln(err.Error())
					}
				}
			}
			/////////////////////////////////////////////////////////// Conquering Stage
			cmd.Banner(fmt.Sprintf(" Round %d | %s | %s ", round, race.GetName(), "Conquering Stage"), '.')
			cmd.Promptln("Choose a region to conquer. " +
				"Insert a coordinate like: `1-2`.\n" +
				"Or insert `f` to finish this stage.")
			for {
				command, _ := cmd.ReadCommand()
				if command == "f" {
					break
				}
				if row, col, err := hobbit.ParseCoord(command); err == nil {
					if region, err1 := atlas.GetRegion(row, col); err1 == nil {
						if !race.CanReach(atlas, row, col) {
							cmd.Promptln("can not reach this region!")
							continue
						}
						if foul := race.Conquer(region); foul == nil {
							fmt.Printf("%s conquers region %d-%d.\n", race.GetName(), row, col)
							Screen(atlas, races)
						} else {
							cmd.Promptln(foul.Error())
						}
						if race.GetSoldiers() == 0 {
							break
						}
					} else {
						cmd.Promptln(err1.Error())
					}
				} else {
					cmd.Promptln(err.Error())
				}
			}
			/////////////////////////////////////////////////////////// Redeploying Stage
			cmd.Banner(fmt.Sprintf(" Round %d | %s | %s ", round, race.GetName(), "Redeploying Stage"), '.')
			cmd.Promptln("For redeploying your troops between regions.\n" +
				"Insert a coordinate pair like `3-0 2-1`.\n" +
				"If you have idle soldiers, insert a single coordinate to deploy it.\n" +
				"Or insert `f` to finish this stage.")
			for {
				command, _ := cmd.ReadCommand()
				if command == "f" {
					break
				}
				fields := strings.Fields(command)
				if len(fields) == 1 {
					if row, col, err := hobbit.ParseCoord(command); err == nil {
						if region, err1 := atlas.GetRegion(row, col); err1 == nil {
							if foul := race.RedeployTo(region); foul == nil {
								fmt.Printf("%s redeploys a idle soldier to %d-%d\n",
									race.GetName(), row, col)
								Screen(atlas, races)
							} else {
								cmd.Promptln(foul.Error())
							}
						} else {
							cmd.Promptln(err1.Error())
						}
					} else {
						cmd.Promptln(err.Error())
					}
				} else if len(fields) == 2 {
					if rSrc, qSrc, err := hobbit.ParseCoord(fields[0]); err == nil {
						if rDst, qDst, err1 := hobbit.ParseCoord(fields[1]); err1 == nil {
							if src, err2 := atlas.GetRegion(rSrc, qSrc); err2 == nil {
								if dst, err3 := atlas.GetRegion(rDst, qDst); err3 == nil {
									if foul := race.RedeployBetween(src, dst); foul == nil {
										fmt.Printf("%s redeploys a soldier from %d-%d to %d-%d\n",
											race.GetName(), rSrc, qSrc, rDst, qDst)
										Screen(atlas, races)
									} else {
										cmd.Promptln(foul.Error())
									}
								} else {
									cmd.Promptln("2nd field: " + err3.Error())
								}
							} else {
								cmd.Promptln("1st field: " + err2.Error())
							}
						} else {
							cmd.Promptln("2nd field: " + err1.Error())
						}
					} else {
						cmd.Promptln("1st field: " + err.Error())
					}
				} else {
					cmd.Promptln("the command must have 2 fields")
				}
			}
			/////////////////////////////////////////////////////////// Scoring Stage
			coins := race.Score(atlas)
			fmt.Printf("%s make %d victory coins\n", race.GetName(), coins)
		}

	}
}

const bannerLarge = `+----------------------------------------------------------------------------+
|   _________              .__  .__     __      __            .__       .___ |
|  /   _____/ _____ _____  |  | |  |   /  \    /  \___________|  |    __| _/ |
|  \_____  \ /     \\__  \ |  | |  |   \   \/\/   /  _ \_  __ \  |   / __ |  |
|  /        \  Y Y  \/ __ \|  |_|  |__  \        (  <_> )  | \/  |__/ /_/ |  |
| /_______  /__|_|  (____  /____/____/   \__/\  / \____/|__|  |____/\____ |  |
|         \/      \/     \/                   \/                         \/  |
|                                V1.0 alpha                                  |
+----------------------------------------------------------------------------+`

/////////////////////////////////////////////////////////// Screen
func Screen(atlas *hobbit.Atlas, races []hobbit.RaceI) {

	asciiString := bannerLarge + "\n"
	hexString := atlas.GetString()
	hexlines := strings.Split(hexString, "\n")

	// Combines the all race cards to one stack
	raceString := ""
	for _, r := range races {
		raceString += raceCard(r)
	}
	racelines := strings.Split(raceString, "\n")

	// Cat the race card and the hex board together
	for i := range hexlines {
		asciiString += hexlines[i]
		if i < len(racelines) {
			asciiString += racelines[i]
		}
		asciiString += "\n"
	}
	// Output to file
	file, err := os.Create("map.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Seek(0, 0)
	writer := bufio.NewWriter(file)
	writer.WriteString(asciiString)
	writer.Flush()
}

func raceCard(r hobbit.RaceI) string {
	card := "|  XXXXXXXXXXX  |\n" +
		"|  YYYYYYYYYYY  |\n" +
		"+---------------+\n"
	line1 := hobbit.FixToLength(r.GetName(), 11, ' ')
	line2 := hobbit.SeveralSymbols(r.GetSymbol(), r.GetSoldiers())
	line2 = hobbit.FixToLength(line2, 11, ' ')
	card = strings.Replace(card, "XXXXXXXXXXX", line1, 1)
	card = strings.Replace(card, "YYYYYYYYYYY", line2, 1)
	return card
}
