import axios from "axios";
import * as ls from  "./ls"

const client = axios.create({
    baseURL: process.env.VUE_APP_URL,
    withCredentials: true
});

client.interceptors.request.use(req => {
    let token = ls.Get("token")
    if (token) {
        req.headers.Authorization = token;
    }
    return req;
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