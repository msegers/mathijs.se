+++
date = '2025-09-20T13:00:46+02:00'
title = 'Setup for my dotfiles'
summary = "Simple automation for my dotfiles."
tags = ['learning', 'linux']
+++

A while ago I reinstalled Arch, I do this every now and then. However recently I first tried to use `archinstall` which is bundled with the ISO, which I always have near me for some reason.
Today I got fed up with having the decision to make a way too small home partition (I assumed most programs wouldn't install in home by default, and understimated some repositories in size).
So it got me thinking about automating things, since setup can be tedious, I started just throwing all command into a bash script, that makes life easier.

- install a list of packages, that I store in my [.dotfiles](https://github.com/msegers/dotfiles), both from AUR and regular.
- installing packages from AUR requires the checkout and make of `yay` so that's automated too.
- setup some things in `.bashrc` so I don't have to remember all the websites where to copy examples from :D.
- Adding my neovim setup as a git submodule (I also use this on my work MacBook, and those dotfiles don't match).

This setup was actually easier than expected, I could have done things even faster using `yay`, but using pacman you can simply fetch your installed packages using `pacman -Qe` which lists only your manually installed packages (not dependencies).
I should have used `yay`, as it actually supports listing packages which are AUR only so it'd be easier to generate two lists.
So basically I did things like `pacman -Qe >> package_list.txt` into my dotfiles. Then in the install script: ``pacman -S `cat package_list.txt` --needed`` which will install what's in the file if it's not installed yet.

The rest of the script is some simple bash logic that makes directories or manipulates files, I don't frequently do this so it's a fun little exercise.

For now the script is quite simple, I'll keep expanding the script if I install things I would always need/want. And hopefully I can find some motivation to make a variation for laptops, I don't think I need to consider different drivers, as `archinstall` already handles graphics drivers etc. (my laptop is intel/nvidia, desktop AMD/AMD).
