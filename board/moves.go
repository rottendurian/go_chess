package board

import (
	"chess/pieces"
	"chess/utils"
	"errors"
	// "fmt"
)

type Vec2 = utils.Vec2

func (b Board) validPosition(pos Vec2) bool {
	return pos.X >= 0 && pos.X < BOARD_SIZE && pos.Y >= 0 && pos.Y < BOARD_SIZE
}

func (b Board) GetPiece(pos Vec2) (pieces.Piece, error) {
	if b.validPosition(pos) {
		return b.Nodes[pos.Y][pos.X], nil
	}
	return pieces.NewNone(), errors.New("GetPiece called on an invalid location")
}

func (b Board) hasCollision(piecePos Vec2, destPos Vec2) (bool, error) {
	delta := utils.GetDelta(piecePos, destPos)

	if delta.X == 0 && delta.Y == 0 {
		return true, errors.New("hasCollision called with the same position")
	}
	
	destPiece, err := b.GetPiece(destPos)
	if err != nil {
		return true, err
	}

	originalPiece, err := b.GetPiece(piecePos)
	if err != nil {
		return true, err
	}
	
	// check if the color of the piece is the same and the piece is not a none piece // early return
	if destPiece.Color == originalPiece.Color && destPiece.PieceType != pieces.NONE {
		return true, nil
	}

	// fmt.Println("Delta is: ", delta)
	
	// assert that move is diagonal or straight
	if (utils.Abs(delta.X) == utils.Abs(delta.Y)) || (delta.X == 0 && delta.Y != 0) || (delta.X != 0 && delta.Y == 0) {

		// check each position in the direction of the move
		// clamp the vectors from -1 to 1
		xOffset := utils.Clamp(delta.X, -1, 1)
		yOffset := utils.Clamp(delta.Y, -1, 1)

		searchPos := piecePos
		
		for searchPos.X != destPos.X || searchPos.Y != destPos.Y {
			searchPos.X += xOffset
			searchPos.Y += yOffset

			piece, err := b.GetPiece(searchPos)
			if err != nil {
				return true, err
			}

			if piece.PieceType != pieces.NONE {
				return true, nil
			}
		}
	}
	
	return false, nil
}

func (b Board) ValidMove(piecePos Vec2, destPos Vec2) bool {
	piece, err := b.GetPiece(piecePos)
	if err != nil {
		return false
	}

	if piece.PieceType == pieces.NONE {
		return false
	}

	
	if piece.ValidMove(utils.GetDelta(piecePos, destPos)) {
		result, err := b.hasCollision(piecePos, destPos)
		
		if err != nil {
			panic(err)
		}
		// fmt.Println("Got past collision check with result: ", result)
		if !result {
			return true
		}
		// fmt.Println("Got past initial checks")
	}

	return false
}

func (b Board) MovePiece(boardPos Vec2, destPos Vec2) (Board, error) {
	if b.ValidMove(boardPos, destPos) {
		piece, err := b.GetPiece(boardPos)
		if err != nil {
			return b, err
		}

		if piece.PieceType == pieces.PAWN {
			piece.FirstMove = false
		}

		b.Nodes[destPos.Y][destPos.X] = piece
		b.Nodes[boardPos.Y][boardPos.X] = pieces.NewNone()

		return b, nil
	}
	return b, errors.New("invalid Move")
}
