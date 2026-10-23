# Mathijs.se

My blog. A small Go program turns the Markdown files in `content/` into static HTML in `public/`.

## Build

Requirements: Go 1.24+.

```sh
go run .
```

The generated site is written to `public/`. To preview it locally:

```sh
go run . --serve
```

Then open <http://localhost:1313/>.

There is no Hugo theme/module setup, no Node toolchain, and no generated-file state kept between builds.

## Writing

- Site title, URL and navigation live in `site.json`.
- Posts live in `content/blog/YYYY/MM/slug/index.md`.
- Pages live in `content/<name>/index.md`.
- An `_index.md` marks a section overview page.
- Set `draft = true` in the front matter to leave a page out of the build.
- Images can sit next to a post's `index.md` and be referenced by filename.
- `image.thumb` is optional. When set, the blog listing uses the thumbnail while the post uses the full `image.src`.
- Files in `static/` are copied to the root of `public/` unchanged.
- HTML layouts are in `templates/`; CSS is in `static/css/main.css`.

### Creating a post with an agent

`.agents/skills/new-blog-post/SKILL.md` is a vendor-neutral agent skill. It interviews you for the title, summary, content, tags, image, and draft state, then creates the post.

The skill can take an image from any path on your filesystem and uses `ffmpeg` to resize it and make a 16:9 thumbnail:

```sh
.agents/skills/new-blog-post/scripts/prepare-image.sh "/path/to/image.jpg" "content/blog/2025/09/my-post"
```


Front matter stays TOML, so existing content works as-is:

```toml
+++
title = 'A post'
date = '2025-09-20T13:00:46+02:00'
draft = false
tags = ['linux']
summary = 'A short description.'
image = { src = 'screenshot.png', alt = 'A screenshot' }
+++
```

Tags and categories become `/tags/...` and `/categories/...` pages automatically. The build also generates `sitemap.xml` and RSS feeds.

## Deploy

The `Dockerfile` builds the site with Go and serves the result with nginx:

```sh
docker build -t mathijs-se .
docker run --rm -p 8080:80 mathijs-se
```

`build.sh` builds the site, publishes a multi-architecture Docker image, and bumps the Helm chart version:

```sh
./build.sh          # patch bump
./build.sh minor
./build.sh major
```

It expects Docker buildx with push access to `msegers/mathijs-se`, so run `docker login` first. The script creates a `docker-container` buildx builder automatically, which is required for the multi-architecture image. The chart in `helm/` still deploys that image to Kubernetes.

The Terraform/SSH setup that actually rolls the deployment out to the Pi lives in the separate [`terraform-pi`](https://github.com/msegers/terraform-pi) repository, not here.
