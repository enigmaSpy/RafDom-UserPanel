import { useState } from "react";

interface LoginViewProps {
    onLoginSuccess: (token: string) => void;
}
const BASE_URL = "https://localhost:8081"

const LoginForm = ({ onLoginSuccess }: LoginViewProps) => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false); 

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (!email || !password) {
            setError('Email i hasło są wymagane');
            return;
        }

        setIsLoading(true);

        try {
            const response = await fetch(`${BASE_URL}/api/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, password })
            });
            
            const data = await response.json();
            
            if (!response.ok) {
                setError(data.error || 'Logowanie nie powiodło się');
                setIsLoading(false);
                return;
            }

            onLoginSuccess(data.token);
            
        } catch (error) {
            setError('Błąd połączenia z serwerem');
            console.error('Login error:', error);
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            {error && (
                <div className="errorCard p-3 rounded">
                    {error}
                </div>
            )}
            <input
                name="email"
                type="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled={isLoading}
                required
                className="p-3 rounded bg-bg-base border border-bg-border focus:outline-none focus:border-primary text-text-main"
            />
            <input
                name="password"
                type="password"
                placeholder="Hasło"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={isLoading}
                required
                className="p-3 rounded bg-bg-base border border-bg-border focus:outline-none focus:border-primary text-text-main"
            />
            <button
                type="submit"
                disabled={isLoading}
                className={`formButton-primary p-3 rounded  text-white font-bold transition-colors ${
                    isLoading ? 'opacity-50 cursor-not-allowed' : ''
                }`}
            >
                {isLoading ? 'Logowanie...' : 'Zaloguj się'}
            </button>
        </form>
    );
}

export const LoginPage = ({ onLoginSuccess }: LoginViewProps) => {
    return (
        <div className="max-w-md mx-auto mt-20 p-6 bg-bg-surface rounded-card border border-bg-border shadow-xl">
            <h2 className="text-2xl font-bold mb-6 text-center text-text-main">RenovManager</h2>
            <LoginForm onLoginSuccess={onLoginSuccess} />
        </div>
    );
}