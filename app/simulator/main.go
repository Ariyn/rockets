package main

import (
	"fmt"
	"github.com/ariyn/rockets"
	"github.com/fogleman/gg"
	"math"
	"sync"
)

const FPS = 60
const seconds = 1000
const drawerCount = 20

const radiusOfEarth = 6.375427e+06 // in meters

var massOfEarth = 5.97219e+24

var cameraSize = rockets.Vector3D{X: 1000, Y: 1000, Z: 0}
var viewportZoomRatio = 1e-5

func main() {
	c := make(chan drawStruct, drawerCount*5)
	wg := sync.WaitGroup{}

	wg.Add(drawerCount)
	for i := 0; i < drawerCount; i++ {
		go drawer(c, &wg)
	}

	objects := make([]*rockets.MassObject, 0)

	obj1 := &rockets.MassObject{Mass: massOfEarth, R: radiusOfEarth, ID: "Earth", IsKinematic: true}
	obj1.Position = rockets.Vector3D{X: 0, Y: 0, Z: 0}
	obj1.Angle = rockets.Vector3D{X: 0, Y: 1, Z: 0}

	obj2 := &rockets.MassObject{Mass: 100, R: 10, ID: "Satellite"}
	obj2.Position = rockets.Vector3D{X: radiusOfEarth + 800*1e3, Y: 0, Z: 0}
	obj2.Velocity = rockets.Vector3D{X: 0, Y: 7.67e3, Z: 0}
	obj2.Angle = rockets.Vector3D{X: 0, Y: 1, Z: 0}
	obj2.AngularVel = rockets.Vector3D{X: 0, Y: 0, Z: 0.2}

	objects = append(objects, obj1, obj2)

	copiedObjects := make([]rockets.MassObject, len(objects))
	for i, obj := range objects {
		copiedObjects[i] = *obj
	}

	c <- drawStruct{objects: copiedObjects, cameraPosition: obj2.Position, index: 0}

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
			c <- drawStruct{objects: copiedObjects, cameraPosition: obj2.Position, index: i/FPS + 1}
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
		ctx.SetRGB(0, 1, 0)
	case 2:
		ctx.SetRGB(0, 0, 1)
	case 3:
		ctx.SetRGB(0, 1, 0)
	case 4:
		ctx.SetRGB(0, 1, 1)
	}
}

type drawStruct struct {
	objects        []rockets.MassObject
	cameraPosition rockets.Vector3D
	index          int
}

func drawer(c <-chan drawStruct, wg *sync.WaitGroup) {
	defer wg.Done()

	for ds := range c {
		ctx := gg.NewContext(int(cameraSize.X), int(cameraSize.Y))
		ctx.SetRGB(0, 0, 0)
		ctx.DrawRectangle(0, 0, cameraSize.X, cameraSize.Y)
		ctx.Fill()

		for i, obj := range ds.objects {
			drawObj(ctx, ds.cameraPosition, &obj, i+1)
		}

		ctx.SetRGB(1, 1, 1)
		ctx.DrawString(fmt.Sprintf("index: %d", ds.index), 30, 870)

		for i, obj := range ds.objects[1:] {
			ctx.DrawString(fmt.Sprintf("obj%d: F(%f, %f) / Alt: %f", i+1, obj.Force.X, obj.Force.Y, (obj.Position.Length()-radiusOfEarth)/1e3), 30, float64(900+i*20))
		}

		err := ctx.SavePNG(fmt.Sprintf("images/%05d.png", ds.index))
		if err != nil {
			panic(err)
		}
	}
}

func cameraCoordinate(cameraPosition, pos rockets.Vector3D) rockets.Vector3D {
	return pos.Sub(cameraPosition).Mul(viewportZoomRatio).Add(cameraSize.Mul(0.5))
}

func drawPositionLine(ctx *gg.Context, cameraPosition rockets.Vector3D, positions []rockets.Vector3D, index int) {
	ctx.SetLineWidth(1)
	setColor(ctx, index)

	prevIndex := 0
	for i, p := range positions {
		if i == 0 {
			continue
		}

		prev := cameraCoordinate(cameraPosition, positions[prevIndex])
		curr := cameraCoordinate(cameraPosition, p)

		if cameraSize.Length() <= prev.Distance(cameraSize.Mul(0.5)) || cameraSize.Length() <= curr.Distance(cameraSize.Mul(0.5)) {
			continue
		}
		if prev.Distance(curr) < 3 {
			continue
		}

		ctx.DrawLine(prev.X, prev.Y, curr.X, curr.Y)
		prevIndex = i
	}
	ctx.Stroke()
}

func drawObj(ctx *gg.Context, cameraPosition rockets.Vector3D, obj *rockets.MassObject, index int) {
	setColor(ctx, index)

	if len(obj.PositionHistory) != 0 {
		drawPositionLine(ctx, cameraPosition, obj.PositionHistory, index)
	}

	newPosition := cameraCoordinate(cameraPosition, obj.Position)
	if cameraSize.Length() <= newPosition.Distance(cameraSize.Mul(0.5)) {
		return
	}

	ctx.DrawPoint(newPosition.X, newPosition.Y, math.Max(obj.R*viewportZoomRatio, 2))
	ctx.Fill()

	headingNorm := obj.Angle.Normalize()
	headingIndicatorStart := cameraCoordinate(cameraPosition, obj.Position.Add(headingNorm.Mul(obj.R)))
	headingIndicatorEnd := cameraCoordinate(cameraPosition, obj.Position.Add(headingNorm.Mul(obj.R+10)))

	ctx.SetLineWidth(5)
	ctx.DrawLine(
		headingIndicatorStart.X,
		headingIndicatorStart.Y,
		headingIndicatorEnd.X,
		headingIndicatorEnd.Y,
	)
	ctx.Stroke()

	ctx.SetRGB(1, 1, 1)
	id := fmt.Sprintf("%d", index)
	w, h := ctx.MeasureString(id)

	ctx.SetRGB(0, 0, 0)

	ctx.DrawString(id, newPosition.X-w/2, newPosition.Y+h/2)
}
