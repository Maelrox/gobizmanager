import { createContext } from 'preact';
import { useContext, useCallback } from 'preact/hooks';
import { signal, useSignal, useComputed } from '@preact/signals';
import axios from 'axios';

const AuthContext = createContext();

export function AuthProvider({ children }) {
  const token = useSignal(localStorage.getItem('token') || '');
  const refreshToken = useSignal(localStorage.getItem('refreshToken') || '');
  const isAuthenticated = useComputed(() => !!token.value);

  const setTokens = useCallback((newToken, newRefreshToken) => {
    token.value = newToken;
    refreshToken.value = newRefreshToken;
    localStorage.setItem('token', newToken);
    localStorage.setItem('refreshToken', newRefreshToken);
  }, []);

  const clearTokens = useCallback(() => {
    token.value = '';
    refreshToken.value = '';
    localStorage.removeItem('token');
    localStorage.removeItem('refreshToken');
  }, []);

  // Configure axios interceptors
  axios.interceptors.request.use(
    (config) => {
      if (token.value) {
        config.headers.Authorization = `Bearer ${token.value}`;
      }
      return config;
    },
    (error) => Promise.reject(error)
  );

  axios.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config;

      if (error.response?.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true;

        try {
          const response = await axios.post('/api/auth/refresh', {
            refresh_token: refreshToken.value,
          });

          const { token: newToken, refresh_token: newRefreshToken } = response.data;
          setTokens(newToken, newRefreshToken);

          originalRequest.headers.Authorization = `Bearer ${newToken}`;
          return axios(originalRequest);
        } catch (refreshError) {
          clearTokens();
          window.location.href = '/login';
          return Promise.reject(refreshError);
        }
      }

      return Promise.reject(error);
    }
  );

  return (
    <AuthContext.Provider value={{ token, refreshToken, isAuthenticated, setTokens, clearTokens }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
} 