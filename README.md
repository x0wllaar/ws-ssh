# ws-ssh

*A small program that forwards TCP connections over websockets.*
    
Can be very useful if you want to hide your SSH (or other pure TCP) server
behind Cloudflare or other CDN.

To use, on the server:
run

    ws-ssh listen --from 127.0.0.1:8822 --to 127.0.0.1:22

then add to nginx config:

    location /ws-ssh {
        proxy_pass http://127.0.0.1:8822/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
    }

and restart nginx for the changes to take effect

On the client:

    ssh -o ProxyCommand="ws-ssh connect --url https://yoursite.com/ws-ssh stdio" yoursite.com


It's also recommended to add frequent SSH keepalives to such connections:

    Host yoursite.com
        ServerAliveInterval 10
        ServerAliveCountMax 2

## Building

ws-ssh uses a unnecessarily complex build process in order to integrate 
the browserproxy page and all of its JS and WASM dependencies into the
executable file.

As a result, it will not build if you simply clone the repo and run 
`go build`.

On the first build, or after any change in the `browserproxy` directory,
you need to run `./build.sh` so that it generates the HTML page correctly.

Please note that running `./build.sh` requires UNIX tools and Perl to be
installed and in PATH.

## browserproxy

browserproxy is an experimental connection method that uses a real web
browser to proxy the WebSocket connection, thus completely hiding the
usage of `ws-ssh` from network observers.

This also allows it to use all the new web privacy/security tech like ECH
and DoH to add even more cloaking to the connection.

The reason behind this is that modern browsers that move extremely 
fast with adding ECH, DoH, post-quantum crypto and all those things. It's 
obvious the parrot is dead, and all attempts to emulate modern web browser
TLS will be futile and fingerprintable.

The only way to avoid this is to replace the parrot with a puppet, and
use the real thing to make the connections instead of trying to emulate.
This is an attempt at doing that.

Caveats:
1. The way it currently works by serving the page from localhost, and
   then connecting to a (CDN shielded) remote. This will lead to the browser
   sending an origin header like `http://127.0.0.1:8822` to the remote.
   On Cloudflare it least, this will immediately lead to a refused connection.
   A workaround is to use a browsr extension like
   [ModHeader](https://addons.mozilla.org/en-US/firefox/addon/modheader-firefox/)
   to remove the origin header from requests to the ws-ssh remote URL.
   It should not add any distinguishers to outside observers, but requires
   Firefox, as the header cannot be removed in Chrome.
   A more permament way to solve this is for each remote to serve their connection
   page from their origin, and have it connect to localhost instead, but this is not
   done/tested yet.

   
