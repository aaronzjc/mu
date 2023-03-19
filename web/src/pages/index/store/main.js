import { defineStore } from 'pinia'

export const useMainStore = defineStore('main', {
    state: () => ({
        userInfo: {
            username: '',
            avatar: ''
        }
    }),
    getters: {
        isLogin: (state) => {
            return state.userInfo.username != ''
        }
    },
    actions: {
        setUser(payload) {
            if (payload.username) {
                this.userInfo.username = payload.username
            }
            if (payload.avatar) {
                this.userInfo.avatar = payload.avatar
            }
        }
    }
})