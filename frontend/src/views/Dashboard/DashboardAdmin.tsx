import { useEffect, useState } from "react"
import type { DashboardStats } from "../../types/dashboard"
import { apiClient } from "../../services/api";
import { Tile } from "../../components/Dashboard/Tile";
import { RecentRenovations } from "../../components/Dashboard/RecentRenovation";

export const DashboardAdmin = ()=>{
    const [stats, setStats] = useState<DashboardStats | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchDashboardData = async () => {
            try {
                const response = await apiClient.get<DashboardStats>('/api/admin/dashboard');
                setStats(response.data);
            } catch (err) {
                console.error("Dashboard error:", err);
                setError('Nie udało się pobrać statystyk.');
            } finally {
                setIsLoading(false);
            }
        };

        fetchDashboardData();
    }, []);

    if (isLoading) {
        return <div className="p-6 text-text-main font-main animate-pulse">Inicjalizacja modułu analitycznego...</div>;
    }

    if (error) {
        return <div className="p-6 text-red-500 font-main bg-red-900/20 rounded border border-red-500">{error}</div>;
    }

    const formattedIncome = stats?.total_income 
        ? `${stats.total_income.toLocaleString('pl-PL')} zł` 
        : '0 zł';

    return (
        <div className="p-6 font-main">
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-text-main">Dashboard</h1>
                <p className="text-text-muted mt-1">Przegląd operacyjny systemu</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                
                <Tile title="Wszystkie Projekty" value={stats?.total_projects.toString()} />
                <Tile title="Aktywne (In Progress)" value={stats?.active_projects.toString()} />
                <Tile title="Zakończone" value={stats?.completed_projects.toString()} />
                <Tile 
                    title="Zaksięgowane wpłaty" 
                    value={formattedIncome} 
                    highlight 
                />
            </div>
            <RecentRenovations/>
        </div>
    );
};
