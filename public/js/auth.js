/**
 * NeuroMetric Session Manager (Go Backend Version)
 * Replaces Firebase with stateless REST API + JWT
 */

const API_BASE = "/api";

export const SessionManager = {
    /**
     * Register a new user
     */
    async register(username, password, userData) {
        try {
            const response = await fetch(`${API_BASE}/auth/register`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    username,
                    password,
                    full_name: userData.fullName,
                    age: parseInt(userData.age),
                    gender: userData.gender
                })
            });

            const result = await response.json();
            if (!response.ok) throw new Error(result.error || "Registration failed");

            // After registration, we automatically sign in
            return await this.signIn(username, password);
        } catch (error) {
            console.error("Registration error:", error);
            throw error;
        }
    },

    /**
     * Sign in an existing user
     */
    async signIn(username, password) {
        try {
            const response = await fetch(`${API_BASE}/auth/login`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password })
            });

            const result = await response.json();
            if (!response.ok) throw new Error(result.error || "Invalid username or password");

            // Store JWT and Profile
            localStorage.setItem('neuro_token', result.token);
            localStorage.setItem('neuro_user', result.user.username);
            localStorage.setItem('neuro_full_name', result.user.full_name);

            return result.user;
        } catch (error) {
            console.error("Sign-in error:", error);
            throw error;
        }
    },

    /**
     * Sign out
     */
    async logout() {
        localStorage.removeItem('neuro_token');
        localStorage.removeItem('neuro_user');
        localStorage.removeItem('neuro_full_name');
        window.location.href = 'index.html';
    },

    /**
     * Check if a user is authenticated (Stateless check)
     */
    async checkSession(redirectIfEmpty = true) {
        const token = localStorage.getItem('neuro_token');
        
        if (!token) {
            if (redirectIfEmpty && !window.location.pathname.endsWith('index.html')) {
                window.location.href = 'index.html';
            }
            return null;
        }

        try {
            // Verify token with backend
            const response = await fetch(`${API_BASE}/user/profile`, {
                headers: { "Authorization": `Bearer ${token}` }
            });

            if (!response.ok) {
                this.logout();
                return null;
            }

            return await response.json();
        } catch (e) {
            console.error("Session check failed:", e);
            return null;
        }
    },

    /**
     * Get current user profile from localStorage
     */
    getCurrentUser() {
        return {
            username: localStorage.getItem('neuro_user'),
            fullName: localStorage.getItem('neuro_full_name'),
            token: localStorage.getItem('neuro_token')
        };
    }
};

// Export dummy db to prevent breakage in other files temporarily
export const db = null;
export const auth = null;
