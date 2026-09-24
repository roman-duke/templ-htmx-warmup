package main

import (
	"log"

	"github.com/roman-duke/templ-htmx-warmup/blog/internal/core"
)

// this just loads the map on startup so that it's not empty during each "hot reload"
// Each entry is created through Service.CreatePost, which assigns a generated id + timestamps and
// validates it against the domain invariants (unique title, content >= 50 chars).
func seedPosts() []core.CreatePostInput {
	return []core.CreatePostInput{
		{
			Title:   "The Second Law of Thermodynamics",
			Content: "Entropy in an isolated system never decreases over time. This single idea gives the arrow of time its direction and explains why heat flows from hot to cold, never the reverse.",
		},
		{
			Title:   "Understanding Quantum Entanglement",
			Content: "When two particles become entangled, measuring one instantly fixes the state of the other, no matter the distance between them. Einstein called it *spooky action at a distance*.",
		},
		{
			Title:   "Newton's Laws of Motion Explained",
			Content: "Three deceptively simple rules govern almost all everyday motion: inertia, F = ma, and equal-and-opposite reaction. Together they carried physics for over two centuries.",
		},
		{
			Title:   "Special Relativity and Time Dilation",
			Content: "As an object approaches the speed of light, time for it slows relative to a stationary observer. Moving clocks tick slower, and this effect is measured every day in particle accelerators.",
		},
		{
			Title:   "What Is the Higgs Boson?",
			Content: "The Higgs boson is the excitation of the Higgs field, the mechanism that gives fundamental particles their mass. Its 2012 discovery at CERN completed the Standard Model.",
		},
		{
			Title:   "The Double-Slit Experiment",
			Content: "Fire electrons one at a time through two slits and an interference pattern still builds up, as if each particle passes through both slits at once. It is quantum weirdness made visible.",
		},
		{
			Title:   "Black Holes and Event Horizons",
			Content: "A black hole's event horizon is the boundary beyond which nothing, not even light, can escape. Cross it and every possible path leads inexorably toward the singularity.",
		},
		{
			Title:   "Maxwell's Equations Demystified",
			Content: "Four elegant equations unify electricity and magnetism and reveal that light itself is an electromagnetic wave. They are the foundation of every antenna, motor, and fiber optic cable.",
		},
		{
			Title:   "The Photoelectric Effect",
			Content: "Light knocks electrons out of a metal only above a threshold frequency, not merely with enough intensity. Einstein explained this by treating light as discrete quanta, or photons.",
		},
		{
			Title:   "Wave-Particle Duality",
			Content: "Matter and light behave as both waves and particles depending on how you observe them. This duality sits at the very heart of quantum mechanics and defies classical intuition.",
		},
	}
}

// seed loads the starter content into the store via the service. A failure here
// is logged rather than fatal: the app can still run with an empty store.
func seed(svc *core.Service) {
	for _, in := range seedPosts() {
		if _, err := svc.CreatePost(in); err != nil {
			log.Printf("seed: failed to create %q: %v", in.Title, err)
		}
	}
}
