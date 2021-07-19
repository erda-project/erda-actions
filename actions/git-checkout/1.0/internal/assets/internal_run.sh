#!/bin/bash

set -e

exec 2>&1

source $(dirname $0)/common.sh

destination=.

PATH=/usr/local/bin:$PATH

#payload=$(mktemp $TMPDIR/git-resource-request.XXXXXX)

#cat > $payload <&0

load_pubkey
load_git_crypt_key
configure_https_tunnel
configure_git_ssl_verification
configure_credentials

# 从 环境变量 获取 params
uri=${ACTION_URI-${GITTAR_REPO}} # 默认用 GITTAR_REPO
branch=${ACTION_BRANCH-${GITTAR_BRANCH}} # 默认用 GITTAR_BRANCH

# represents the current execution pipeline is retrying the failed node or retrying the entire process
if [ "$PIPELINE_TYPE" == "rerun-failed" ] || [ "$PIPELINE_TYPE" == "rerun" ]; then
    # gittar address can rerun
    if [ "$uri" == "$GITTAR_REPO" ] && [ "$branch" == "$GITTAR_BRANCH" ]; then
       branch=${GITTAR_COMMIT-${GITTAR_BRANCH}}
    fi
fi

git_config_payload=${ACTION_GIT_CONFIG-[]}
ref=${ACTION_REF-HEAD}
depth=${ACTION_DEPTH-1}
fetch=${ACTION_FETCH}
submodules=${ACTION_SUBMODULES-all}
submodule_recursive=${ACTION_SUBMODULE_RECURSIVE-true}
submodule_remote=${ACTION_SUBMODULE_REMOTE-false}
gpg_keyserver=${ACTION_GPG_KEYSERVER-hkp://ipv4.pool.sks-keyservers.net/}
disable_git_lfs=${ACTION_DISABLE_GIT_LFS-false}
clean_tags=${ACTION_CLEAN_TAGS-false}
short_ref_format=${ACTION_SHORT_REF_FORMAT-%s}

configure_git_global "${git_config_payload}"

if [ -z "$uri" ]; then
  echo "missing uri" >&2
  exit 1
fi

branchflag=""
if [ -n "$branch" ]; then
  branchflag="--branch $branch"
fi

depthflag=""
if test "$depth" -gt 0 2> /dev/null; then
  depthflag="--depth $depth"
fi

isCommit=false

if [ ${#branch} -eq 40 ]; then
    isCommit=true
fi

if [ $isCommit = true ]; then
    # commit checkout
    echo 'fetch mode'
    mkdir -p $destination
    cd $destination
    git init
    git remote add origin $uri
    git fetch --depth 1 origin $branch
    git checkout FETCH_HEAD
else
    echo 'clone mode'
    pwd
    git clone --single-branch $depthflag $uri $branchflag $destination
    cd $destination
    git fetch origin refs/notes/*:refs/notes/*

    # A shallow clone may not contain the Git commit $ref:
    # 1. The depth of the shallow clone is measured backwards from the latest
    #    commit on the given head (master or branch), and in the meantime there may
    #    have been more than $depth commits pushed on top of our $ref.
    # 2. If there's a path filter (`paths`/`ignore_paths`), then there may be more
    #    than $depth such commits pushed to the head (master or branch) on top of
    #    $ref that are not affecting the filtered paths.
    #
    # In either case we try to deepen the shallow clone until we find $ref, reach
    # the max depth of the repo, or give up after a given depth and resort to deep
    # clone.
    if [ "$depth" -gt 0 ]; then
      readonly max_depth=128

      d="$depth"
      while ! git checkout -q $ref &>/dev/null; do
        # once the depth of a shallow clone reaches the max depth of the origin
        # repo, Git silenty turns it into a deep clone
        if [ ! -e .git/shallow ]; then
          echo "Reached max depth of the origin repo while deepening the shallow clone, it's a deep clone now"
          break
        fi

        echo "Could not find ref ${ref} in a shallow clone of depth ${d}"

        (( d *= 2 ))

        if [ "$d" -gt $max_depth ]; then
          echo "Reached depth threshold ${max_depth}, falling back to deep clone..."
          git fetch --unshallow origin

          break
        fi

        echo "Deepening the shallow clone to depth ${d}..."
        git fetch --depth $d origin
      done
    fi
    git checkout -q $ref
fi

invalid_key() {
  echo "Invalid GPG key in: ${commit_verification_keys}"
  exit 2
}

commit_not_signed() {
  commit_id=$(git rev-parse ${ref})
  echo "The commit ${commit_id} is not signed"
  exit 1
}

if [ ! -z "${commit_verification_keys}" ] || [ ! -z "${commit_verification_key_ids}" ] ; then
  if [ ! -z "${commit_verification_keys}" ]; then
    echo "${commit_verification_keys}" | gpg --batch --import || invalid_key "${commit_verification_keys}"
  fi
  if [ ! -z "${commit_verification_key_ids}" ]; then
    echo "${commit_verification_key_ids}" | \
      xargs --no-run-if-empty -n1 gpg --batch --keyserver $gpg_keyserver --recv-keys
  fi
  git verify-commit $(git rev-list -n 1 $ref) || commit_not_signed
fi

if [ "$disable_git_lfs" != "true" ]; then
  git lfs fetch
  git lfs checkout
fi

echo $(git log -1 --oneline)
git clean --force --force -d
git submodule sync

if [ -f $GIT_CRYPT_KEY_PATH ]; then
  echo "unlocking git repo"
  git-crypt unlock $GIT_CRYPT_KEY_PATH
fi


submodule_parameters=""
if [ "$submodule_remote" != "false" ]; then
    submodule_parameters+=" --remote "
fi
if [ "$submodule_recursive" != "false" ]; then
    submodule_parameters+=" --recursive "
fi

if [ "$submodules" == "all" ]; then
  git submodule update --init  $depthflag $submodule_parameters
elif [ "$submodules" != "none" ]; then
  submodules=$(echo $submodules | jq -r '(.[])')
  for submodule in $submodules; do
    git submodule update --init $depthflag $submodule_parameters $submodule
  done
fi

if [ "$disable_git_lfs" != "true" ]; then
  git submodule foreach "git lfs fetch && git lfs checkout"
fi

for branch in $fetch; do
  git fetch origin $branch
  git branch $branch FETCH_HEAD
done

if [ "$ref" == "HEAD" ]; then
  return_ref=$(git rev-parse HEAD)
else
  return_ref=$ref
fi

# Store committer email in .git/committer. Can be used to send email to last committer on failed build
# Using https://github.com/mdomke/concourse-email-resource for example
git --no-pager log -1 --pretty=format:"%ae" > .git/committer

# Store git-resource returned version ref .git/ref. Useful to know concourse
# pulled ref in following tasks and resources.
echo "${return_ref}" > .git/ref

# Store short ref with templating. Useful to build Docker images with
# a custom tag
echo "${return_ref}" | cut -c1-7 | awk "{ printf \"${short_ref_format}\", \$1 }" > .git/short_ref

# Store commit message in .git/commit_message. Can be used to inform about
# the content of a successfull build.
# Using https://github.com/cloudfoundry-community/slack-notification-resource
# for example
git log -1 --format=format:%B > .git/commit_message

metadata=$(git_metadata)

if [ "$clean_tags" == "true" ]; then
  git tag | xargs git tag -d
fi

jq -n "{
  version: {ref: $(echo $return_ref | jq -R .)},
  metadata: $metadata
}" > ${METAFILE}
