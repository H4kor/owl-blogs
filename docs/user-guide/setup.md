# Setup

## Create an author

Most configuration of owl-blogs is done on the website, but you need an account to access the admin pages. An author account has to be created via the CLI with:

```
owl new-author -u <name> -p <password>
```

## Basic Configuration

Open the blog in your browser. In the footer you can find a link to **Editor**.
Alternatively you can access the [/admin](http://localhost:3000/admin) page directly.
Login with the author you just created.

Open the [Site Settings](http://localhost:3000/site-config).
First set the "FullUrl" to the correct URL of your blog, e.g. `https://blog.example.com`.
Multiple features of owl-blogs depend on this URL, e.g to generate correct links.

While you're at it you can also change the title and color of the blog.
The primary color should be dark enough for white text be easily readable on this color.

At the bottom of the site settings you can select the types of entries which should be shown on the main page of your blog.
I'd recommend to select all except "Note" and "Page".
You can learn more about the entry types in the [First Post Guide](first_post.md)