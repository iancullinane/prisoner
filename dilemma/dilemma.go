package dilemma

import "github.com/iancullinane/prisoner/entity"

// Once will play the prisoners dilemma one time
func Once(player1, player2 *entity.Entity) {

	moveA := player1.Play()
	moveB := player2.Play()

	a, b := Compute(moveA, moveB)

	player1.RecordMoves(moveA, moveB)
	player2.RecordMoves(moveB, moveA)

	player1.Score += a
	player2.Score += b

}

// Select the number of rounds two entities should engage in 'the prisoners
// dilemma' and update their score
func PlayRepeated(player1, player2 *entity.Entity, rounds int) {

	for i := 1; i <= rounds; i++ {
		Once(player1, player2)
	}

}

func PlayOneTournament(agents []*entity.Entity, turns int) {

	for i := 0; i < len(agents); i++ {
		agents[i].Reset()
	}

	// Round robin
	for i := 0; i < len(agents); i++ {
		playerA := agents[i]

		for j := i + 1; j < len(agents); j++ {
			playerB := agents[j]
			PlayRepeated(playerA, playerB, 10)
		}
	}

}

// PD.playOneTournament = function(agents, turns){

// 	// Round robin!
// 	for(var i=0; i<agents.length; i++){
// 		var playerA = agents[i];
// 		for(var j=i+1; j<agents.length; j++){
// 			var playerB = agents[j];
// 			PD.playRepeatedGame(playerA, playerB, turns);
// 		}
// 	}

// };
