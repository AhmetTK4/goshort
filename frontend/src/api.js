import axios from 'axios';

const API_BASE = "http://localhost:8080";

export const shortenUrl = async (url) => {
    const response = await axios.post(`${API_BASE}/api/shorten`, { url });
    return response.data;
};

export const getStats = async (shortCode) => {
    const response = await axios.get(`${API_BASE}/api/stats/${shortCode}`);
    return response.data;
};

