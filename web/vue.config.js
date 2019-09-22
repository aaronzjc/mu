module.exports = {
    outputDir: "../dist",
    assetsDir: "static",
    pages: {
        index: {
            entry: "src/index/main.js",
            template: "public/index.html",
            filename: "index.html",
            title: "首页",
            chunks: ['chunk-vendors', 'chunk-common', 'index']
        },
        admin: {
            entry: "src/admin/main.js",
            template: "public/admin.html",
            filename: "admin.html",
            title: "后台",
            chunks: ['chunk-vendors', 'chunk-common', 'admin']
        }
    },
    configureWebpack:{
        performance: {
            hints: false
        }
    }
};