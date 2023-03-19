import axios from "axios";

const client = axios.create({
    baseURL: import.meta.env.VITE_APP_URL,
    withCredentials: true
});

client.interceptors.request.use(req => {
    let token = localStorage.getItem(import.meta.env.VITE_TOKEN_KEY)
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