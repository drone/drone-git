local windows_pipe = '\\\\\\\\.\\\\pipe\\\\docker_engine';
local windows_pipe_volume = 'docker_pipe';
local versions = [
  //'1803',
  '1809',
];
local trigger = {
  ref: [
    'refs/heads/master',
    'refs/tags/**',
  ],
};
local pipeline_name(version) = 'Windows ' + version;

local pipeline(version, arch) = {
  kind: 'pipeline',
  name: pipeline_name(version),

  platform: {
    os: 'windows',
    arch: arch,
    version: version,
  },

  steps: [{
    name: 'git',
    image: 'plugins/docker:windows-1809',  // TODO: This should just use the manifest
    settings: {
      repo: 'drone/git',
      dockerfile: 'docker/Dockerfile.windows.' + version,
      auto_tag: true,
      auto_tag_suffix: 'windows-' + version + '-' + arch,

      username: { from_secret: 'docker_username' },
      password: { from_secret: 'docker_password' },

      // Windows specific options
      daemon_off: true,
      purge: 'false',  // TODO: Fix bug where setting false won't generate the yaml value
    },
    volumes: [{ name: windows_pipe_volume, path: windows_pipe }],
  }],

  volumes: [{ name: windows_pipe_volume, host: { path: windows_pipe } }],
  trigger: trigger,
};

[
  pipeline(version, 'amd64')
  for version in versions
] + [
  {
    kind: 'pipeline',
    name: 'Image Manifest',

    steps: [{
      name: 'manifest',
      image: 'plugins/manifest',
      settings: {
        spec: 'docker/manifest.tmpl',
        ignore_missing: true,

        username: { from_secret: 'docker_username' },
        password: { from_secret: 'docker_password' },
      },
    }],

    depends_on: [
      pipeline_name(version)
      for version in versions
    ],
    trigger: trigger,
  },
]
