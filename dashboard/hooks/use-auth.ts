import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface User {
  id: string;
  name: string;
  email: string;
  role: string;
  org_id: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (user: User) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      login: (user) => {
        const token = `token-${Date.now()}`;
        if (typeof window !== 'undefined') {
          localStorage.setItem('token', token);
          document.cookie = `auth-token=${token};path=/;max-age=86400`;
        }
        set({ user, token, isAuthenticated: true });
      },
      logout: () => {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('token');
          document.cookie = 'auth-token=;path=/;max-age=0';
        }
        set({ user: null, token: null, isAuthenticated: false });
      },
    }),
    {
      name: 'auth-storage',
    }
  )
);
