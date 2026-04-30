const { createProxyMiddleware } = require('http-proxy-middleware');

const proxyTarget = (
  process.env.API_PROXY_TARGET || 'http://localhost:8080'
).replace(/\/$/, '');

module.exports = function setupProxy(app) {
  [
    '/news',
    '/auth',
    '/api',
    '/feed',
    '/posts',
    '/post-likes',
    '/post-comments',
    '/following',
    '/profile',
    '/users',
  ].forEach((path) => {
    app.use(
      path,
      createProxyMiddleware({
        target: proxyTarget,
        changeOrigin: true,
      })
    );
  });
};
