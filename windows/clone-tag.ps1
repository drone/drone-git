
Set-Variable -Name "FLAGS" -Value ""
if ($Env:PLUGIN_DEPTH) {
    Set-Variable -Name "FLAGS" -Value "--depth=$Env:PLUGIN_DEPTH"
}

if (!(Test-Path .git)) {
	git init
	git remote add origin $Env:DRONE_REMOTE_URL
}

git fetch $FLAGS origin "+refs/tags/${Env:DRONE_TAG}:"
git checkout -qf FETCH_HEAD

if ($Env:PLUGIN_RECURSIVE) {
    git submodule update --init --recursive
}
