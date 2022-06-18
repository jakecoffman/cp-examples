package main

import (
	"bytes"
	"cp-examples"
	"embed"
	. "github.com/jakecoffman/cp"
	"image"
	"image/png"
	"log"
)

func main() {
	space := NewSpace()
	space.Iterations = 10
	space.SetGravity(Vector{0, -100})
	space.SetDamping(.9)

	var shape *Shape

	walls := []Vector{
		{-320, -240}, {-320, 240},
		{320, -240}, {320, 240},
		{-320, -240}, {320, -240},
	}

	for i := 0; i < len(walls)-1; i += 2 {
		shape = space.AddShape(NewSegment(space.StaticBody, walls[i], walls[i+1], 0))
		shape.SetElasticity(1)
		shape.SetFriction(0)
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
	addFruit(space, img)

	examples.Main(space, 1.0/180.0, update, examples.DefaultDraw)
}

func update(space *Space, dt float64) {
	space.Step(dt)
}

func addFruit(space *Space, img image.Image) {
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
	}

	lineSet := MarchHard(bb, 10_000, 10_000, 0.2, PolyLineCollectSegment, sampleFunc)
	//MarchSoft(bb, 300, 300, 0.5, PolyLineCollectSegment, &lineSet, sampleFunc)

	if len(lineSet.Lines) > 1 {
		log.Panicln("Num lines:", len(lineSet.Lines))
	}
	line := lineSet.Lines[0].SimplifyCurves(0.2)
	offset := Vector{float64(b.Max.X-b.Min.X) / 2., float64(b.Max.Y-b.Min.Y) / 2.}
	// center the verts on origin
	for i, l := range line.Verts {
		line.Verts[i] = l.Sub(offset)
	}

	body := space.AddBody(NewBody(10, MomentForPoly(10, len(line.Verts), line.Verts, Vector{}, 1)))
	fruit := space.AddShape(NewPolyShape(body, len(line.Verts), line.Verts, NewTransformIdentity(), 0))
	fruit.SetElasticity(.5)
	for i := 0; i < len(line.Verts)-1; i++ {
		a := line.Verts[i]
		b := line.Verts[i+1]
		AddSegment(space, body, a, b, 0)
	}
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
