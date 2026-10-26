+++
title = 'Learning Blender while building a character for Tribute'
date = '2026-10-25'
draft = true
tags = ['learning', 'development']
summary = 'Notes from learning Blender to model and rig the first character for my indie game, Tribute.'
image = { src = 'image.png', thumb = 'thumb.png', alt = 'A low-poly Tribute character in Blender with reference images behind it' }
+++

Working notes for a post about learning Blender while making the first character for Tribute. Not a tutorial.

[TODO: rewrite in my own voice and decide how much of the process to show.]

## Why Blender?

- Tribute is an indie game in Rust, built around a settlement of NPCs.
- I want characters with a distinctive, slightly absurd, stylized low-poly appearance.
- I could have leaned entirely on premade assets. I decided to learn Blender and model the first character myself.
- The goal was not realism. It was a silhouette that reads at game distance and looks like it belongs in Tribute.
- First character brief: an exaggerated, top-heavy silhouette.
  - Oversized torso.
  - Small legs.
  - A somewhat egg-shaped head.
  - Chunky hands and boots.
- [TODO: expand on why this shape felt right for Tribute.]

## Learning the basics

- Started with a beveled cube for the head.
- Experimented with rounded geometry and low polygon counts.
- Initially modelled a torso, then realized the tunic itself could serve as the visible body.
  - The torso underneath was mostly doing work nobody would see.
- Used an open cylinder with horizontal edge loops to create a tunic.
- Added a belt and adjusted the silhouette.

> [TODO: insert Blender screenshot — tunic blockout]

- Things I had to learn, or unlearn:
  - Edit Mode versus Object Mode.
  - Origins versus the 3D cursor.
  - Edge loops.
  - Modifiers.
  - Selection.
- Accidentally added new objects to existing meshes more than once.
- Scaled geometry flat to zero and then discovered that scaling it back up does not restore its dimensions.
- Figured out mirroring for symmetrical limbs.

## The hands were unexpectedly annoying

- Tried modelling a simple mitten-like hand with a separate thumb.
- Attempted finger grooves, bevels and other adjustments.
- Ended up with messy geometry.
- Eventually simplified the shape.

> [TODO: insert Blender screenshot — mitten hand]

- Realized individual fingers and finger animation probably are not necessary for Tribute.
- A fixed, loose gripping pose should handle many basic interactions.
- [TODO: expand on where the hand ended up vs. where I wanted it.]

## Boots and limbs

- Built stylized arms and legs from cylinders with a few carefully shaped edge loops.
- The legs ended up somewhat like chicken drumsticks.
- Created chunky wedge-shaped boots.
- Experimented with a rolled boot cuff.

> [TODO: insert Blender screenshot — leg]

- Learned that a recognizable silhouette often matters more than geometric detail.

## Rigging

- Created a basic armature manually.
- Added spine, pelvis, neck, head, shoulder, upper-arm, forearm, hand, thigh, shin and foot bones.
- Used Automatic Weights.
- Tested bending elbows, knees and moving the neck.
- The initial results were surprisingly usable.

> [TODO: insert Blender screenshot — finished character]

- Some deformation issues remain, especially around boot cuffs.
- A proper walking animation has not been created yet.

> [TODO: insert Blender screenshot — first rigging test]

## Decisions about Tribute's art direction

- Prefer distinctive silhouettes over realistic anatomy.
- Use simple solid-coloured materials initially.
- Avoid modelling hidden anatomy beneath clothing unnecessarily.
- Keep clothes, hands and other elements modular.
- Potentially use Shape Keys for variations in body weight.
- Eventually support different body shapes, clothing and hairstyles.
- Balance low-poly appearance and performance without obsessively optimizing every polygon.

## Lessons and reflections

Writing prompts, rather than conclusions:

- How does learning a visual tool differ from programming? What transferred, and what did not?
- Why was experimenting and making mistakes useful here? Which mistakes were worth keeping?
- How easy was it to spend too long perfecting tiny details?
- When did I know it was time to move on?
- What did it feel like to see a character move with a skeleton I built myself?
- Does this change how I think about custom animations versus premade assets?
- [TODO: add the specific moment that made the process click, if there was one.]

## What's next?

- Make a short, sturdy, dwarflike walking animation.
- Test the rig more thoroughly.
- Add facial features and hair.
- Experiment with materials.
- Eventually import the character into Tribute.
