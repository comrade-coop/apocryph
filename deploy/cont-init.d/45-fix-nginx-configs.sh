#!/command/with-contenv sh

sed -i -Ee 's|\$GLOBAL_HOST_CONSOLE|'$GLOBAL_HOST_CONSOLE'|g;s|\$GLOBAL_HOST|'$GLOBAL_HOST'|g' /etc/nginx/nginx.conf /etc/nginx/sites-enabled/default /etc/nginx/static/apocryphLogin/index.html

echo 'daemon off;' >> /etc/nginx/nginx.conf
