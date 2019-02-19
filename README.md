# confluence

Confluence is a client server torrent streaming solution.

All the torrent operations happens on the server side, which means it is only seens as a streaming solution from the client.

This repository holds both the server and client code.

## Example

magnet: magnet:?xt=urn:btih:3945AAE317B45183150400E128D007806EB4CEA1&dn=Conan+the+Barbarian+%281982%29+%5B1080p%5D+%5BYTS.AM%5D&tr=udp%3A%2F%2Fglotorrents.pw%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fp4p.arenabg.ch%3A1337&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337
magnet: magnet:?xt=urn:btih:00102086B401F8CE049BE55410FF9C69D87BB740&dn=Deadpool+2+%282018%29+%5B720p%5D+%5BYTS.AM%5D&tr=udp%3A%2F%2Fglotorrents.pw%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fp4p.arenabg.ch%3A1337&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337


## Know issues

#### OpenSubtitles API

401 Unauthorized: even if opensubtitles [documentation](http://trac.opensubtitles.org/projects/opensubtitles/wiki/XMLRPC#LogIn) claims to support anonymous login, it might reject anonymous login with 01 errors. In this case use the `-osUser` and `-osPassword` to provide credentials.

414 Unknown User Agent: this means the user-agent in use as to be changed. At the time I am writing those line the current user-agent used by github.com/oz/osdb is "osdb-go 0.2". There might be other working user-agent out there like "Subdownloader 1.2.4" from [there](http://trac.opensubtitles.org/projects/opensubtitles/wiki/DevReadFirst#Howtorequestanewuseragent).
