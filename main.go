package main

import (
	"log"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glh"
	mgl "github.com/go-gl/mathgl/mgl64"
)

func MakeCube() *glh.MeshBuffer {
	cube := glh.NewMeshBuffer(
		glh.RenderBuffered,
		glh.NewPositionAttr(3, gl.FLOAT, gl.STATIC_DRAW),
	)

	cube.Add([]float32{
		-1, -1, -1, -1, -1, 1,
		-1, -1, -1, -1, 1, -1,
		-1, -1, -1, 1, -1, -1,

		-1, 1, 1, -1, -1, 1,
		-1, 1, 1, -1, 1, -1,
		-1, 1, 1, 1, 1, 1,

		1, -1, 1, -1, -1, 1,
		1, -1, 1, 1, -1, -1,
		1, -1, 1, 1, 1, 1,

		1, 1, -1, -1, 1, -1,
		1, 1, -1, 1, -1, -1,
		1, 1, -1, 1, 1, 1,
	})

	return cube
}

func LoadMatrix(mat mgl.Mat4) {
	arg := [16]float64(mat)
	gl.LoadMatrixd(&arg)
}

func GetMatrix(which gl.GLenum) mgl.Mat4 {
	switch which {
	case gl.MODELVIEW:
		which = gl.MODELVIEW_MATRIX
	case gl.PROJECTION:
		which = gl.PROJECTION_MATRIX
	}
	var mat [16]float64
	gl.GetDoublev(which, mat[:])
	return mat
}

func MulMatrix(mat mgl.Mat4) {
	var which [1]int32
	gl.GetIntegerv(gl.MATRIX_MODE, which[:])
	prev := mgl.Mat4(GetMatrix(gl.GLenum(which[0])))
	// LoadMatrix(mat.Mul4(prev))
	LoadMatrix(prev.Mul4(mat))
}

func LoadQuat(quat mgl.Quat) {
	LoadMatrix(quat.Normalize().Mat4())
}

func Vertex(v mgl.Vec3) {
	gl.Vertex3d(v[0], v[1], v[2])
}

func main() {

	const (
		startFullscreen = false
		windowTitle     = "ToggleFullscreen"
		defW            = 512
		defH            = 512
	)

	window, err := NewWindow(windowTitle, startFullscreen, defW, defH, nil)
	if err != nil {
		log.Fatalln("Unable to create window: ", err)
	}
	defer window.Destroy()

	window.SetFramebufferSizeCallback(func(_ *glfw.Window, w, h int) {
		gl.Viewport(0, 0, w, h)

		ratio := float64(h) / float64(w)

		p := mgl.Ortho(-1, 1, -ratio, ratio, -10, 10)
		const s = 0.2
		p = p.Mul4(mgl.Scale3D(s, s, s))

		// p := mgl.Frustum(-1, 1, -ratio, ratio, 4, 100)

		p = p.Mul4(mgl.Translate3D(0, 0, -20))

		gl.MatrixMode(gl.PROJECTION)
		LoadMatrix(p)
	})

	arcball := NewArcball(window.Window)

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action != glfw.Press {
			return
		}
		switch key {
		case glfw.KeyR:
			arcball.Reset()
		}
	})

	cube := MakeCube()

	// ident := mgl.QuatRotate(0.2, mgl.Vec3{1, 0, 0})
	ident := mgl.QuatIdent()
	// delta := mgl.QuatRotate(0, mgl.Vec3{0, 1, 0})

	for !window.ShouldClose() {

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()

		LoadQuat(ident.Mul(arcball.Rotation()))

		// gl.Color4f(1, 1, 1, 1)
		// cube.Render(gl.LINES)
		_ = cube

		arcball.Draw()
		// gl.LoadIdentity()
		// i = [16]float64(arcball.Rotation().Mat4())
		// gl.MultMatrixd(&i)
		// gl.Rotated(90, 1, 0, 0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
