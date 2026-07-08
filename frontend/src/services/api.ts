import axios from "axios";
import {config} from '../config'

export const apiClient = axios.create({
    baseURL: config.baseUrl||"https://localhost:8081",
    headers:{
        'Content-Type': 'application/json'
    },
});

apiClient.interceptors.request.use(
    (config)=>{
        const token = localStorage.getItem('token');
        if(token){
            config.headers.Authorization =`Bearer ${token}`
        }
        return config;
    },
    (error)=>{
        return Promise.reject(error)
    }
)

apiClient.interceptors.response.use(
    (response)=>{
        return response
    },
    (error)=>{
        if(error.response && error.response.status === 401){
            console.warn("Nieaktywny token");
            localStorage.removeItem('token')
            window.location.href = '/'
        }
        return Promise.reject(error);
    }
)