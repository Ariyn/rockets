package main

import (
	"fmt"
	"github.com/ariyn/rockets"
	"github.com/fogleman/gg"
	"log"
)

const gravity = 9.8
const FPS = 100
const seconds = 2000

func main() {
	objects := make([]*rockets.MassObject, 0)

	dc := gg.NewContext(1000, 1000)

	obj1 := &rockets.MassObject{Mass: 20, R: 10, ID: "obj1"}
	obj1.Position = rockets.Vector3D{X: 500, Y: 500, Z: 0}
	obj1.Velocity = rockets.Vector3D{X: -0.3, Y: 0.5, Z: 0}

	obj2 := &rockets.MassObject{Mass: 10, R: 10, ID: "obj2"}
	obj2.Position = rockets.Vector3D{X: 700, Y: 500, Z: 0}
	obj2.Velocity = rockets.Vector3D{X: 0, Y: 0.4, Z: 0}

	obj3 := &rockets.MassObject{Mass: 15, R: 10, ID: "obj3"}
	obj3.Position = rockets.Vector3D{X: 300, Y: 500, Z: 0}
	obj3.Velocity = rockets.Vector3D{X: 0, Y: -0.4, Z: 0}

	obj4 := &rockets.MassObject{Mass: 50, R: 10, ID: "obj4"}
	obj4.Position = rockets.Vector3D{X: 500, Y: 200, Z: 0}
	obj4.Velocity = rockets.Vector3D{X: -0.1, Y: 0.2, Z: 0}

	objects = append(objects, obj1, obj2, obj3, obj4)
	//objects = append(objects, obj1, obj3)

	for i := 0; i < seconds*FPS; i++ {
		for _, obj := range objects {
			calculateForce(obj, objects)
		}

		// 1초당 한 번씩 이미지 캡처
		if i%FPS == 0 {
			dc.DrawRectangle(0, 0, 1000, 1000)
			dc.SetRGB(0, 0, 0)
			dc.Fill()

			drawObj(dc, obj1, 1)
			drawObj(dc, obj2, 2)
			drawObj(dc, obj3, 3)
			drawObj(dc, obj4, 4)

			dc.SetRGB(1, 1, 1)
			dc.DrawString(fmt.Sprintf("index: %d", i), 30, 870)
			dc.DrawString(fmt.Sprintf("obj1: %f, %f, %f", obj1.Force.X, obj1.Acceleration.X, obj1.Velocity.X), 30, 900)
			dc.DrawString(fmt.Sprintf("obj2: %f, %f, %f", obj2.Force.X, obj2.Acceleration.X, obj2.Velocity.X), 30, 920)
			dc.DrawString(fmt.Sprintf("obj3: %f, %f, %f", obj3.Force.X, obj3.Acceleration.X, obj3.Velocity.X), 30, 940)
			dc.DrawString(fmt.Sprintf("obj4: %f, %f, %f", obj4.Force.X, obj4.Acceleration.X, obj4.Velocity.X), 30, 960)

			dc.DrawString(fmt.Sprintf("obj1: %f, %f / obj3: %f, %f", obj1.Position.X, obj1.Position.Y, obj3.Position.X, obj3.Position.Y), 300, 900)
			dc.SavePNG(fmt.Sprintf("app/simulator/images/%05d.png", i/FPS))

			dc.Clear()

			log.Println(i)
		}
	}
}

func calculateForce(obj *rockets.MassObject, objects []*rockets.MassObject) {
	defer obj.Step(1 / float64(FPS))

	force := rockets.ZeroVector
	for _, obj2 := range objects {
		if obj == obj2 {
			continue
		}

		p := obj2.Position.Sub(obj.Position)
		r := p.Length()
		force = force.Add(p.Normalize().Mul(gravity * obj.Mass * obj2.Mass / (r * r)))
	}

	obj.Force = force
}

func setColor(ctx *gg.Context, index int) {
	switch index {
	case 1:
		ctx.SetRGB(1, 0, 0)
	case 2:
		ctx.SetRGB(0, 0, 1)
	case 3:
		ctx.SetRGB(0, 1, 0)
	case 4:
		ctx.SetRGB(0, 1, 1)
	}
}

func drawObj(ctx *gg.Context, obj *rockets.MassObject, index int) {
	setColor(ctx, index)
	if len(obj.PositionHistory) != 0 {
		ctx.SetLineWidth(1)
		for i, p := range obj.PositionHistory[:len(obj.PositionHistory)-1] {
			if i == 0 || i%FPS != 0 {
				continue
			}

			prev := obj.PositionHistory[i-FPS]
			ctx.DrawLine(prev.X, prev.Y, p.X, p.Y)
		}
		ctx.Stroke()
	}

	ctx.DrawPoint(obj.Position.X, obj.Position.Y, obj.R)
	ctx.Fill()

	ctx.SetRGB(1, 1, 1)
	id := fmt.Sprintf("%d", index)
	w, h := ctx.MeasureString(id)

	ctx.SetRGB(0, 0, 0)
	ctx.DrawString(id, obj.Position.X-w/2, obj.Position.Y+h/2)
}
