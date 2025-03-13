#!/usr/bin/env bash

set -e

LAUNCH_URL="http://s3-aapp.local/app/"
FIREFOX_PROFILE="/tmp/apocryph-firefox-profile" ##$(mktemp -d)
mkdir -p "$FIREFOX_PROFILE"

if [ ! -f "$FIREFOX_PROFILE/prefs.js" ]; then

  cat >"$FIREFOX_PROFILE/prefs.js" <<EOF
user_pref("network.proxy.socks", "localhost");
user_pref("network.proxy.socks_port", 1080);
user_pref("network.proxy.type", 1);
user_pref("extensions.ui.lastCategory", "addons://list/extension");
user_pref("trailhead.firstrun.didSeeAboutWelcome", true);
user_pref("browser.aboutwelcome.didSeeFinalScreen", true);
user_pref("browser.bookmarks.restore_default_bookmarks", false);
user_pref("browser.newtabpage.activity-stream.feeds.topsites", false);
user_pref("toolkit.telemetry.reportingpolicy.firstRun", false);
user_pref("toolkit.telemetry.enabled", false);
user_pref("toolkit.telemetry.rejected", true);
EOF

  mkdir "$FIREFOX_PROFILE/extensions/"
  curl https://addons.mozilla.org/firefox/downloads/file/4392249/ether_metamask-12.6.2.xpi -L -o "$FIREFOX_PROFILE/extensions/webextension@metamask.io.xpi"
  
  firefox --profile "$FIREFOX_PROFILE" "$LAUNCH_URL" &
  sleep 5
  firefox --profile "$FIREFOX_PROFILE" "about:addons"
  wait
else
  firefox --profile "$FIREFOX_PROFILE" "$LAUNCH_URL"
fi

