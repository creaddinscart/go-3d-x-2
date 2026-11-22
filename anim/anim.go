package anim

import (
	"math"

	m32 "creaddinscart.studio.go.sl/3d-x-2/math32"
	"creaddinscart.studio.go.sl/3d-x-2/scene"
)

// Tween animates a single float from -> to over a duration.
type Tween struct {
	From, To, Duration float32
	Ease               EaseFunc
}

// NewTween builds a tween with the given easing.
func NewTween(from, to, duration float32, ease EaseFunc) Tween {
	if ease == nil {
		ease = EaseLinear
	}
	return Tween{From: from, To: to, Duration: duration, Ease: ease}
}

// Value returns the tweened value at local time t (clamped to the duration).
func (tw Tween) Value(t float32) float32 {
	if tw.Duration <= 0 {
		return tw.To
	}
	return m32.Lerp(tw.From, tw.To, tw.Ease(m32.Clamp01(t/tw.Duration)))
}

// Keyframe is a single (time, value) sample for a Vec3 track.
type Keyframe struct {
	Time  float32
	Value m32.Vec3
	Ease  EaseFunc
}

// Vec3Track interpolates a sequence of keyframes.
type Vec3Track struct{ Keys []Keyframe }

// Sample evaluates the track at absolute time t.
func (tr Vec3Track) Sample(t float32) m32.Vec3 {
	if len(tr.Keys) == 0 {
		return m32.Zero3
	}
	if t <= tr.Keys[0].Time {
		return tr.Keys[0].Value
	}
	last := tr.Keys[len(tr.Keys)-1]
	if t >= last.Time {
		return last.Value
	}
	for i := 0; i+1 < len(tr.Keys); i++ {
		a, b := tr.Keys[i], tr.Keys[i+1]
		if t >= a.Time && t <= b.Time {
			u := m32.InvLerp(a.Time, b.Time, t)
			if a.Ease != nil {
				u = a.Ease(u)
			}
			return a.Value.Lerp(b.Value, u)
		}
	}
	return last.Value
}

// Clip is a named set of position tracks played together.
type Clip struct {
	Name      string
	Duration  float32
	Positions map[string]Vec3Track
}

// NewClip returns an empty clip.
func NewClip(name string) *Clip { return &Clip{Name: name, Positions: map[string]Vec3Track{}} }

// AddPosition registers a position track for a node name.
func (c *Clip) AddPosition(node string, keys []Keyframe) {
	c.Positions[node] = Vec3Track{Keys: keys}
	if len(keys) > 0 && keys[len(keys)-1].Time > c.Duration {
		c.Duration = keys[len(keys)-1].Time
	}
}

// Animator plays a clip and applies results to a scene by node name.
type Animator struct {
	Clip   *Clip
	Time   float32
	Speed  float32
	Loop   bool
	Scene  *scene.Scene
	byName map[string]*scene.Node
}

// NewAnimator binds a clip to a scene, indexing nodes by name.
func NewAnimator(clip *Clip, s *scene.Scene) *Animator {
	an := &Animator{Clip: clip, Speed: 1, Scene: s, byName: map[string]*scene.Node{}}
	var walk func(n *scene.Node)
	walk = func(n *scene.Node) {
		an.byName[n.Name] = n
		for _, c := range n.Children() {
			walk(c)
		}
	}
	walk(s.Root)
	return an
}

// Advance steps the animation clock by dt seconds and applies the pose.
func (an *Animator) Advance(dt float32) {
	an.Time += dt * an.Speed
	if an.Clip == nil {
		return
	}
	t := an.Time
	if an.Clip.Duration > 0 {
		if an.Loop {
			t = float32(math.Mod(float64(t), float64(an.Clip.Duration)))
		} else if t > an.Clip.Duration {
			t = an.Clip.Duration
		}
	}
	for name, tr := range an.Clip.Positions {
		if n, ok := an.byName[name]; ok {
			n.Position = tr.Sample(t)
		}
	}
}
