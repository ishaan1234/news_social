const { createProxyMiddleware } = require('http-proxy-middleware');

const proxyTarget = (process.env.API_PROXY_TARGET || 'http://localhost:8080').replace(
  /\/$/,
  ''
);

module.exports = function setupProxy(app) {
  app.use(
    '/news',
    createProxyMiddleware({
      target: proxyTarget,
      changeOrigin: true,
    })
  );

  app.use(
    '/auth',
    createProxyMiddleware({
      target: proxyTarget,
      changeOrigin: true,
    })
  );

  app.use(
    '/api',
    createProxyMiddleware({
      target: proxyTarget,
      changeOrigin: true,
    })
  );
};
