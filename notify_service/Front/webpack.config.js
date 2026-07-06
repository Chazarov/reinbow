const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    // Точка входа — главный JS-файл вашего приложения
    entry: './static/index.js',
    
    // Настройка выходных файлов
    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'dist'),
        clean: true, // очищает dist/ при каждой сборке
    },
    
    // Правила обработки разных типов файлов
    module: {
        rules: [
            {
                test: /\.css$/i,      // все файлы .css
                use: ['style-loader', 'css-loader'], // обрабатываются этими загрузчиками
            },
        ],
    },
    
    // Плагины
    plugins: [
        new HtmlWebpackPlugin({
            template: './static/index.html', // берём ваш HTML как шаблон
        }),
    ],
    
    mode: 'development', // режим разработки (для продакшена — 'production')
};