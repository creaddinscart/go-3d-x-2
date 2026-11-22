package render

import (
	"math"

	m32 "creaddinscart.studio.go.sl/3d-x-2/math32"
	"creaddinscart.studio.go.sl/3d-x-2/scene"
)

// Renderer draws scenes into a framebuffer.
type Renderer struct {
	FB           *Framebuffer
	Camera       *scene.Camera
	CullBackface bool
	Wireframe    bool
	Samples      int // 1 = off, 2 = 4x AA, 4 = 16x AA

	cur    *Framebuffer // active draw target (FB or supersample buffer)
	sup    *Framebuffer
	lights []scene.Light
	fog    scene.Fog
	camPos m32.Vec3
}

// New builds a renderer targeting fb with cam.
func New(fb *Framebuffer, cam *scene.Camera) *Renderer {
	return &Renderer{FB: fb, Camera: cam, CullBackface: true, Samples: 1}
}

// WithAA enables supersampling anti-aliasing (factor 2 or 4).
func (r *Renderer) WithAA(factor int) *Renderer {
	if factor < 1 {
		factor = 1
	}
	r.Samples = factor
	return r
}

// Render draws the whole scene into the renderer's output framebuffer.
func (r *Renderer) Render(s *scene.Scene) {
	r.lights = s.Lights
	r.fog = s.Fog
	r.camPos = r.Camera.Position

	if r.Samples > 1 {
		sw, sh := r.FB.W*r.Samples, r.FB.H*r.Samples
		if r.sup == nil || r.sup.W != sw || r.sup.H != sh {
			r.sup = NewFramebuffer(sw, sh)
		}
		r.cur = r.sup
		r.renderInto(r.sup, s)
		r.resolve()
		r.cur = r.FB
	} else {
		r.cur = r.FB
		r.renderInto(r.FB, s)
	}
}

func (r *Renderer) renderInto(fb *Framebuffer, s *scene.Scene) {
	fb.Clear(s.Background, 1)
	vp := r.Camera.ViewProjection()
	for _, it := range s.VisibleEntities() {
		mvp := vp.Mul(it.World)
		r.drawMesh(it.Ent, mvp, it.World)
	}
	r.drawPoints(vp, s)
}

// resolve box-downsamples the supersample buffer into the output framebuffer.
func (r *Renderer) resolve() {
	n := r.Samples
	sx, sy := r.sup.W, r.sup.H
	for y := 0; y < r.FB.H; y++ {
		for x := 0; x < r.FB.W; x++ {
			var sr, sg, sb, sa float32
			for j := 0; j < n; j++ {
				for i := 0; i < n; i++ {
					c := r.sup.Pix[(y*n+j)*sx+(x*n+i)]
					sr += c.R
					sg += c.G
					sb += c.B
					sa += c.A
				}
			}
			k := float32(n * n)
			r.FB.Pix[y*r.FB.W+x] = m32.NewRGBA(sr/k, sg/k, sb/k, sa/k)
		}
	}
	_ = sy
}

// shade evaluates Blinn-Phong lighting plus fog for one fragment.
func (r *Renderer) shade(mat scene.Material, worldPos, normal m32.Vec3, uv m32.Vec2) m32.Color {
	n := normal.Normalize()
	view := r.camPos.Sub(worldPos).Normalize()

	albedo := mat.Diffuse
	ambientTerm := mat.Ambient
	if mat.Texture != nil {
		tex := mat.Texture.Sample(uv.X, uv.Y)
		albedo = albedo.Modulate(tex)
		ambientTerm = ambientTerm.Modulate(tex)
	}

	out := mat.Emissive
	out.R += ambientTerm.R
	out.G += ambientTerm.G
	out.B += ambientTerm.B

	for _, L := range r.lights {
		var lightDir m32.Vec3
		atten := float32(1)
		switch L.Kind {
		case scene.LightAmbient:
			a := L.Color.Mul(L.Intensity)
			out.R += ambientTerm.R * a.R
			out.G += ambientTerm.G * a.G
			out.B += ambientTerm.B * a.B
			continue
		case scene.LightDirectional:
			lightDir = L.Direction.Neg()
		case scene.LightPoint:
			toL := L.Position.Sub(worldPos)
			d := toL.Length()
			if d == 0 {
				continue
			}
			lightDir = toL.Div(d)
			if L.Range > 0 {
				atten = m32.Clamp01(1 - d/L.Range)
			}
		case scene.LightSpot:
			toL := worldPos.Sub(L.Position)
			d := toL.Length()
			if d == 0 {
				continue
			}
			toL = toL.Div(d)
			if toL.Dot(L.Direction) < L.SpotCos {
				continue
			}
			lightDir = toL.Neg()
			if L.Range > 0 {
				atten = m32.Clamp01(1 - d/L.Range)
			}
		}

		ndl := n.Dot(lightDir)
		if ndl < 0 {
			ndl = 0
		}
		lc := L.Color.Mul(L.Intensity * atten)
		out.R += albedo.R * lc.R * ndl
		out.G += albedo.G * lc.G * ndl
		out.B += albedo.B * lc.B * ndl
		if ndl > 0 {
			hv := lightDir.Add(view).Normalize()
			if spe := n.Dot(hv); spe > 0 {
				spe = float32(math.Pow(float64(spe), float64(mat.Shininess)))
				out.R += mat.Specular.R * lc.R * spe
				out.G += mat.Specular.G * lc.G * spe
				out.B += mat.Specular.B * lc.B * spe
			}
		}
	}

	out.R = m32.Clamp01(out.R)
	out.G = m32.Clamp01(out.G)
	out.B = m32.Clamp01(out.B)
	out.A = mat.Opacity

	// fog blends toward the fog color by view distance
	if r.fog.Mode != scene.FogNone {
		depth := worldPos.Sub(r.camPos).Length()
		f := r.fog.Factor(depth)
		if f > 0 {
			out.R = m32.Lerp(out.R, r.fog.Color.R, f)
			out.G = m32.Lerp(out.G, r.fog.Color.G, f)
			out.B = m32.Lerp(out.B, r.fog.Color.B, f)
		}
	}
	return out
}
