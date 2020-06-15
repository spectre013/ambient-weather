module.exports = {
    devServer: {
        proxy: {
            '/api': {
                target: 'http://localhost:3000',
                changeOrigin: true
            },
            // '/api/ws': {
            //     target: 'ws://localhost:3000',
            //     ws: true
            // },
        }
    }
}