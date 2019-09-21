import axios from "axios";

const client = axios.create({
    baseURL: process.env.VUE_APP_API,
    timeout: 1000
});

export function Get(url, params, headers) {
    if (!params) {
        params = {};
    }

    let config = {
        method: "get",
        url: url,
        params: params
    };
    if (headers) {
        config.headers = headers;
    }

    return client(config)
}

export function Post(url, data, headers) {
    let config= {
        method: 'post',
        url: url,
        params: {}
    };
    if (headers) {
        config.headers = headers;
    }
    if (data) {
        config.data = data;
    }
    return client(config);
}

export default client;