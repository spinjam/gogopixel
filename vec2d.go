package main

type Vec2D[T int | float64] struct {
	X, Y T
}

func (v *Vec2D[T]) Add(other Vec2D[T]) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vec2D[T]) AddDelta(delta T) {
	v.X += delta
	v.Y += delta
}

/*

 */
