const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const Dotenv = require('dotenv-webpack');
const dotenv = require('dotenv');

module.exports = (env, argv) => {
  const isProduction = argv.mode === 'production';
  console.log("FeedServiceFrontend: загружен " + argv.mode + " конфиг")

  dotenv.config({
    path: `./.env.${argv.mode}`
  });

  return {
    entry: './js/index.js',

    output: {
      filename: 'bundle.[contenthash].js',
      path: path.resolve(__dirname, 'dist'),
      clean: true,
    },

    devServer: {
      static: './dist',
      hot: true,
      port: 3000,
      open: false,
      proxy: [
        {
          context: ['/api'],
          target: 'http://localhost:8080',
          changeOrigin: true,
          pathRewrite: {
            '^/api': '',
          }
        }
      ]
    },

    plugins: [
      new HtmlWebpackPlugin({
        template: './index.html',
        filename: 'index.html',
      }),
      new Dotenv({
        path: `./.env${isProduction ? '.production' : '.development'}`
      }),
    ],

    module: {
      rules: [
        {
          test: /\.(png|jpe?g|gif|svg)$/i,
          type: 'asset/resource',
        },
        {
          test: /\.css$/i,
          use: ['style-loader', 'css-loader'],
        },
      ],
    },
  }
};