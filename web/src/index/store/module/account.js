const state = {
    id: 0,
    username: "",
    avatar: ""
};

const getters = {
    isLogin(state) {
        return state.id > 0;
    },
    getUsername(state) {
        return state.username;
    },
    getAvatar(state) {
        return state.avatar;
    }
};

const actions = {
    initUser({ commit }, info) {
        commit('initUser', info)
    }
};

const mutations = {
    initUser(state, info) {
        console.log(info);
        state.id = info.id;
        state.username = info.username;
        state.avatar = info.avatar;
    }
};

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}