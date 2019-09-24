import axios from "axios";

const client = axios.create({
    baseURL: process.env.VUE_APP_API_ADMIN,
    timeout: 1000,
    withCredentials: true
});

client.interceptors.response.use(resp => {
    let res = resp.data;
    if (res.code === 10002) {
        window.location.href = res.data.url;
        return false;
    }

    return resp;
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