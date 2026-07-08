export interface UserClaims{
    userID: string;
    role: 'admin'|'client';
    email: string;
    exp:number;
}

export interface AuthState{
    token: string | null;
    user: UserClaims | null;
    isAuthenticated: boolean;
}