# prisoner

This is v2 of my personal repo that is both portfolio project (I try to use different features of Go for professional development) and fun side project. 

This both a library to model a [Prisoner's Dilemma](https://en.wikipedia.org/wiki/Prisoner%27s_dilemma), a CLI tool to execute those models as different kinds of game (including group runs like round robin) and different "peronalities" that can play, and a deployment that allows you to play against real humans (or are they?).

## Technologies used

As I said above this is a portfolio, and as such it has mish mash of all my favorite (and sometimes not so favorite) tools.

In no particular order and with a note on why this particular tech:
- [Go](https://go.dev/) - We use this language at [work](https://en.wikipedia.org/wiki/WB_Games_Boston) so it has been the language I have invested the most hours working with 
	- I make a best effort to only use the standard libray
- [AWS](https://aws.amazon.com/) - The next most common thing I do at work is general 👋  "AWS" stuff 👋 so I know a lot about it, most importantly though I know IaC (Mostly CloudFormation and CDK) 
- [Postgres](https://www.postgresql.org/) - Yup, you guessed it we use this at work, though I like it....I guess
	- Data strorage is a requirement of mine and I needed to practice on something of mine so...
	- [SQLC](https://github.com/sqlc-dev/sqlc) - Finally somethine we don't use at work, though the reason I found it is because we basically have this at work but a bespoke version and I wanted something off-the-shelf, I like it 👍
	- [testcontainers](https://golang.testcontainers.org/modules/postgres/) - This is new, part of a general effort to just get better at testing
- [make](https://en.wikipedia.org/wiki/Make_(software)) - We actually don't use this at work which is _wild_ because we have this complicated mess of insanity instead
- [Docker](https://en.wikipedia.org/wiki/Docker_(software)) - You would be crazy not to that 


## Project Structure

I try to follow Go standards here but I am not perfect I am sure I mess things up. In general I think it is pretty clear. Files to look our for are `Makefile` and `sqlc.yaml`. `sqlc` is used to take the db schema and generate the golang ready code. That is why `prisonerdb` is in the `.gitignore` but you can change it if you want the package will still be named `data`.