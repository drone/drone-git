
Set-Variable -Name "FLAGS" -Value ""
if ($Env:PLUGIN_DEPTH) {
    Set-Variable -Name "FLAGS" -Value "--depth=$Env:PLUGIN_DEPTH"
}

if (!(Test-Path .git)) {
	Write-Host 'git init';
	git init
	Write-Host "git remote add origin $Env:DRONE_REMOTE_URL"
	git remote add origin $Env:DRONE_REMOTE_URL
}

Write-Host "git fetch $FLAGS origin +refs/heads/${Env:DRONE_COMMIT_BRANCH}:";
git fetch $FLAGS origin "+refs/heads/${Env:DRONE_COMMIT_BRANCH}:";
Write-Host "git checkout $Env:DRONE_COMMIT_SHA -f $Env:DRONE_COMMIT_BRANCH";
git checkout $Env:DRONE_COMMIT_SHA -b $Env:DRONE_COMMIT_BRANCH;

if ($Env:PLUGIN_RECURSIVE) {
    git submodule update --init --recursive
}
