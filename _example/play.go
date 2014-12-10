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
	humans := &sw.Humans{hobbit.Race{"Humans", 'H', 10}}
	orcs := &sw.Orcs{hobbit.Race{"Orcs", 'O', 10}}
	mermans := &sw.Mermans{hobbit.Race{"Mermans", 'M', 10}}

	// Create skills
	foresting := &sw.Foresting{}

	races := []hobbit.RaceI{}
	races = append(races, mermans, humans, orcs)

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
					row, col, err := hobbit.ParseCoord(command)
					if err != nil {
						cmd.Promptln(err.Error())
						continue
					}
					foul := race.RecallFrom(atlas, row, col)
					if foul != nil {
						cmd.Promptln(foul.Error())
						continue
					}
					fmt.Printf("%s recall the soldier from region %d-%d.\n", race.GetName(), row, col)
					Screen(atlas, races)
				} // for
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
				row, col, err := hobbit.ParseCoord(command)
				if err != nil {
					cmd.Promptln(err.Error())
					continue
				}
				// Check the reachability before conquering.
				if _, err := race.CanReach(atlas, row, col); err != nil {
					cmd.Promptln(err.Error())
					continue
				}
				// Check whether the race has enough idle soldiers to conquer
				// the region.
				defense := race.GetDefenseOver(atlas, row, col)
				if race.GetSoldiers() < defense {
					cmd.Promptln("not enough soldiers!")
					continue
				}
				// Defeat the troop on the region.
				region, _ := atlas.GetRegion(row, col)
				if troop := region.GetTroop(); troop != nil {
					if troop.Race == race {
						cmd.Promptln("can not attack own region!")
						continue
					} else if troop.Race != nil {
						troop.Race.Defeat(troop.Soldiers)
					}
				}
				// The new race resides soliders on the region
				race.Reside(atlas, row, col, defense)
				fmt.Printf("%s conquers region %d-%d.\n", race.GetName(), row, col)
				Screen(atlas, races)
				if race.GetSoldiers() == 0 {
					break // Running out of soldiers
				}
			} // for
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
				if fields := strings.Fields(command); len(fields) == 1 {
					row, col, err := hobbit.ParseCoord(command)
					if err != nil {
						cmd.Promptln(err.Error())
						continue
					}
					foul := race.RedeployTo(atlas, row, col)
					if foul != nil {
						cmd.Promptln(foul.Error())
						continue
					}
					fmt.Printf("%s redeploys a idle soldier to %d-%d\n", race.GetName(), row, col)
					Screen(atlas, races)
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
					foul := race.RedeployBetween(atlas, rSrc, qSrc, rDst, qDst)
					if foul != nil {
						cmd.Promptln(foul.Error())
						continue
					}
					fmt.Printf("%s redeploys a soldier from %d-%d to %d-%d\n",
						race.GetName(), rSrc, qSrc, rDst, qDst)
					Screen(atlas, races)
				} else {
					cmd.Promptln("the command must have 2 fields")
				}
			} // for
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
