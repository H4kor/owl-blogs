![Mascot](assets/owl.png)

# Owl Blogs

A simple web server for blogs generated from Markdown files

## Repository

A repository holds all data for a web server. It contains multiple users.

## User

A user has a collection of posts.
Each directory in the `/users/` directory of a repository is considered a user.

### User Directory structure

```
<user-name>/
  \- public/
       \- <post-name>
            \- index.md
                -- This will be rendered as the blog post.
                -- Must be present for the blog post to be valid.
                -- All other folders will be ignored
            \- media/
                -- Contains all media files used in the blog post.
                -- All files in this folder will be publicly available
  \- meta/
       \- base.html
            -- The template used to render all sites
       \- VERSION
            -- Contains the version string.
            -- Used to determine compatibility in the future
  \- config.yml
        -- Contains settings global to the user.
        -- For example: page title and style options
```

### Post

Posts are Markdown files with a mandatory metadata head.

- The `title` will be added to the web page and does not have to be reapeated in the body. It will be used in any lists of posts.
- `aliases` are optional. They are used as permanent redirects to the actual blog page.

```
---
title: My new Post
date: 13 Aug 2022 17:07 UTC
aliases:
     - /my/new/post
     - /old_blog_path/
---

Actual post

```
