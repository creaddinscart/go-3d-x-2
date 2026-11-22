package anim

import (
	"math/rand"

	m32 "creaddinscart.studio.go.sl/3d-x-2/math32"
	"creaddinscart.studio.go.sl/3d-x-2/scene"
	"creaddinscart.studio.go.sl/3d-x-2/texture"
)

// Particle is a single simulated point sprite.
type Particle struct {
	Pos, Vel m32.Vec3
	Life     float32
	MaxLife  float32
	Color    m32.Color
	Size     float32
	Additive bool
	Tex      *texture.Texture
}

// ParticleConfig configures a ParticleSystem emitter.
type ParticleConfig struct {
	Origin    m32.Vec3
	Spread    m32.Vec3 // random velocity spread
	Direction m32.Vec3 // base emission direction (normalized)
	Speed     float32
	Rate      float32 // particles per second
	Lifetime  float32 // seconds
	Size      float32
	Color     m32.Color
	Gravity   float32 // downward accel (units/s^2)
	FadeOut   bool    // shrink alpha toward end of life
	Additive  bool
	Max       int
	Tex       *texture.Texture
}

// ParticleSystem spawns and integrates billboard particles.
type ParticleSystem struct {
	cfg   ParticleConfig
	parts []Particle
	acc   float32
	rng   *rand.Rand
}

// NewParticleSystem builds an emitter from a config (Max defaults to 2000).
func NewParticleSystem(cfg ParticleConfig) *ParticleSystem {
	if cfg.Max <= 0 {
		cfg.Max = 2000
	}
	if cfg.Direction.LengthSq() == 0 {
		cfg.Direction = m32.UnitY
	}
	cfg.Direction = cfg.Direction.Normalize()
	return &ParticleSystem{cfg: cfg, rng: rand.New(rand.NewSource(1))}
}

func (s *ParticleSystem) randn() float32 { return float32(s.rng.NormFloat64()) }

// Update advances the simulation by dt seconds, spawning and culling particles.
func (s *ParticleSystem) Update(dt float32) {
	s.acc += s.cfg.Rate * dt
	for s.acc >= 1 && len(s.parts) < s.cfg.Max {
		s.acc -= 1
		v := s.cfg.Direction.Mul(s.cfg.Speed).Add(m32.V3(
			s.randn()*s.cfg.Spread.X,
			s.randn()*s.cfg.Spread.Y,
			s.randn()*s.cfg.Spread.Z,
		))
		s.parts = append(s.parts, Particle{
			Pos: s.cfg.Origin, Vel: v,
			Life: s.cfg.Lifetime, MaxLife: s.cfg.Lifetime,
			Color: s.cfg.Color, Size: s.cfg.Size, Additive: s.cfg.Additive, Tex: s.cfg.Tex,
		})
	}
	alive := s.parts[:0]
	for _, p := range s.parts {
		p.Life -= dt
		if p.Life <= 0 {
			continue
		}
		p.Vel.Y -= s.cfg.Gravity * dt
		p.Pos = p.Pos.Add(p.Vel.Mul(dt))
		alive = append(alive, p)
	}
	s.parts = alive
}

// AppendPoints writes the live particles into dst as billboard scene points and
// returns the extended slice.
func (s *ParticleSystem) AppendPoints(dst []scene.Point) []scene.Point {
	for _, p := range s.parts {
		col := p.Color
		if s.cfg.FadeOut {
			col.A = p.Color.A * (p.Life / p.MaxLife)
		}
		dst = append(dst, scene.Point{
			Position: p.Pos, Color: col, Size: p.Size, Texture: p.Tex, Additive: p.Additive,
		})
	}
	return dst
}

// Count returns the number of live particles.
func (s *ParticleSystem) Count() int { return len(s.parts) }
