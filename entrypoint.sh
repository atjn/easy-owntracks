#!/bin/sh

STORAGEDIR=/owntracks-storage
CONFIGFILE=${STORAGEDIR}/configuration.conf

AUTHELIA_STORAGEDIR=${STORAGEDIR}/authelia-store
AUTHELIA_SECRETSDIR=${AUTHELIA_STORAGEDIR}/secrets
# These are read directly by Authelia, do not change their names.
# https://www.authelia.com/configuration/methods/secrets/
export AUTHELIA_JWT_SECRET_FILE=${AUTHELIA_SECRETSDIR}/jwt
export AUTHELIA_SESSION_SECRET_FILE=${AUTHELIA_SECRETSDIR}/session
export AUTHELIA_STORAGE_ENCRYPTION_KEY_FILE=${AUTHELIA_SECRETSDIR}/storage_encryption_key

RECORDER_STORAGEDIR=${STORAGEDIR}/recorder-store
# This is needed to make the recorder publish view urls
export OTR_HTTPPREFIX="https://${SITE_ADDRESS}/"

# Makes sure that '/', '|' and '&' is escaped, allowing a string to work properly in the substituion part of a SED command
# This solution is total crap. It is wild to me that there is no better way to safely pass dirty strings to SED.
alias escape="sed 's_/_\\\\/_g' | sed 's_|_\\\\|_g' | sed 's_&_\\\\&_g'"

# Configs
if ! [ -f ${CONFIGFILE} ]; then
    echo "setup: No configuration file found, creating a new one"
    cp -f /configs/configuration.default.template.conf ${CONFIGFILE}
	NEW_PASSWORD=$(xkcdpass --acrostic='track' --delimiter='-' | escape)
	sed -i "s|INSERT_NEW_PASSWORD|${NEW_PASSWORD}|g" ${CONFIGFILE}
fi
echo "setup: You can see and edit the login details for the server in ${CONFIGFILE}"

cd /configs
source ${CONFIGFILE}

# Ensure config is valid


echo "setup: Generating Caddy configuration"
cp -f Caddyfile.template Caddyfile
sed -i "s|INSERT_SITE_ADDRESS|${SITE_ADDRESS}|g" Caddyfile

echo "setup: Generating Authelia configuration"
cp -f authelia.template.yml authelia.yml
sed -i "6i default_redirection_url: https://${SITE_ADDRESS}/" authelia.yml
sed -i "s|INSERT SITE_ADDRESS|${SITE_ADDRESS}|g" authelia.yml

echo "setup: Generating recorder configuration"
cp -f recorder.template.conf recorder.conf
sed -i "3i SITE_ADDRESS = \"${SITE_ADDRESS}\"" recorder.conf
cd /

# Authelia store
if ! [ -d ${AUTHELIA_SECRETSDIR} ]; then
    mkdir -p ${AUTHELIA_SECRETSDIR}
fi
alias generate_simple_secret="tr -cd '[:alnum:]' < /dev/urandom | fold -w 64 | head -n 1 | tr -d '\n'"
cd ${AUTHELIA_SECRETSDIR}
if ! [ -f jwt ]; then
    echo "setup: No Authelia JWT secret found, creating a new one"
    generate_simple_secret > jwt
fi
if ! [ -f session ]; then
    echo "setup: No Authelia session secret found, creating a new one"
    generate_simple_secret > session
fi
if ! [ -f storage_encryption_key ]; then
    echo "setup: No Authelia storage encryption key found, creating a new one"
    generate_simple_secret > storage_encryption_key
fi
cd /configs
echo "setup: Generating Authelia user configuration"
WEB_HASHED_PASSWORD=$(authelia hash-password "${WEB_PASSWORD}" | sed "s|^Password hash: ||" | escape)
cp -f authelia-users.template.yml authelia-users.yml
sed -i "s|INSERT_USERNAME|${WEB_USERNAME}|g" authelia-users.yml
sed -i "s|INSERT_HASHED_PASSWORD|${WEB_HASHED_PASSWORD}|g" authelia-users.yml
cd /

# Recorder store
if ! [ -f ${RECORDER_STORAGEDIR}/ghash/data.mdb ]; then
    echo "setup: No recorder data store found, setting up a new one"
    ot-recorder --initialize
	mkdir -p ${RECORDER_STORAGEDIR}/last
fi

cd /configs

caddy run &
authelia --config /configs/authelia.yml &
ot-recorder
