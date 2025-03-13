#!/command/with-contenv bash

set -e
shopt -s globstar

# HACK

if [ "$FIXUP_VITE_STORAGE_SYSTEM" == "true" ]; then
  . /shared_secrets/backend_key.env

  VITE_STORAGE_SYSTEM=$(apocryph-s3-backend get-public-key)
  
  echo "Fixing-up storage system key in frontend!"

  sed -i -Ee 's|\$\$\$VITE_STORAGE_SYSTEM\$\$\$|'$VITE_STORAGE_SYSTEM'|g' /usr/share/nginx/html/**/*.js
fi
