---
name: new-blog-post
description: Interview the user and scaffold a new blog post for this site, including optional image copying, resizing, and thumbnail generation.
---

# New blog post

Use this skill when the user wants to add, draft, or write a new blog entry for `mathijs.se`.

## Interview first

Do not create or edit files until the interview is complete.

Ask one question at a time. Use an interactive choice prompt if your agent supports one; otherwise ask in plain text and wait for the answer. Never invent facts, tags, dates, or image alt text. If the user does not know, propose a sensible default and ask them to confirm it.

Ask in this order:

1. **Title** — the working title.
2. **Summary** — one sentence for the blog listing, RSS, and previews. Offer to derive it from their opening paragraph.
3. **Content** — ask for notes, an outline, or the full text. Preserve their wording; do not pad or rewrite it unless they ask you to.
4. **Tags** — show existing tags first:
   ```sh
   grep -rh '^tags' content | sed -E "s/.*\[(.*)\].*/\1/" | tr ',' '\n' | tr -d " '[]" | sort -u
   ```
   Offer existing tags and allow new ones.
5. **Image** — ask whether they want one. If yes, ask for an absolute path on their filesystem and for the alt text. If no, continue without one.
6. **Draft or published?** — default to a draft. Only set `draft = false` when they explicitly say publish.
7. **Confirm** — show the derived date, slug, path, tags, and draft state before writing.

## Create the post

1. Get today's local date:
   ```sh
   date +%Y-%m-%dT%H:%M:%S%:z
   ```
2. Derive a slug from the title: lowercase, replace anything that is not a letter or digit with `-`, collapse dashes, trim dashes, and keep it under about 60 characters.
3. Use this path:
   ```text
   content/blog/YYYY/MM/<slug>/index.md
   ```
4. If the user supplied an image, run the image helper from the repository root:
   ```sh
   ./.agents/skills/new-blog-post/scripts/prepare-image.sh "/absolute/path/to/image" "content/blog/YYYY/MM/<slug>"
   ```
   It writes `image.<ext>` and `thumb.<ext>`. The full image is for the post body; the thumbnail is used in the blog listing.
5. Write the front matter and body. Use double quotes for any value containing an apostrophe:

   ```toml
   +++
   title = 'The title'
   date = '2025-09-20T13:00:46+02:00'
   draft = true
   tags = ['linux', 'learning']
   summary = 'One sentence describing the post.'
   image = { src = 'image.jpg', thumb = 'thumb.jpg', alt = 'What the image shows' }
   +++

   Opening paragraph.

   ![What the image shows](image.jpg)
   *Optional caption.*
   ```

   Omit the `image` line entirely if there is no image.

6. Put the user's content in the body as Markdown. If an image was added, include it in the body with `![alt](image.jpg)` at the point where it belongs.
7. Run the site build and fix any problems:
   ```sh
   go run .
   ```
8. Report the file path and whether the post is a draft or published. If the local server is running, give the preview URL.

## Image helper

`prepare-image.sh` requires `ffmpeg`. It accepts `jpg`, `jpeg`, `png`, and `webp`. It:

- resizes the full image down to at most 1600px on its longest side (never upscales),
- creates a center-cropped 800x450 thumbnail for the blog listing,
- keeps the original format and prints the resulting file names and sizes.

If `ffmpeg` is unavailable, tell the user rather than silently copying a huge file. They can install it or copy the image manually and omit the thumbnail.
