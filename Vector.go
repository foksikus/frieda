package main

type Vector struct {
	X, Y int
}

// returns new vector
func NewVector(x, y int) *Vector {
	return &Vector{X: x, Y: y}
}

type VectorPool struct {
	width, height int
	pool          []*Vector
}

// return new vector pool
func NewVectorPool(width, height int) *VectorPool {
	size := width * height
	pool := make([]*Vector, size)
	for i := range pool {
		x, y := calcXY(i, width)
		pool[i] = NewVector(x, y)
	}
	return &VectorPool{width, height, pool}
}

func calcXY(index, width int) (x, y int) {
	return index % width, index / width
}

// return vector from pool with x ,y
func (pool *VectorPool) Get(x, y int) *Vector {
	return pool.pool[y*pool.width+x]
}
