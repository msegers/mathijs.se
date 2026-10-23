+++
title = 'The Java bug I should have used in interviews'
date = '2026-10-23'
draft = true
tags = ['java', 'debugging']
summary = 'A debugging story from a service-to-monolith migration that I remembered ten years too late.'
+++

More than ten years ago, when I was still a junior developer (at least according to my job title), we were consolidating a bunch of Java services into a monolith.

Yes, we were moving away from separate services before microservices became everybody's favourite architectural decision.

The services had packages and dependencies that weren't always in sync, and bringing everything together introduced some interesting runtime problems.

I remember debugging one particularly annoying issue. After following the code paths and looking at how everything was wired together, I eventually traced it to two libraries doing similar things and getting in each other's way.

I believe they were annotation-driven Java REST frameworks, with some sort of routing or registration collision. It's been over a decade, so I'm not going to pretend I remember the exact libraries or exception.

The interesting part was that everything looked reasonable in isolation. It was the interaction between components that broke things.

That's the kind of debugging I actually enjoy: understanding the system well enough to find the assumption that turned out to be wrong.

And now I realise this would have been an excellent story for a technical interview.

Naturally, I remembered it about ten years too late.
