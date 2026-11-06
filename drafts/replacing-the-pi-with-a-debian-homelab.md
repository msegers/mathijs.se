+++
title = 'Replacing the Pi with a Debian homelab'
date = '2026-11-06T22:27:58+01:00'
draft = true
tags = ['infrastructure', 'linux']
summary = 'Moving my homelab from a Raspberry Pi to a used EliteDesk running Debian and Ansible.'
+++

My Raspberry Pi finally has a replacement. After previously setting up the Pi with Terraform and K3s, I've gone in a slightly different direction: a second-hand office PC, Debian 13, Ansible, and Docker Compose.

The machine is an HP EliteDesk 800 G5 SFF, bought for approximately 2,500 SEK. It now runs my website, Pi-hole, and a Valheim dedicated server, with a small Rust dashboard for managing the latter. There is also a disk mirror that will eventually have something more interesting to do than exist.

This is a hobby project. Unfortunately, putting the household's DNS on it means that not everyone gets to choose when they participate.

## Less cluster, more Debian

The Pi setup involved Terraform, K3s, Helm, and eventually Flux. That was a learning exercise, but this time I'm using Ansible to configure a headless Debian host and Compose to run services. The old Kubernetes configuration is still in the homelab repository; it isn't the architecture of the new machine.

Why not TrueNAS? If the main job were serving disks through a storage appliance, that would be a reasonable direction. But storage is only part of what I'm building. There is a website, DNS, a game server, and custom tooling. Debian gives me a general-purpose Linux host for those, without making the NAS the centre of the setup.

Ansible keeps the host configuration alongside the service definitions. Packages, systemd units, Compose templates, and storage configuration are things I can inspect and change in the repository rather than remember from a sequence of SSH sessions. It doesn't make mistakes impossible. It does make them repeatable, which is occasionally useful.

I'm not claiming this is less work than a dedicated NAS OS. It leaves more of the integration to me. For this project, working on that integration is part of the point.

## First, secure the disks. Physically.

The EliteDesk has an Intel Core i5-8500, 16 GB of Corsair Vengeance DDR4, and a Samsung 850 EVO 500 GB SSD for Debian. The memory is two 8 GB sticks rated for 2666 MT/s, currently running at 2133 MT/s. That's the setup I have, not a tuning success story.

There are also two 1 TB hard drives: a WD Blue and a Seagate Constellation enterprise drive. They arrived loose inside the case, without proper mounting. Before configuring storage, I had to obtain the correct mounting screws. Infrastructure as code still has some gaps.

Then one drive disappeared because of a broken SATA cable. An additional troublesome 4 TB drive prevented the computer from booting at all. Those were separate problems: the cable explained the disappearing drive; the 4 TB disk was associated with the boot failure. I don't have a confirmed explanation for why that disk stopped the machine booting, so I'm not going to assign one.

The two 1 TB disks had rather different histories. The WD Blue reported roughly 36,000 power-on hours, while the Seagate reported around nine. Both passed extended SMART tests. That's reassuring, but a passing test isn't a promise about how long either drive will last. Neither is the word “enterprise”.

## The services, and their dependencies

The main Compose stack lives under `/opt/homelab` and contains Pi-hole, the website, and Caddy. The website is still a small Go-generated static site, served from a container. Caddy sits in front of it and handles HTTPS.

For certificates, Caddy uses Cloudflare DNS challenges through a custom build containing the Cloudflare DNS module. Proving control through DNS also lets the private services have valid HTTPS certificates without opening those services to the internet.

`lab.mathijs.se` is the private Home Tools entry point, including the Valheim dashboard, and `pihole.mathijs.se` provides Pi-hole's admin interface. Pi-hole supplies their local DNS records. Caddy restricts access to the configured LAN subnet; having a certificate does not make an application public.

Valheim has its own Compose definition. The custom Rust dashboard runs as a systemd service on the host, behind Caddy, rather than as another container in the main stack. Its web process doesn't get access to the Docker socket; a separate broker handles the permitted operations.

That distinction matters because the dashboard needs to know which proxy to trust. Recreating Caddy changed its Docker IP, and the dashboard's trusted-proxy configuration still referred to the old address. The private dashboard broke until Ansible updated it.

The role inspects the configured Caddy container and trusts its exact private IPv4 addresses, not the entire Docker subnet. That's deliberately narrow, but it also means replacing the proxy has a configuration consequence outside its own container. Compose service names solve some addressing problems, not all of them.

## Running is not the same as reachable

The most irritating software problem has been recurring networking failures after reboot with Docker 29.9.3. Sometimes containers were running without a usable network attachment or without their published ports.

That is an awkward failure mode when you're looking at a list of running containers and wondering why nothing answers. The intended port bindings in the configuration aren't proof that Docker actually published them.

Pi-hole made this particularly noticeable. With the household using it for DNS, Pi-hole going down effectively broke name resolution for everyone. The internet connection hadn't necessarily failed, but that distinction was not especially helpful to somebody trying to use it.

There is a plausible boot-order concern here: Docker publishes ports on a particular LAN address, and this host uses Wi-Fi with ifupdown. The repository now includes safeguards to wait for that interface and check that the intended address exists before Docker starts. It also has read-only verification of actual network endpoints and published ports.

Those changes address a concrete dependency and detect the symptoms. They do **not** conclusively prove the root cause of every failure, establish a particular upstream Docker bug, or demonstrate that the problem is permanently fixed.

For investigating it, these are more useful than another glance at a green status indicator:

```sh
sudo journalctl -b -u docker.service
sudo docker network inspect homelab_default
```

I've managed to make DNS depend on the same machine I'm experimenting with. That's a compromise worth stating plainly, rather than describing this as reliable household infrastructure.

## A mirror, not a backup

The two 1 TB disks form an `mdadm` RAID 1 mirror, with ext4 mounted at `/srv/storage`. Usable capacity is roughly one disk: about 931 GiB before filesystem overhead, not 2 TB. Debian remains on the separate SSD.

Ansible identifies the drives through stable `/dev/disk/by-id` paths and checks their identities before initialization. The destructive setup is explicitly gated, rather than something an ordinary deployment should repeat. Formatting the wrong disk would be a fairly effective demonstration of automation.

The initial resync still had to run even though the disks were practically empty. RAID works below the filesystem: it doesn't treat “there are no files I care about” as evidence that both disks contain matching blocks. I could watch its progress with:

```sh
cat /proc/mdstat
```

RAID 1 protects against a disk failure. It does not protect against accidental deletion, and it does not remove the need for backups. Delete a file and the mirror faithfully reflects that decision on both disks.

Samba and Immich are plans, with Jellyfin a possibility. They aren't already using this storage, and building the mirror hasn't quietly solved backups for them either.

For now, the Pi's replacement is a Debian box serving a few useful things, alongside some dependencies I've learned about by breaking them. Before adding more services, I'd like a reboot to be a less interesting event.
