package main

import (
	"fmt"
	"github.com/ariyn/rockets"
	"github.com/fogleman/gg"
	"sync"
)

const FPS = 60
const seconds = 300
const drawerCount = 20

func main() {
	c := make(chan drawStruct, drawerCount*5)
	wg := sync.WaitGroup{}

	wg.Add(drawerCount)
	for i := 0; i < drawerCount; i++ {
		go drawer(c, &wg)
	}

	objects := make([]*rockets.MassObject, 0)

	obj1 := &rockets.MassObject{Mass: 30, R: 20, ID: "obj1", IsKinematic: true}
	obj1.Position = rockets.Vector3D{X: 500, Y: 500, Z: 0}
	obj1.Angle = rockets.Vector3D{X: 0, Y: 1, Z: 0}

	obj2 := &rockets.MassObject{Mass: 10, R: 10, ID: "obj2"}
	obj2.Position = rockets.Vector3D{X: 700, Y: 500, Z: 0}
	obj2.Velocity = rockets.Vector3D{X: 0, Y: 1.2, Z: 0}
	obj2.Angle = rockets.Vector3D{X: 0, Y: 1, Z: 0}
	obj2.AngularVel = rockets.Vector3D{X: 0, Y: 0, Z: 0.2}

	objects = append(objects, obj1, obj2)

	for i := 0; i < seconds*FPS; i++ {
		for _, obj := range objects {
			calculateForce(obj, objects)
		}

		for _, obj := range objects {
			obj.Step(1 / float64(FPS))
		}

		// 1초당 한 번씩 이미지 캡처
		if i%FPS == 0 {
			copiedObjects := make([]rockets.MassObject, len(objects))
			for i, obj := range objects {
				copiedObjects[i] = *obj
			}
			c <- drawStruct{objects: copiedObjects, index: i / FPS}
		}
	}

	close(c)
	wg.Wait()
}

func calculateForce(obj *rockets.MassObject, objects []*rockets.MassObject) {
	force := rockets.ZeroVector
	for _, obj2 := range objects {
		if obj == obj2 {
			continue
		}

		obj.CalculateForce(obj2)
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

type drawStruct struct {
	objects []rockets.MassObject
	index   int
}

func drawer(c <-chan drawStruct, wg *sync.WaitGroup) {
	defer wg.Done()

	for ds := range c {
		ctx := gg.NewContext(1000, 1000)
		ctx.SetRGB(0, 0, 0)
		ctx.DrawRectangle(0, 0, 1000, 1000)
		ctx.Fill()

		for i, obj := range ds.objects {
			drawObj(ctx, &obj, i+1)
		}

		ctx.SetRGB(1, 1, 1)
		ctx.DrawString(fmt.Sprintf("index: %d", ds.index), 30, 870)

		for i, obj := range ds.objects {
			ctx.DrawString(fmt.Sprintf("obj%d: %f, %f", i+1, obj.Force.X, obj.Force.Y), 30, float64(900+i*20))
		}

		err := ctx.SavePNG(fmt.Sprintf("images/%05d.png", ds.index))
		if err != nil {
			panic(err)
		}
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

	headingNorm := obj.Angle.Normalize()
	headingIndicatorStart := obj.Position.Add(headingNorm.Mul(obj.R))
	headingIndicatorEnd := obj.Position.Add(headingNorm.Mul(obj.R + 10))

	ctx.SetLineWidth(5)
	ctx.DrawLine(headingIndicatorStart.X, headingIndicatorStart.Y, headingIndicatorEnd.X, headingIndicatorEnd.Y)
	ctx.Stroke()

	ctx.SetRGB(1, 1, 1)
	id := fmt.Sprintf("%d", index)
	w, h := ctx.MeasureString(id)

	ctx.SetRGB(0, 0, 0)
	ctx.DrawString(id, obj.Position.X-w/2, obj.Position.Y+h/2)
}
