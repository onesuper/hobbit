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
	farmland := hobbit.NewTerrain(sw.Farmland, '+')
	sea := hobbit.NewTerrain(sw.Sea, '~')
	swamp := hobbit.NewTerrain(sw.Swampland, '.')
	forest := hobbit.NewTerrain(sw.Forestland, 'o')
	mountain := hobbit.NewTerrain(sw.Mountain, '^')
	plain := hobbit.NewTerrain(sw.Flatland, ' ')

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

	// Race list
	var races []hobbit.RaceI
	races = append(races,
		sw.NewHumans(),
		sw.NewElves(),
		sw.NewOrcs(),
		sw.NewSkeletons(),
		sw.NewMermans(),
	)

	// Skill list
	var skills []hobbit.SkillI
	skills = append(skills,
		sw.NewFlying(),
		sw.NewSwamp(),
		sw.NewAlchemist(),
		sw.NewMerchant(),
		sw.NewForest(),
	)

	var totalCoins []int

	for i, race := range races {
		race.AddSoldiers(skills[i].GetSoldiers())
		totalCoins = append(totalCoins, 0)
	}

	// Dealing the interaction part
	cmd := hobbit.NewCmd()
	fmt.Println(bannerLarge)
	for round := 1; round <= 6; round++ {
		for i, race := range races {
			/////////////////////////////////////////////////////////// Recalling Stage
			race.GatherSoldiers(atlas)
			Screen(atlas, races, skills)
			if race.OccupiedRegions(atlas) > 0 {
				cmd.Banner(fmt.Sprintf(" Round %d | %s | %s ", round, race.GetName(), "Recalling Stage"), '.')
				cmd.Promptln("Choose a region to recall soldier from. Insert a coordinate like: `1-2`.\n" +
					"Or insert `f` to finish this stage.")
				for {
					command, _ := cmd.ReadCommand()
					if command == "f" {
						break
					}
					row, col, err := hobbit.ParseCoord(command)
					if err != nil {
						cmd.Promptln(err.Error())
						continue
					}
					if foul := race.RecallFrom(atlas, row, col, 0); foul != nil {
						cmd.Promptln(foul.Error())
						continue
					}
					fmt.Printf("%s recalls 1 soldier from region %d-%d.\n", race.GetName(), row, col)
					Screen(atlas, races, skills)
				} // for
			}
			/////////////////////////////////////////////////////////// Conquering Stage
			Screen(atlas, races, skills)
			cmd.Banner(fmt.Sprintf(" Round %d | %s | %s ", round, race.GetName(), "Conquering Stage"), '.')
			cmd.Promptln("Choose a region to conquer. Insert a coordinate like: `1-2`.\n" +
				"Or insert `f` to finish this stage.")
			for {
				command, _ := cmd.ReadCommand()
				if command == "f" {
					break
				}
				row, col, err1 := hobbit.ParseCoord(command)
				if err1 != nil {
					cmd.Promptln(err1.Error())
					continue
				}
				region, err2 := atlas.GetRegion(row, col)
				if err2 != nil {
					cmd.Promptln(err2.Error())
					continue
				}
				if !race.CanReach(atlas, row, col) && !skills[i].CanReach(atlas, row, col) {
					cmd.Promptln("can not reach this region!")
					continue
				}
				if race.HasOccupied(atlas, row, col) {
					cmd.Promptln("can not attack your own region!")
					continue
				}
				defense := race.GetDefenseOver(atlas, row, col)
				if race.GetSoldiers() < defense {
					cmd.Promptln("not enough soldiers!")
					continue
				}
				if troop := region.GetTroop(); troop != nil {
					// Find the to defeat
					for _, r := range races {
						if r.GetSymbol() == troop.Symbol {
							r.Defeat(troop.Soldiers)
							break
						}
					}
					race.AfterEachDefeat()
				}
				race.Reside(atlas, row, col, defense)
				fmt.Printf("%s conquers region %d-%d.\n", race.GetName(), row, col)
				Screen(atlas, races, skills)
				if race.GetSoldiers() == 0 {
					break // Running out of soldiers
				}
			} // for
			/////////////////////////////////////////////////////////// Redeploying Stage
			race.AfterConquest()
			Screen(atlas, races, skills)
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
				if fields := strings.Fields(command); len(fields) == 1 {
					row, col, err := hobbit.ParseCoord(command)
					if err != nil {
						cmd.Promptln(err.Error())
						continue
					}
					if foul := race.DeployTo(atlas, row, col); foul != nil {
						cmd.Promptln(foul.Error())
						continue
					}
					fmt.Printf("%s redeploys a idle soldier to %d-%d\n", race.GetName(), row, col)
					Screen(atlas, races, skills)
				} else if len(fields) == 2 {
					rSrc, qSrc, err1 := hobbit.ParseCoord(fields[0])
					if err1 != nil {
						cmd.Promptln("1st field: " + err1.Error())
						continue
					}
					rDst, qDst, err2 := hobbit.ParseCoord(fields[1])
					if err2 != nil {
						cmd.Promptln("2nd field: " + err2.Error())
						continue
					}
					// In redeploying stage, the race must leave at least 1 soldier on each region.
					if foul := race.RecallFrom(atlas, rSrc, qSrc, 1); foul != nil {
						cmd.Promptln(foul.Error())
						continue
					}
					if foul := race.DeployTo(atlas, rDst, qDst); foul != nil {
						cmd.Promptln(foul.Error())
						continue
					}
					fmt.Printf("%s redeploys a soldier from %d-%d to %d-%d\n", race.GetName(), rSrc, qSrc, rDst, qDst)
					Screen(atlas, races, skills)
				} else {
					cmd.Promptln("the command must have 2 fields")
				}
			} // for
			/////////////////////////////////////////////////////////// Scoring Stage
			coins := races[i].Score(atlas)
			coins += skills[i].Score(atlas, race)
			fmt.Printf("%s make %d victory coins\n", race.GetName(), coins)
			totalCoins[i] += coins
		} // races
	} // round
	cmd.Banner(" Scoreboard ", '.')
	for i, race := range races {
		fmt.Printf("%s: %d\n", skills[i].GetName()+" "+race.GetName(), totalCoins[i])
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
func Screen(atlas *hobbit.Atlas, races []hobbit.RaceI, skills []hobbit.SkillI) {

	asciiString := bannerLarge + "\n"
	hexString := atlas.GetString()
	hexlines := strings.Split(hexString, "\n")

	// Combines the all race cards to one stack
	raceString := ""
	for i, _ := range races {
		raceString += raceCard(races[i], skills[i])
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

func raceCard(r hobbit.RaceI, s hobbit.SkillI) string {
	card :=
		"|  XXXXXXXXXXXXXXXXXXX  |\n" +
			"|  YYYYYYYYYYYYYYYYYYY  |\n" +
			"+-----------------------+\n"
	line1 := hobbit.FixToLength(s.GetName()+" "+r.GetName(), 19, ' ')
	line2 := hobbit.SeveralSymbols(r.GetSymbol(), r.GetSoldiers())
	line2 = hobbit.FixToLength(line2, 19, ' ')
	card = strings.Replace(card, "XXXXXXXXXXXXXXXXXXX", line1, 1)
	card = strings.Replace(card, "YYYYYYYYYYYYYYYYYYY", line2, 1)
	return card
}
