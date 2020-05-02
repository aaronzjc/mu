const PREFIX = "mu_cache_"

function formatKey(key) {
    return PREFIX + key
}

function ttlKey(key) {
    return formatKey(key) + "_ttl"
}

export function Set(key, val, ttl) {
    let expireAt = -1;
    if (ttl > 0) {
        expireAt = (new Date()).getTime() + parseInt(ttl) * 1000
    }
    if (expireAt > 0) {
        localStorage.setItem(ttlKey(key), "" + expireAt)
    }
    localStorage.setItem(formatKey(key), val)
}

/**
 * @return {string|boolean}
 */
export function Get(key) {
    let expire = localStorage.getItem(ttlKey(key))
    if (expire === "" || expire == null) {
        // do nothing
    } else if (expire != -1 && parseInt(expire) < (new Date()).getTime()) {
        Del(key)
        return false
    }
    let res = localStorage.getItem(formatKey(key))
    if (res == undefined || res == null) {
        return false;
    }
    return res
}

export function Del(key) {
    localStorage.removeItem(formatKey(key))
    localStorage.removeItem(ttlKey(key))
}