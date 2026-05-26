package db

import "embed"

//go:embed schema.sql
var SchemaSQL string

//go:embed migrations/*.sql
var Migrations embed.FS
