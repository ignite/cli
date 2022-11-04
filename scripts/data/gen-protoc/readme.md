# Gen protoc

## Some background on the current approach

Currently, the protocolbuffers/protobuf binaries are only building releases for linux-amd64, linux-arm64 and osx-amd64.

There is an issue on their repo that seems to suggest that osx-arm64 builds are coming soon: <https://github.com/protocolbuffers/protobuf/issues/9397>

When prebuilt binaries exist, we can use the script outlined at the bottom of this file instead of the current version (scripts/gen-protoc)

Until then, the current script downloads the latest version and builds from source. This is triggered in the .github/workflows/gen_protoc.yml file (4 separate jobs, one for each architecture combination)

## Versions

This folder has 4 files that are used to keep track of which version was the latest to be built. This is done to not have to rebuild (which is slow) binaries every time.

The binaries are also non-deterministic (in that they get different file hashes every time you build), so building them on push or a schedule would produce new PR's every time.

# The script we can use later

As mentioned above, when the osx-arm64 binary is being released, we can simplify our approach significantly by having a single job in .github/workflows/gen_protoc.yml that uses the script below:

```bash
#!/bin/bash

# Downloads latest protoc libraries, unpacks and puts them in the right place

set -e

[[ $(command -v wget) ]] || { echo "'wget' not found!" ; dep_check="false" ;}
[[ $(command -v unzip) ]] || { echo "'unzip' not found!" ; dep_check="false" ;}
[[ $(command -v jq) ]] || { echo "'jq' not found!" ; dep_check="false" ;}

[[ ${dep_check} = "false" ]] && { echo "Some dependencie(s) isn't installed yet. Please install that dependencie(s)" ; exit 1 ;}

gh_protoc_release_url="https://api.github.com/repos/protocolbuffers/protobuf/releases/latest"
setdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)" # this line powered by stackoverflow

# Check dir else create save dir
if [[ $(basename ${setdir}) = "scripts" ]] ; then
    if [[ $(basename $(dirname "${setdir}")) = "starport" ]] ; then
        [[ -d $(dirname "${setdir}/pkg/protoc/data") ]] || mkdir -p "$(dirname "${setdir}")/starport/pkg/protoc/data"
    else
        echo "Attention: you are running the script out of the startport project please run it this script in: https://github.com/tendermint/starport"
        exit 1
    fi
else
    echo "$setdir"
    echo "Attention: you are running the script out of the startport project please run it this script in: https://github.com/tendermint/starport"
    exit 1
fi

# Check and Create Temp Directory
[[ -d "/tmp/${0}" ]] && rm -rf "/tmp/${0}"
mkdir -p "/tmp/${0}" && cd "/tmp/${0}"

# Fetch releases, go through assets (release artifacts) and find the relevant ones
wget -O - ${gh_protoc_release_url} \
  | jq --raw-output '.assets[] | select(.name | test("protoc-.*-(linux-x86_64|linux-aarch_64|osx-x86_64)\\.zip")) | .browser_download_url' > filesToDownload.txt

mkdir downloads
wget -P downloads -i filesToDownload.txt

cd downloads
for f in *.zip; do
  name=${f%.zip}
  unzip "$f" -d "$name"


  case "$name" in
  *"linux-x86_64") fname="protoc-linux-amd64"
    ;;
  *"linux-aarch_64") fname="protoc-linux-arm64"
    ;;
  *"osx-x86_64") fname="protoc-darwin-amd64"
    ;;
  *"osx-aarch_64") fname="protoc-darwin-arm64" # TODO: Check that this is actually the name of the osx arm64 binary, it is not released yet
    ;;
  *) echo "No known type was found in $name"; exit 1;
    ;;
  esac

  mv "${f%.zip}/bin/protoc" "$(dirname ${setdir})/starport/pkg/protoc/data/${fname}"
done

echo "/tmp/${0}"
ls -la
rm -rf "/tmp/${0}"
```
