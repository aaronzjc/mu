const CACHE_NAME = "mu-pwa-1";

const OFFLINE_PAGE = "/pwa/offline.html";

const FILES_TO_CACHE = [
    "/index.manifest",
    "/favicon.png",
    "index.html",
    "/pwa/imgs/android-chrome-512x512.png",
    "/pwa/imgs/android-chrome-192x192.png",
    OFFLINE_PAGE
];

self.oninstall = function(evt) {
    evt.waitUntil(
        caches.open(CACHE_NAME).then((cache) => {
            console.log('[SW] Pre-caching resources');
            return cache.addAll(FILES_TO_CACHE);
        })
    );
}

self.onactivate = function(evt) {
    evt.waitUntil(
        caches.keys().then((keyList) => {
            return Promise.all(keyList.map((key) => {
                if (key !== CACHE_NAME) {
                    console.log('[SW] Removing old cache', key);
                    return caches.delete(key);
                }
            }));
        })
    );
}

self.onfetch = function(evt) {
    if (evt.request.mode !== 'navigate') {
        return;
    }
    evt.respondWith(
        fetch(evt.request).catch(() => {
            return caches.open(CACHE_NAME).then((cache) => {
                return cache.match(OFFLINE_PAGE);
            });
        })
    );
}