
if (!(Test-Path .git)) {
	git init
	git remote add origin $Env:DRONE_REMOTE_URL
}

git fetch origin "+refs/tags/${Env:DRONE_TAG}:"
git checkout -qf FETCH_HEAD
