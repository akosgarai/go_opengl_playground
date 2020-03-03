package primitives

import (
	"testing"
)

func TestNewCubeByVectorAndLength(t *testing.T) {
	testData := []struct {
		V Vector
		L float64
	}{
		{Vector{0, 0, 0}, 1.0},
		{Vector{0, 0, 0}, 0.5},
		{Vector{1, 1, 1}, 0.5},
	}
	for _, tt := range testData {
		cube := NewCubeByVectorAndLength(tt.V, tt.L)
		// Front left bottom (H)
		if cube.H.Coordinate.X != tt.V.X || cube.H.Coordinate.Y != tt.V.Y || cube.H.Coordinate.Z != tt.V.Z+tt.L {
			t.Log(cube.H.Coordinate)
			t.Log(tt.V)
			t.Error("Invalid 'H' Point - NewCubeByVectorAndLength")
		}
		// Front right bottom (G)
		if cube.G.Coordinate.X != tt.V.X+tt.L || cube.G.Coordinate.Y != tt.V.Y || cube.G.Coordinate.Z != tt.V.Z+tt.L {
			t.Log(cube.G.Coordinate)
			t.Log(tt.V)
			t.Error("Invalid 'G' Point - NewCubeByVectorAndLength")
		}
		// Front left top (E)
		if cube.E.Coordinate.X != tt.V.X || cube.E.Coordinate.Y != tt.V.Y+tt.L || cube.E.Coordinate.Z != tt.V.Z+tt.L {
			t.Log(cube.E.Coordinate)
			t.Log(tt.V)
			t.Error("Invalid 'E' Point - NewCubeByVectorAndLength")
		}
		// Front right top (F)
		if cube.F.Coordinate.X != tt.V.X+tt.L || cube.F.Coordinate.Y != tt.V.Y+tt.L || cube.F.Coordinate.Z != tt.V.Z+tt.L {
			t.Log(cube.F.Coordinate)
			t.Log(tt.V)
			t.Error("Invalid 'F' Point - NewCubeByVectorAndLength")
		}
		// Back left bottom (A)
		if cube.A.Coordinate.X != tt.V.X || cube.A.Coordinate.Y != tt.V.Y || cube.A.Coordinate.Z != tt.V.Z {
			t.Log(cube.A.Coordinate)
			t.Log(tt.V)
			t.Error("Invalid 'A' Point - NewCubeByVectorAndLength")
		}
		// Back right bottom (B)
		if cube.B.Coordinate.X != tt.V.X+tt.L || cube.B.Coordinate.Y != tt.V.Y || cube.B.Coordinate.Z != tt.V.Z {
			t.Log(cube.B.Coordinate)
			t.Log(tt.V)
			t.Error("Invalid 'B' Point - NewCubeByVectorAndLength")
		}
		// Back left top (D)
		if cube.D.Coordinate.X != tt.V.X || cube.D.Coordinate.Y != tt.V.Y+tt.L || cube.D.Coordinate.Z != tt.V.Z {
			t.Log(cube.D.Coordinate)
			t.Log(tt.V)
			t.Error("Invalid 'D' Point - NewCubeByVectorAndLength")
		}
		// Back right top (C)
		if cube.C.Coordinate.X != tt.V.X+tt.L || cube.C.Coordinate.Y != tt.V.Y+tt.L || cube.C.Coordinate.Z != tt.V.Z {
			t.Log(cube.C.Coordinate)
			t.Log(tt.V)
			t.Error("Invalid 'C' Point - NewCubeByVectorAndLength")
		}
	}
}
func TestSetupVao(t *testing.T) {
	cube := NewCubeByVectorAndLength(Vector{0, 0, 0}, 1.0)
	// back (a,b,c,d -> abc, acd)
	// right (b,g,f,c -> bgf, bfc)
	// top (c,f,e,d -> cfe, ced)
	// front (g,f,e,h -> gfe, geh)
	// left (e,d,a,h -> eda, eah)
	// bottom (a,b,g,h -> abg, agh)
	expected := []Vector{
		cube.A.Coordinate,
		cube.B.Coordinate,
		cube.C.Coordinate,
		cube.A.Coordinate,
		cube.C.Coordinate,
		cube.D.Coordinate,

		cube.B.Coordinate,
		cube.G.Coordinate,
		cube.F.Coordinate,
		cube.B.Coordinate,
		cube.F.Coordinate,
		cube.C.Coordinate,

		cube.C.Coordinate,
		cube.F.Coordinate,
		cube.E.Coordinate,
		cube.C.Coordinate,
		cube.E.Coordinate,
		cube.D.Coordinate,

		cube.G.Coordinate,
		cube.F.Coordinate,
		cube.E.Coordinate,
		cube.G.Coordinate,
		cube.E.Coordinate,
		cube.H.Coordinate,

		cube.E.Coordinate,
		cube.D.Coordinate,
		cube.A.Coordinate,
		cube.E.Coordinate,
		cube.A.Coordinate,
		cube.H.Coordinate,

		cube.A.Coordinate,
		cube.B.Coordinate,
		cube.G.Coordinate,
		cube.A.Coordinate,
		cube.G.Coordinate,
		cube.H.Coordinate,
	}
	vao := cube.setupVao()
	// px,py,pz,cx,cy,cz
	length := len(vao)
	if length%6 != 0 || length/6 != len(expected) {
		t.Log(length)
		t.Log(len(expected))
		t.Error("Invalid length")
	}
	for i := 0; i < len(expected); i++ {
		if float32(expected[i].X) != vao[i*6] || float32(expected[i].Y) != vao[i*6+1] || float32(expected[i].Z) != vao[i*6+2] {
			t.Log(expected[i])
			t.Log(vao[i*6], vao[i*6+1], vao[i*6+2])
			t.Error("Invalid values")
		}
	}
}
