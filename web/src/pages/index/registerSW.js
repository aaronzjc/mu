/**
 * 注册service worker。
 *
 * 为什么不用Vue-PWA插件来做PWA呢？因为我的项目是多页面。如果使用那个插件，
 * 它会将我所有页面都插入PWA的那些内容。事实上，我只想首页支持PWA。
 */
window.addEventListener("load", function () {
    if ('serviceWorker' in navigator) {
        navigator.serviceWorker.register("/sw.js").then(function(registration) {
            console.log('Registration successful, scope is: ', registration.scope);
        }).catch(function(error) {
            console.log('Service worker registration failed, error: ', error);
        });
    }
})