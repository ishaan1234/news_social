process.env.BABEL_ENV = 'development';
process.env.NODE_ENV = 'development';

const { execFileSync } = require('child_process');
const path = require('path');
const { defineConfig } = require('cypress');
const { devServer } = require('@cypress/webpack-dev-server');

const reactScriptsRoot = path.dirname(
  require.resolve('react-scripts/package.json', { paths: [__dirname] })
);
const craWebpackConfig = require(
  path.join(reactScriptsRoot, 'config', 'webpack.config.js')
);
const cypressDir = path.resolve(__dirname, 'cypress');
const componentCssInput = path.resolve(__dirname, 'src', 'index.css');
const componentCssOutput = path.resolve(
  __dirname,
  'cypress',
  'support',
  'component.css'
);
const tailwindCli = require.resolve('tailwindcss/lib/cli.js', {
  paths: [__dirname],
});

const webpackConfig = craWebpackConfig('development');
const oneOfRule = webpackConfig.module.rules.find((rule) => Array.isArray(rule.oneOf));

const buildComponentStyles = () => {
  try {
    execFileSync(
      process.execPath,
      [tailwindCli, '-i', componentCssInput, '-o', componentCssOutput],
      {
        cwd: __dirname,
        env: {
          ...process.env,
          NODE_ENV: 'development',
        },
        stdio: 'pipe',
      }
    );
  } catch (error) {
    const message =
      error instanceof Error ? error.message : 'Unknown Tailwind build failure.';

    throw new Error(`Failed to build Cypress component styles. ${message}`);
  }
};

if (oneOfRule) {
  oneOfRule.oneOf.forEach((rule) => {
    const isBabelLoader =
      typeof rule.loader === 'string' && rule.loader.includes('babel-loader');

    if (!isBabelLoader) {
      return;
    }

    const existingInclude = Array.isArray(rule.include)
      ? rule.include
      : rule.include
        ? [rule.include]
        : [];

    rule.include = [...existingInclude, cypressDir];
  });
}

webpackConfig.resolve.plugins = (webpackConfig.resolve.plugins || []).filter(
  (plugin) => plugin.constructor?.name !== 'ModuleScopePlugin'
);

module.exports = defineConfig({
  component: {
    specPattern: 'cypress/component/**/*.cy.{js,jsx}',
    supportFile: 'cypress/support/component.js',
    indexHtmlFile: 'cypress/support/component-index.html',
    devServer(devServerConfig) {
      buildComponentStyles();

      return devServer({
        ...devServerConfig,
        framework: 'react',
        webpackConfig,
      });
    },
  },
});
