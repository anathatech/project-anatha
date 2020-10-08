#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/anathad/${BINARY:-anathad-manager}
ID=${ID:-0}
LOG=${LOG:-anathad.log}
LOG_LEVEL="main:info,state:info,x/crisis:info,x/hra:info,x/upgrade:info,x/gov:info,x/governance:info,x/treasury:info,x/distribution:debug,x/mint:debug,x/astaking:debug,*:error"

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'anathad' E.g.: -e BINARY=anathad_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

export DAEMON_NAME=anathad
export DAEMON_HOME="/anathad/node${ID}/anathad"
export DAEMON_ALLOW_DOWNLOAD_BINARIES=on
export DAEMON_RESTART_AFTER_UPGRADE=on

mkdir -p ${DAEMON_HOME}/upgrade_manager/genesis/bin
mkdir -p ${DAEMON_HOME}/upgrade_manager/upgrades

if ! [ -f "${DAEMON_HOME}/upgrade_manager/genesis/bin/${DAEMON_NAME}" ]; then
	cp /anathad/${DAEMON_NAME} ${DAEMON_HOME}/upgrade_manager/genesis/bin/${DAEMON_NAME}
fi

##
## Run binary with all parameters
##
export ANATHADHOME="/anathad/node${ID}/anathad"

if [ -d "$(dirname "${ANATHADHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${ANATHADHOME}" "$@" --log_level  "${LOG_LEVEL}" | tee "${ANATHADHOME}/${LOG}"
else
  "${BINARY}" --home "${ANATHADHOME}" "$@" --log_level  "${LOG_LEVEL}"
fi
