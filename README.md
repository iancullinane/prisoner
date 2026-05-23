> [!NOTE] 
> This is super old and out-of-date but I am overhauling from the ground up and doing grug things like writing tests. See [here]([https://github.com/iancullinane/prisoner/tree/icullinane/add-sqlc-testcontainer](https://github.com/iancullinane/prisoner/tree/icullinane/version2)) until I merge it in.

prisoner
========

Implementation of the prisoners dilemma in go. 

## As a package

```
package main

import (
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

func main() {
	game := prisoner.New()
	players := game.GetPersonalities(10)
	game.PlayTournament(players, 5)
	game.PrintResults(players)
}
```

## As a command

```sh
go get -u -v github.com/iancullinane/cmd/prisoner/...

prisoner -players <number of players> -rounds <number of rounds>

2021/04/06 12:02:05 Name        Score   Personality
2021/04/06 12:02:05 ------------------------------------
2021/04/06 12:02:05 Aubrey      70      random
2021/04/06 12:02:05 Daniel      81      random
2021/04/06 12:02:05 Ella        73      copycat
2021/04/06 12:02:05 Avery       97      random
2021/04/06 12:02:05 Elijah      66      niceguy
2021/04/06 12:02:05 Jayden      105     cheater
2021/04/06 12:02:05 Olivia      102     revenge
2021/04/06 12:02:05 James       94      copycat
2021/04/06 12:02:05 Ella        103     cheater
2021/04/06 12:02:05 Isabella    98      copycat

```
