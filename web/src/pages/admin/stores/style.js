import { defineStore } from 'pinia'

export const useStyleStore = defineStore('style', {
    state: () => ({
        isAsideMobileOpen: false,
        isNavMobileOpen: false
    }),
    actions: {
        toggleAside() {
            this.isAsideMobileOpen = !this.isAsideMobileOpen
        },
        toggleNav() {
            this.isNavMobileOpen = !this.isNavMobileOpen
        }
    }
})
