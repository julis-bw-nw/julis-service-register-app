const path = require('path');
const HtmlWebpackPlugin = require("html-webpack-plugin");
const CopyPlugin = require("copy-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin").default;
const webpack = require("webpack");

module.exports = {
    entry: './src/app.js',
    output: {
        path: __dirname + '/dist/build',
        filename: "app.js"
    },
    devServer: {
        static: {
            directory: path.join(__dirname, 'src')
        },
        port: 9000
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
            },
            {
                test: /\.css$/i,
                use: [MiniCssExtractPlugin.loader, 'css-loader']
            }
        ]
    },
    plugins: [
        new webpack.ProvidePlugin({
            $: 'jquery',
            jQuery: 'jquery'
        }),
        new HtmlWebpackPlugin({
            template: 'src/index.html'
        }),
        new MiniCssExtractPlugin(),
        new CopyPlugin({
            patterns: [
                {
                    from: path.resolve(__dirname, 'node_modules/@shoelace-style/shoelace/dist/assets'),
                    to: path.resolve(__dirname, 'dist/build/shoelace/assets')
                },
                {
                    from: path.resolve(__dirname, 'src/assets'),
                    to: path.resolve(__dirname, 'dist/build/assets')
                },
            ]
        })
    ]
}