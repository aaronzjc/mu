import { createApp } from 'vue'
import Toast from '@adm/components/toast/Toast.vue'

const toastContainer = document.createElement('div')
toastContainer.className = 'notice-container'

const ToastInstance = {
    show(text, level = '') {
        toastContainer.innerHTML = ''
        const toastApp = createApp(Toast, { text, level }).mount(
            document.createElement('div')
        )
        toastContainer.insertAdjacentElement('afterbegin', toastApp.$el)
        document.getElementById('app').appendChild(toastContainer)
    },
    error(text) {
        this.show(text, 'danger')
    },
    warn(text) {
        this.show(text, 'warning')
    },
    success(text) {
        this.show(text, 'success')
    }
}

Toast.install = (app, options) => {
    app.config.globalProperties.$toast = ToastInstance
}

function useToast() {
    return ToastInstance
}

export {Toast, useToast}