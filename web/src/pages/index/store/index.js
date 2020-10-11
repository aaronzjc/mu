import { createStore } from 'vuex'

import account from "./module/account"

export const store = createStore({
    modules: {
        account,
    }
})
