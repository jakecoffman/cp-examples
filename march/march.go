package main

import (
	"bytes"
	"cp-examples"
	"embed"
	. "github.com/jakecoffman/cp"
	"image"
	"image/png"
)

func main() {
	space := NewSpace()
	space.Iterations = 10
	space.SetGravity(Vector{0, -100})

	var shape *Shape

	walls := []Vector{
		{-320, -240}, {-320, 240},
		{320, -240}, {320, 240},
		{-320, -240}, {320, -240},
	}

	for i := 0; i < len(walls)-1; i += 2 {
		shape = space.AddShape(NewSegment(space.StaticBody, walls[i], walls[i+1], 0))
		shape.SetElasticity(1)
		shape.SetFriction(1)
		shape.SetFilter(examples.NotGrabbableFilter)
	}

	//list, err := fruits.ReadDir("fruits")
	//if err != nil {
	//	panic(err)
	//}
	//for _, item := range list {
	data, err := fruits.ReadFile("fruits/apple.png")
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	if err = addFruit(space, img); err != nil {
		//panic("For fruit " + item.Name() + ": " + err.Error())
	}
	//}

	examples.Main(space, 1.0/180.0, update, examples.DefaultDraw)
}

func update(space *Space, dt float64) {
	space.Step(dt)
}

func addFruit(space *Space, img image.Image) error {
	b := img.Bounds()
	bb := BB{float64(b.Min.X), float64(b.Min.Y), float64(b.Max.X), float64(b.Max.Y)}

	sampleFunc := func(point Vector) float64 {
		x := point.X
		y := point.Y
		rect := img.Bounds()

		if x < float64(rect.Min.X) || x > float64(rect.Max.X) || y < float64(rect.Min.Y) || y > float64(rect.Max.Y) {
			return 0.0
		}
		_, _, _, a := img.At(int(x), int(y)).RGBA()
		return float64(a) / 0xffff
		//if a == 0xffff {
		//	return 1.0
		//}
		//return 0.0
	}

	lineSet := MarchHard(bb, 1_000, 1_000, 0.01, PolyLineCollectSegment, sampleFunc)
	//MarchSoft(bb, 300, 300, 0.5, PolyLineCollectSegment, &lineSet, sampleFunc)

	//if len(lineSet.Lines) > 1 {
	//	return fmt.Errorf("too many lines: %v", len(lineSet.Lines))
	//}
	for _, line := range lineSet.Lines {
		newLine := line.SimplifyCurves(0.0)
		body := space.AddBody(NewBody(1, MomentForPoly(1, len(newLine.Verts), newLine.Verts, Vector{}, 0)))
		//space.AddShape(NewPolyShape(body, len(newLine.Verts), newLine.Verts, NewTransformIdentity(), 0))
		for i := 0; i < len(newLine.Verts)-1; i++ {
			a := newLine.Verts[i]
			b := newLine.Verts[i+1]
			AddSegment(space, body, a, b, 0)
		}
	}
	return nil
}

//go:embed fruits
var fruits embed.FS

func AddSegment(space *Space, body *Body, a, b Vector, radius float64) *Shape {
	// swap so we always draw the same direction horizontally
	if a.X < b.X {
		a, b = b, a
	}

	seg := NewSegment(body, a, b, radius).Class.(*Segment)
	shape := space.AddShape(seg.Shape)
	shape.SetElasticity(0)
	shape.SetFriction(0.7)

	return shape
}
