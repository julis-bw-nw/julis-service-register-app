const path = require('path');
const HtmlWebpackPlugin = require("html-webpack-plugin");


module.exports = {
    entry: './src/app.js',
    output: {
        path: __dirname + '/dist/build',
        filename: "bundle.js"
    },
    devServer: {
        contentBase: path.join(__dirname, 'src'),
        port: 8000,
        watchContentBase: true,
    },
    module: {
        rules: [
            {
                test: /\.scss$/,
                use: [
                    {loader: "style-loader"},
                    {loader: "css-loader"},
                    {loader: "sass-loader"}
                ]
            }
        ]
    },
    plugins: [
        new HtmlWebpackPlugin({
            template: 'src/index.html'
        })
    ]
}