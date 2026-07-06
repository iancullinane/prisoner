> [!NOTE] 
> 🏗️ This is an in-progress refactor that totally changes the original. See `main` for more information but I have changed the default branch to `icullinane/version2` because my previous implementation is not reflective of my modern knowledge of the Go programming language, AWS, and the other technologies I plan or am using.

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

## Database migrations

Migrations are a deliberate, standalone step — the server **never** migrates on its own. This keeps schema changes controlled and safe for a managed database like AWS RDS (no concurrent DDL from multiple replicas, no automatic destructive changes, and the running server doesn't need schema-altering privileges).

The `migrate` command uses [goose](https://github.com/pressly/goose) with the SQL files embedded from `db/migrations`. It reads the connection string from `DATABASE_URL`:

```sh
prisoner migrate status   # show what's applied / pending
prisoner migrate up       # apply all pending migrations
prisoner migrate down     # roll back the last migration
```

The container image doesn't migrate either. Its default command runs the server; migrations are the *same* image run with the command overridden, so they map onto a one-off ECS task, a Kubernetes Job, or a plain `docker run`:

```sh
docker run --rm -e DATABASE_URL=... <image> migrate up   # one-off, before rolling the server
docker run       -e DATABASE_URL=... <image>              # defaults to `server`
```

The intended deploy flow is: run `migrate up` as a one-off task against the database, then roll out the new server image. Because containers no longer self-migrate, a fresh database must be migrated once before the server will function. Locally that's a separate step too — bring up the database with `make start-db`, then apply migrations before starting the server:

```sh
make start-db
DATABASE_URL='postgres://prisoner:prisoner@127.0.0.1:5432/prisoner?sslmode=disable' ./bin/prisoner migrate up
make run-server-postgres
```

## On AI

When I first started this refactor, I resolved to do nothing with AI. I am happy I took this route but as I progressed I found myself honestly too excited by the work to constrain myself to that working method.

To me the largest promise of AI is just how much I can produce (thoughtfully) and I don't want to slow myself down too much. To that effect I have since my initial start allowed more AI into my workflow. I am constraining myself to Claude to learn it more deeply (Similar to how I am constraining myself to the standard library). In addition I am not going too deep and set up a specific `/learning-review` skill to make sure that I am always focused on my original goal of deep learning.
