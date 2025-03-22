import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import Toast from "vue-toastification";
import "vue-toastification/dist/index.css";

// Create the app
const app = createApp(App);

// Add plugins before mounting
app.use(Toast, {
    transition: "Vue-Toastification__bounce",
    maxToasts: 3,
    newestOnTop: true,
    position: "bottom-right",
    timeout: 8000,
    closeOnClick: true,
    pauseOnFocusLoss: true,
    pauseOnHover: true,
    hideProgressBar: false,

    filterBeforeCreate: (toast, toasts) => {
        // Don't show duplicate toasts with the same content
        if (toasts.filter(t => t.content === toast.content).length > 0) {
            return false;
        }
        return toast;
    }
});

// Mount the app after configuration is complete
app.mount('#app');