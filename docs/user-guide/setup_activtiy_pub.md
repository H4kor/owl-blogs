# Setup ActivityPub

Owl-blogs will also publish your posts via ActivityPub.
By default the user name of your blog will be `@blog@<your.domain.com>`.

In the Module Configuration section of the admin page, click the "activity_pub" button.
Here you can change the the user name (`blog` by default) to any other name you prefer. 
The public and private keys don't have to be changed and should not be touched for the moment.

To test the ActivityPub integration open any ActivityPub app, such as [Mastodon](https://joinmastodon.org/).
Search for your blogs user name (e.g. `@blog@example.com`).
Notice that existing posts will potentially no be listed in the app, as the server you are using does not retrieved posts if no one follows the blog.
Follow your blog to fix this.

In the "Followers" admin view of your blog you should see your account listed now.
The next post will show up in you feed and your blog's profile page when opened.