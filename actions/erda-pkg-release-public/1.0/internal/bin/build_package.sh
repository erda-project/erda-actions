cd "$(dirname "${BASH_SOURCE[0]}")"
set -o errexit -o nounset -o pipefail

## get main.beta version of erda
function get_main_beta_version() {
    local main_beta_version="$1"

    ## filter v in version
    if echo "$main_beta_version" | grep "v" > /dev/null 2>&1; then
        main_beta_version="${main_beta_version#*v}"
    fi

    ## filter -rc in version
    if echo "$main_beta_version" | grep "\-rc" > /dev/null 2>&1; then
        main_beta_version="${main_beta_version%-rc*}"
    fi

    ## filter patch version in version
    dotNum=$(echo "$main_beta_version" | awk -F"." '{print NF-1}')
    if [ "$dotNum" == "2" ]; then
        main_beta_version="${main_beta_version%.*}"
    fi

    echo "$main_beta_version"
}

## dir to storage installing package of erda
rm -rf package
mkdir package

## ERDA_VERSION validate
if ! env | grep ERDA_VERSION > /dev/null 2>&1; then
    echo "no specify env ERDA_VERSION"
fi
if [[ -z "$ERDA_VERSION" ]]; then
    echo "ERDA_VERSION is empty"
    exit
fi

## erda actions and addons
rm -rf ./erda/scripts/erda-actions
cp -a /tmp/"$ERDA_VERSION"/extensions/erda-actions ./erda/scripts/

rm -rf ./erda/scripts/erda-addons
cp -a /tmp/"$ERDA_VERSION"/extensions/erda-addons ./erda/scripts/

CURRENT_PATH="$PWD"
VERSION_PATH="$CURRENT_PATH"/version
ERDA_YAML="$CURRENT_PATH/erda/erda/templates/erda/erda.yaml"
MAIN_BETA_VERSION=$(get_main_beta_version "$ERDA_VERSION")

## compose erda.yaml
if [ -f "$ERDA_YAML" ]; then
    rm -rf "$ERDA_YAML"
fi
cd "$VERSION_PATH" &&
./compose.sh "$ERDA_VERSION"
cp -rf erda.yaml "$ERDA_YAML"
cd "$CURRENT_PATH" &&
