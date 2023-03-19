import { defineStore } from 'pinia'

export const useMainStore = defineStore('main', {
    state: () => ({
        token: '',
        userInfo: {
            username: '',
            avatar: ''
        }
    }),
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
