+++
title = 'The one-FPS microscope and our performance test suite'
date = '2026-10-23'
draft = true
tags = ['performance', 'debugging']
summary = 'Silverlight, AngularJS, IE9, and an enterprise interpretation of performance testing.'
+++

At one point I worked on rewriting a Silverlight viewer for digitised biopsy microscope slides into AngularJS.

The slides were essentially enormous images that you could zoom into and navigate around. Actually quite interesting software, as far as enterprise applications go.

The catch? Most users were still on Internet Explorer 9.

The viewer could end up rendering at something like **one frame per second**.

And then there was the performance test suite.

We had a release-cycle test script that took weeks to complete. There were only a few performance-related tests, and one of them was roughly:

> Open the laptop, go into the Windows system settings, and check whether the CPU and RAM are above the specified minimums.

That's it. Check the computer's specifications.

Meanwhile, the application was effectively a slideshow. Which is especially unfortunate when the entire point is to view microscope slides interactively.

I understand why hardware requirements matter. But checking whether a laptop meets the minimum requirements isn't testing whether your software performs acceptably.

I don't remember enough to claim how the project ended or whether we eventually fixed the frame rate.

What stuck with me is the wonderful enterprise mismatch between what a test says it's testing and what the user actually experiences.

Somewhere, a spreadsheet probably had a green cell.
