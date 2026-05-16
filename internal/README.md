# store

There is one important design philosophy with this project, which is that I wanted it to be very flexible and be multiple _kinds_ of projects in one. Which meant implements different ways to store player data as that would be the most interesting part.

The initial version of this project did everything in memory. Its typical output from a "round robin" style game would be the following:

```shell
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

This was all simulated however, including the names. The personalities chosen at random and hard coded by me.

## V2

If version 2 was going to be interesting, we were going to need to write things down and expose more than just a command line, we need the internet. So the goal here is to satisfy all of these things as swappable. Step one is this, the store. 

I maintain three types of store, and you can pass a flag depending on which one you want. If you want to use `postgres` you will (obviously) need to start the server. `docker-compose.yml` will get one going for you.

The memory stores are:
- In memory 
- File I/O
- Postgres

In memory is for quick simulations, File I/O was for learning and improving language skill, and postgres is the one that really matters. It is used for the deployed game. 