This is the quickest and dirtiest of podcast feed generators.

If you really want it,

    go get github.com/joeshaw/podfeed

Usage:

    podfeed "Feed Title" "Feed Description" "http://example.com/some-podcast" *.mp3 > rss.xml

It's expected that the files are accessible on the web with the URL + "/" +
the file name (although podfeed will do the proper URL escaping for you).

The dates are pulled from the mtime of the file, so you might need to
adjust them if the ordering is wrong.

I've tested this in Overcast, but no other podcast players.
