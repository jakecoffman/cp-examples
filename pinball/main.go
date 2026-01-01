package main

import (
	"flag"
	"math"
	"math/rand/v2"
	"time"

	examples "github.com/jakecoffman/cp-examples"
	. "github.com/jakecoffman/cp/v2"
)

const (
	ballRadius  = 15
	ballMass    = 0.5
	leverMass   = 10
	leverLength = 60
	leverRadius = 10
	impluseX    = 5000
	swing60     = math.Pi / 3.0
)

func main() {
	flag.Parse()

	space := NewSpace()
	space.Iterations = 10
	space.SetGravity(Vector{0, -200})

	posA := Vector{10, 60}
	posB := Vector{200, 60}
	boxOffset := Vector{0, -120}

	var body1, body2 *Body

	// left
	body1 = addLever(space, posA, boxOffset)
	space.AddConstraint(NewPivotJoint(body1, space.StaticBody, boxOffset.Add(posA)))

	staticRotaryA := space.AddBody(NewStaticBody())
	space.AddConstraint(NewRotaryLimitJoint(body1, staticRotaryA, -swing60*2, -swing60))

	// right
	body2 = addLever(space, posB, boxOffset)
	space.AddConstraint(NewPivotJoint(body2, space.StaticBody, boxOffset.Add(posB)))

	space.AddConstraint(NewRotaryLimitJoint(body2, space.StaticBody, swing60, 2*swing60))

	// use "a" and "l" keys to trigger left and right lever
	examples.HandleKeyFunc = func(char rune) {
		if char == 'a' {
			body1.ApplyImpulseAtLocalPoint(Vector{impluseX, 0}, Vector{0, -leverLength / 2})
		}
		if char == 'l' {
			body2.ApplyImpulseAtLocalPoint(Vector{-impluseX, 0}, Vector{0, -leverLength / 2})
		}
	}

	// spawn balls every two seconds
	go func() {
		for {
			time.Sleep(2000 * time.Millisecond)
			dx := posB.X - posA.X
			ball := space.AddBody(NewBody(ballMass, MomentForCircle(ballMass, 0, ballRadius, Vector{0, 0})))
			ball.SetPosition(Vector{(posA.X+posB.X)/2 - dx/2 + rand.Float64()*dx, 300})
			shape := space.AddShape(NewCircle(ball, ballRadius, Vector{0, 0}))
			shape.SetElasticity(0.5)
			shape.SetFriction(0.7)
		}
	}()

	examples.Main(space, 1.0/60.0, update, draw)
}

func draw(space *Space) {
	examples.DefaultDraw(space)
	examples.DrawString(Vector{-300, 100}, "Use the keys 'a' and 'l' to hit the levers")
}

func update(space *Space, dt float64) {
	space.Step(dt)
}

func addLever(space *Space, pos, boxOffset Vector) *Body {
	a := Vector{0, 15}
	b := Vector{0, -leverLength}

	body := space.AddBody(NewBody(leverMass, MomentForSegment(leverMass, a, b, 0)))
	body.SetPosition(pos.Add(boxOffset.Add(Vector{0, -15})))

	shape := space.AddShape(NewSegment(body, a, b, leverRadius))
	shape.SetElasticity(0)
	shape.SetFriction(0.7)

	return body
}
