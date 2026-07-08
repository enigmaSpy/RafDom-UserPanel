import { createFileRoute } from "@tanstack/react-router";
import { useAuth } from "../context/AuthContext";
import { LoginPage } from "../views/Login/LoginPage";
import { DashboardAdmin } from "../views/Dashboard/DashboardAdmin";
import { DashboardClient } from "../views/Dashboard/DashboardClient";

export const Route = createFileRoute('/')({
    component: IndexComponent,
})

function IndexComponent(){
    const {isAuthenticated, user, login} = useAuth()
    if(!isAuthenticated){
        return <LoginPage onLoginSuccess={login}/>
    }
    return user?.role === 'admin' ? <DashboardAdmin/>:<DashboardClient/>
}