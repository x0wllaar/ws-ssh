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