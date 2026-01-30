import axios from 'axios';

const authApi = axios.create({
  baseURL: '/api/auth',
});

export interface LoginRequest {
    email: string;
    password: string;
}

export interface LoginResponse {
    token: string;
    user: AuthUser;
}

export interface AuthUser {
    id: string;
    email: string;
    name: string;
    role: string;
}

export const authApiClient = {
    login: async (credentials: LoginRequest): Promise<LoginResponse> => {
        const response = await authApi.post('/auth/login', credentials);
        return response.data;
    },

    register: async (credentials: LoginRequest & { name: string }): Promise<LoginResponse> => {
        const response = await authApi.post('/auth/register', credentials);
        return response.data;
    },

    verifyToken: async (token: string): Promise<AuthUser> => {
        const response = await authApi.get('/auth/verify', {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        return response.data.user;
    },
};

export default authApi;
