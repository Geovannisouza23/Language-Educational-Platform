module.exports = {
  apps: [
    {
      name: 'language-platform-web',
      script: 'node_modules/next/dist/bin/next',
      args: 'start',
      cwd: './frontend/web',
      instances: 'max',
      exec_mode: 'cluster',
      env: {
        NODE_ENV: 'production',
        PORT: 3000,
      },
    },
  ],
};
