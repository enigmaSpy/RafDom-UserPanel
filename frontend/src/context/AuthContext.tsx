import{jwtDecode} from 'jwt-decode'
import type { AuthState, UserClaims } from '../types/auth'
import React, { createContext, useContext, useEffect, useState } from 'react';

interface AuthConfigType extends AuthState{
    login: (token: string)=>void;
    logout: ()=>void;
}
type AuthProviderProps = {
    children: React.ReactNode
}

const AuthContext = createContext<AuthConfigType| undefined>(undefined);

export const AuthProvider=({children}:AuthProviderProps)=>{
    const [auth, setAuth] = useState<AuthState>({
        token: null,
        user: null,
        isAuthenticated: false,
    });

    useEffect(()=>{
        const savedToken = localStorage.getItem('token');
        if (savedToken){
            try {
                const decoded = jwtDecode<UserClaims>(savedToken)
                if(decoded.exp * 1000 > Date.now()){
                    setAuth({
                        token: savedToken,
                        user: decoded,
                        isAuthenticated: true
                    })
                }else{
                    localStorage.removeItem('token');
                }
            } catch (error) {
                localStorage.removeItem('token');
            }
        }
    },[]);

    const login = (token: string)=>{
        const decoded = jwtDecode<UserClaims>(token);
        localStorage.setItem('token',token);
        setAuth({token, user: decoded, isAuthenticated:true});
    }
    const logout = ()=>{
        localStorage.removeItem('token');
        setAuth({ token: null, user: null, isAuthenticated: false });
    }
    return (
        <AuthContext.Provider value={{...auth, login, logout}}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth=()=>{
    const context = useContext(AuthContext);
    if (!context) throw new Error('useAuth must be used within an AuthProvider');
    return context;
};

