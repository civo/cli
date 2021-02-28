#!/bin/sh
set -e

# Copyright Civo Ltd 2020, all rights reserved

export VERIFY_CHECKSUM=0
export OWNER="civo"
export REPO="cli"
export SUCCESS_CMD="$OWNER version"
export BINLOCATION="/usr/local/bin"

###############################
# Get the last version        #
###############################

get_last_version() {
  VERSION=""

  echo "Finding latest version from GitHub"
  VERSION=$(curl -sI https://github.com/$OWNER/$REPO/releases/latest | grep -i "location:" | awk -F"/" '{ printf "%s", $NF }' | tr -d '\r')
  VERSION_NUMBER=$(echo "$VERSION" | cut -d "v" -f 2)
  echo "$VERSION_NUMBER"

  if [ ! "$VERSION" ]; then
    echo "Failed while attempting to install $REPO. Please manually install:"
    echo ""
    echo "1. Open your web browser and go to https://github.com/$OWNER/$REPO/releases"
    echo "2. Download the latest release for your platform. Call it '$REPO'."
    echo "3. chmod +x ./$REPO"
    echo "4. mv ./$REPO $BINLOCATION"
    if [ -n "$ALIAS_NAME" ]; then
      echo "5. ln -sf $BINLOCATION/$REPO /usr/local/bin/$ALIAS_NAME"
    fi
    exit 1
  fi
}

###############################
# Check for curl              #
###############################
hasCurl() {
  which curl
  if [ "$?" = "1" ]; then
    echo "You need curl to use this script."
    exit 1
  fi
}

# --- set arch and suffix, fatal if architecture not supported ---
setup_verify_arch() {
  if [ -z "$ARCH" ]; then
    ARCH=$(uname -m)

  fi
  case $ARCH in
  amd64)
    ARCH=-amd64
    SUFFIX=
    ;;
  x86_64)
    ARCH=-amd64
    SUFFIX=
    ;;
  arm64)
    ARCH=-arm64
    ;;
  aarch64)
    ARCH=-arm64
    ;;
  arm*)
    ARCH=-arm
    SUFFIX=
    ;;
  *)
    fatal "Unsupported architecture $ARCH"
    ;;
  esac
}

setup_verify_os() {
  if [ -z "$SUFFIX" ]; then
    SUFFIX=$(uname -s)
  fi
  case $SUFFIX in
  "Darwin")
    SUFFIX="-darwin"
    ;;
  "MINGW"*)
    SUFFIX="windows"
    ;;
  "Linux")
    SUFFIX="-linux"
    ;;
  *)
    fatal "Unsupported OS $SUFFIX"
    ;;
  esac
}

download() {
  URL=https://github.com/$OWNER/$REPO/releases/download/$VERSION/$OWNER-$VERSION_NUMBER$SUFFIX$ARCH.tar.gz
  TARGETFILE="/tmp/$OWNER-$VERSION_NUMBER$SUFFIX$ARCH.tar.gz"
  echo "Downloading package $URL to $TARGETFILE"

  curl -sSL "$URL" --output "$TARGETFILE"

  if [ "$VERIFY_CHECKSUM" = "1" ]; then
    check_hash
  fi

  tar -xf /tmp/$OWNER-$VERSION_NUMBER$SUFFIX$ARCH.tar.gz -C /tmp
  chmod +x /tmp/$OWNER

  echo "Download complete."

  if [ ! -w "$BINLOCATION" ]; then

    echo
    echo "============================================================"
    echo "  The script was run as a user who is unable to write"
    echo "  to $BINLOCATION. To complete the installation the"
    echo "  following commands may need to be run manually."
    echo "============================================================"
    echo
    echo "  sudo mv /tmp/civo $BINLOCATION/$OWNER"

    if [ -n "$ALIAS_NAME" ]; then
      echo "  sudo ln -sf $BINLOCATION/$OWNER $BINLOCATION/$ALIAS_NAME"
    fi

    echo

  else

    echo
    echo "Running with sufficient permissions to attempt to move $OWNER to $BINLOCATION"

    if [ ! -w "$BINLOCATION/$OWNER" ] && [ -f "$BINLOCATION/$OWNER" ]; then

      echo
      echo "================================================================"
      echo "  $BINLOCATION/$OWNER already exists and is not writeable"
      echo "  by the current user.  Please adjust the binary ownership"
      echo "  or run sh/bash with sudo."
      echo "================================================================"
      echo
      exit 1

    fi

    mv /tmp/$OWNER $BINLOCATION/$OWNER

    if [ "$?" = "0" ]; then
      echo "New version of $OWNER installed to $BINLOCATION"
    fi

    if [ -e TARGETFILE ]; then
      rm TARGETFILE
      rm /tmp/$OWNER
    fi

    ${SUCCESS_CMD}
  fi

}

check_hash() {
  SHACMD="sha256sum"

  if [ ! -x "$(command -v $SHACMD)" ]; then
    SHACMD="shasum -a 256"
  fi

  if [ -x "$(command -v "$SHACMD")" ]; then

    TARGETFILEDIR=${TARGETFILE%/*}

    (cd "$TARGETFILEDIR" && curl -sSL https://github.com/$OWNER/$REPO/releases/download/$VERSION/$OWNER-$VERSION_NUMBER-checksums.sha256 | $SHACMD -c >/dev/null)

    if [ "$?" != "0" ]; then
      rm TARGETFILE
      echo "Binary checksum didn't match. Exiting"
      exit 1
    fi
  fi
}

# Error: Show error message in red and exit
fatal() {
  printf "Error: \033[31m${1}\033[39m\n"
  exit 1
}


{
  hasCurl
  setup_verify_arch
  setup_verify_os
  get_last_version
  download
}
