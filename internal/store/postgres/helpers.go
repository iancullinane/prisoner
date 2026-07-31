package postgres

import (
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
	sqlcdb "github.com/iancullinane/prisoner/prisonerdb"
)

// interaction from row
func interactionFromRow(i sqlcdb.Interaction) types.Interaction {
	return types.Interaction{
		ID:          i.ID,
		PlayerA:     i.PlayerAID,
		PlayerB:     i.PlayerBID,
		PlayerAMove: prisoner.Move([]rune(i.PlayerAMove)[0]),
		PlayerBMove: prisoner.Move([]rune(i.PlayerBMove)[0]),
		PlayedAt:    i.PlayedAt.Time,
	}
}

func prettyFromRow(r sqlcdb.GetPrettyHistoryRow) types.PrettyInteraction {
	return types.PrettyInteraction{
		PlayerAName: r.PlayerAName,
		PlayerBName: r.PlayerBName,
		PlayerAMove: prisoner.Move([]rune(r.PlayerAMove)[0]),
		PlayerBMove: prisoner.Move([]rune(r.PlayerBMove)[0]),
		PlayedAt:    r.PlayedAt.Time,
	}
}
