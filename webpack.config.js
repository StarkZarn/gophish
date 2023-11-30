const path = require('path');

module.exports = {
    mode: 'production', // or 'development'
    context: path.resolve(__dirname, 'static', 'js', 'src', 'app'),
    entry: {
        passwords: './passwords',
        users: './users',
        webhooks: './webhooks',
    },
    output: {
        path: path.resolve(__dirname, 'static', 'js', 'dist', 'app'),
        filename: '[name].min.js'
    },
    resolve: {
        extensions: ['.js']
    },
    module: {
        rules: [{
            test: /\.js$/,
            exclude: /node_modules/,
            use: {
                loader: "babel-loader"
            }
        }]
    }
};
