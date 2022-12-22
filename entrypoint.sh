#!/bin/sh

STORAGEDIR=/owntracks-storage

AUTHELIA_STORAGEDIR=${STORAGEDIR}/authelia-store
AUTHELIA_SECRETSDIR=${AUTHELIA_STORAGEDIR}/secrets
# These are read directly by Authelia, do not change their names.
# https://www.authelia.com/configuration/methods/secrets/
export AUTHELIA_JWT_SECRET_FILE=${AUTHELIA_SECRETSDIR}/jwt
export AUTHELIA_SESSION_SECRET_FILE=${AUTHELIA_SECRETSDIR}/session
export AUTHELIA_STORAGE_ENCRYPTION_KEY_FILE=${AUTHELIA_SECRETSDIR}/storage_encryption_key


# Make sure everything is set up to run correctly
setup.sh


# Run all services with fault tolerance

# Ensures that the given service never stops running. If there is an error, it is simply restarted.
persist() {
	until $3; do
		echo -e "\033[31mSYSTEM: $2 just crashed. This is not supposed to happen, please check your logs and see if there is an issue you can fix. The process will restart in $1 seconds.\033[0m" >&2
		sleep $1
	done
}

cd /configs
persist 10 'Caddy webserver' 'caddy run' &
persist 10 'Authelia' 'authelia --config /configs/authelia.yml' &
persist 10 'Owntracks extended API' 'extapi' &
persist 10 'OwnTracks recorder' 'ot-recorder'

