import axios from 'axios';

const api = axios.create({
    baseURL: '/kong-admin',
});

export const getServices = async () => {
    const response = await api.get('/services');
    return response.data.data;
};

export const getRoutes = async () => {
    const response = await api.get('/routes');
    return response.data.data;
};

export default api;
