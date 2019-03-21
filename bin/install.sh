#!/usr/bin/env bash
NAME="stufy"
REPO_NAME="stufy"
OS=""
OWNER="ArthurHlt"
: "${TMPDIR:=${TMP:-$(CDPATH=/var:/; cd -P tmp)}}"
cd -- "${TMPDIR:?NO TEMP DIRECTORY FOUND!}" || exit
cd -
echo "Installing ${NAME}..."
if [[ "$OSTYPE" == "linux-gnu" || "$(uname -s)" == "Linux" ]]; then
    OS="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="darwin"
elif [[ "$OSTYPE" == "cygwin" ]]; then
    OS="windows"
elif [[ "$OSTYPE" == "msys" ]]; then
    OS="windows"
elif [[ "$OSTYPE" == "win32" ]]; then
    OS="windows"
else
    echo "Os not supported by install script"
    exit 1
fi
VERSION=$(curl -s https://api.github.com/repos/${OWNER}/${REPO_NAME}/releases/latest | grep tag_name | head -n 1 | cut -d '"' -f 4)
ARCHNUM=`getconf LONG_BIT`
ARCH=""
CPUINFO=`uname -m`
if [[ "$ARCHNUM" == "32" ]]; then
    ARCH="386"
else
    ARCH="amd64"
fi
if [[ "$CPUINFO" == "armv5"* ]]; then
    ARCH="armv5"
fi
if [[ "$CPUINFO" == "armv6"* ]]; then
    ARCH="armv6"
fi
if [[ "$CPUINFO" == "armv7"* ]]; then
    ARCH="armv5"
fi
if [[ "$CPUINFO" == "arm64"* ]]; then
    ARCH="arm64"
fi
FILENAME="${NAME}_${OS}_${ARCH}"
if [[ "$OS" == "windows" ]]; then
    FILENAME="${FILENAME}.exe"
fi
LINK="https://github.com/${OWNER}/${NAME}/releases/download/${VERSION}/${FILENAME}"
if [[ "$OS" == "windows" ]]; then
    FILEOUTPUT="${FILENAME}"
else
    FILEOUTPUT="${TMPDIR}/${FILENAME}"
fi
RESPONSE=200
if hash curl 2>/dev/null; then
    RESPONSE=$(curl --write-out %{http_code} -L -o "${FILEOUTPUT}" "$LINK")
else
    wget -o "${FILEOUTPUT}" "$LINK"
    RESPONSE=$?
fi

if [ "$RESPONSE" != "200" ] && [ "$RESPONSE" != "0" ]; then
    echo "File ${LINK} not found, so it can't be downloaded."
    rm "$FILEOUTPUT"
    exit 1
fi

chmod +x "$FILEOUTPUT"
if [[ "$OS" == "windows" ]]; then
    mv "$FILEOUTPUT" "${NAME}"
else
    mv "$FILEOUTPUT" "/usr/local/bin/${NAME}"
fi
echo "${NAME} has been installed."