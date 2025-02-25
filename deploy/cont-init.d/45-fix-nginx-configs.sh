#!/command/with-contenv sh

sed -Eie 's|\$GLOBAL_HOST|'$GLOBAL_HOST'|g;s|\$GLOBAL_HOST_CONSOLE|'$GLOBAL_HOST_CONSOLE'|g' /etc/nginx/nginx.conf /etc/nginx/static/apocryphLogin.html

echo 'daemon off;' >> /etc/nginx/nginx.conf
