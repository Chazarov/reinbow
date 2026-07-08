const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    entry: './static/index.js',

    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'dist'),
        clean: true,
    },

    module: {
        rules: [
            {
                test: /\.css$/i,
                use: ['style-loader', 'css-loader'],
            },
        ],
    },

    resolve: {
        alias: {
            '@utils': path.resolve(__dirname, 'static/utils'),
            '@notifications': path.resolve(__dirname, 'static/notifications'),
            '@auth': path.resolve(__dirname, 'static/auth'),
            '@settings': path.resolve(__dirname, 'static/settings'),
            '@webhooks': path.resolve(__dirname, 'static/webhooks'),
        },
    },

    plugins: [
        new HtmlWebpackPlugin({
            template: './static/index.html',
        }),
    ],

    mode: 'development',
};