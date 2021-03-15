#!/bin/bash

package="github.com/wttw/markdownish/cmd/mdtohtml"
platforms=("windows/amd64" "darwin/amd64" "linux/amd64")

# ----------
# Error handling
# ----------

die() {
    echo "^^^ +++"
    echo "${BASH_SOURCE[1]}: line ${BASH_LINENO[0]}: ${FUNCNAME[1]}: ${1-Died}" >&2
    exit 1
}

set -o pipefail -o noclobber -o nounset

# ----------
# Find our source tree
# ----------

# The base directory of the checkout is where this script lives
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

cd ${BASEDIR}

rm -rf builds
mkdir builds || die
cd builds || die

package_split=(${package//\// })
package_name=${package_split[-1]}

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_name=$package_name
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    dir="${GOOS}-${GOARCH}"
    mkdir ${dir} || die
    
    env GOOS=$GOOS GOARCH=$GOARCH go build -o "${dir}/${output_name}" ${package}
    if [ $? -ne 0 ]; then
        die 'Build failed'
    fi
    cd "${dir}" || die
    zip -9 "../${dir}.zip" ${output_name} || die "Zip failed"
    cd "${BASEDIR}/builds" || die
done
