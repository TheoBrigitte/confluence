# confluence

Confluence is a client server torrent streaming solution.

All the torrent operations happens on the server side, which means it is only seens as a streaming solution from the client.

This repository holds both the server and client code.

## Know issues

#### OpenSubtitles API

401 Unauthorized: even if opensubtitles [documentation](http://trac.opensubtitles.org/projects/opensubtitles/wiki/XMLRPC#LogIn) claims to support anonymous login, it might reject anonymous login with 01 errors. In this case use the `-osUser` and `-osPassword` to provide credentials.

414 Unknown User Agent: this means the user-agent in use as to be changed. At the time I am writing those line the current user-agent used by github.com/oz/osdb is "osdb-go 0.2". There might be other working user-agent out there like "Subdownloader 1.2.4" from [there](http://trac.opensubtitles.org/projects/opensubtitles/wiki/DevReadFirst#Howtorequestanewuseragent).
