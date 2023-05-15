package main

import (
	"fmt"
	"pokeapis"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var cliCommands = map[string]cliCommand{}

func commandHelp() error {
	fmt.Println("Welcome to the pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Printf("help:%v\n", cliCommands["help"].description)
	fmt.Printf("exit:%v\n", cliCommands["exit"].description)
	fmt.Printf("map:%s\n", cliCommands["map"].description)
	fmt.Printf("mapb:%s\n", cliCommands["mapb"].description)
	fmt.Println("")
	return nil
}
func exitHelp() error {
	return nil
}

func updateCli() map[string]cliCommand {
	data := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the poke dex cli",
			callback:    exitHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 area locations",
			callback:    displayLocations,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    displayBackwardLocations,
		},
	}
	return data
}

func displayLocations() error {
	err := pokeapis.DisplayLocations()
	if err != nil {
		return err
	}
	return nil
}

func displayBackwardLocations() error {
	err := pokeapis.DisplayBackwardLocations()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	//pokecache.NewCache(5 * time.Second)
	cli := "pokedex >"
	var command string = ""
	cliCommands = updateCli()
	for {
		fmt.Printf("%s", cli)
		_, _ = fmt.Scanln(&command)
		switch command {
		case "help":
			err := cliCommands["help"].callback()
			if err != nil {
				fmt.Println(err)
				return
			}
		case "exit":
			err := cliCommands["exit"].callback()
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		case "map":
			err := cliCommands["map"].callback()
			if err != nil {
				fmt.Println(err)
				return
			}
		case "mapb":
			err := cliCommands["mapb"].callback()
			if err != nil {
				fmt.Println(err)
				return
			}
		default:
			fmt.Println("Please provide help or exit to find usage")
		}
	}

}
