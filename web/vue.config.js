module.exports = {
    outputDir: "../public",
    assetsDir: "static",
    pages: {
        index: {
            entry: "src/pages/index/main.js",
            template: "public/index.html",
            filename: "index.html",
            title: "首页",
        },
        admin: {
            entry: "src/pages/admin/main.js",
            template: "public/admin.html",
            filename: "admin.html",
            title: "后台",
        }
    },
    chainWebpack: config => {
        config.performance.hints = false
        config.optimization.delete('splitChunks')
    },
    devServer: {
        disableHostCheck: true
    }
};