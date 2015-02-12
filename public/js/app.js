var postTemplate = _.template("<li><a class=\"u-full-width\" href=\"<%= href %>\"><%= title %></a> <br /><small>(<a href=\"<%= comments_href %>\"><%= time %> ago, <%= comments %> comments</a>)</small></li>");
var jobTemplate = _.template("<li><a href=\"<%= href %>\"><%= title %></a> <br /><small>(<%= time %> ago)</small></li>");

$.getJSON("/posts/1", function(data) {
    console.log(data);
    posts = _.map(data, function(post) {
        if (post.comments == 0 && post.points == 0) {
            return jobTemplate(post);
        } else {
            post.comments_href = 'http://news.ycombinator.com/item?id=' + post.id;
            return postTemplate(post);
        }
    });
    $("#items").append(posts);
});
