
Set-Variable -Name "FLAGS" -Value ""
if ($Env:PLUGIN_DEPTH) {
    Set-Variable -Name "FLAGS" -Value "--depth=$Env:PLUGIN_DEPTH" 
}

if (!(Test-Path .git)) {
	git init
	git remote add origin $Env:DRONE_REMOTE_URL
}

git fetch $FLAGS origin "+refs/heads/${Env:DRONE_COMMIT_BRANCH}:"
git checkout $Env:DRONE_COMMIT_BRANCH

git fetch origin $Env:DRONE_COMMIT_REF:
git rebase $Env:DRONE_COMMIT_SHA
