import { useEffect, useState } from "react";
import { apiClient } from "../../services/api";
import { Link } from "@tanstack/react-router";
import { StatusBadge } from "../ui/StatusBadge";
import { useNavigate } from '@tanstack/react-router';
import type { BackendRenovation, RenovationUI } from "../../types/dashboard";


interface ApiResponse {
    data: BackendRenovation[];
}

export const RecentRenovations = () => {
    const [renovations, setRenovations] = useState<RenovationUI[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(false);

    useEffect(() => {
        const fetchRecent = async () => {
            try {
                const response = await apiClient.get<ApiResponse>('/api/renovations/list');
                
                const rawArray = response.data.data;
                const safeData = Array.isArray(rawArray) ? rawArray : [];
                
                
                const adaptedData: RenovationUI[] = safeData.map(item => {
                    const calculatedBudget = item.LaborTasks 
                        ? item.LaborTasks.reduce((sum, task) => sum + task.Amount, 0)
                        : 0;

                    return {
                        id: item.ID,
                        name: item.Name,
                        client_name: `${item.Client?.Name || ''} ${item.Client?.Surname || ''}`.trim(),
                        status: item.Status,
                        total_budget: calculatedBudget
                    };
                });

                setRenovations(adaptedData.slice(0, 5));
                
            } catch (error) {
                console.error("Błąd pobierania listy remontów", error);
                setError(true);
            } finally {
                setIsLoading(false);
            }
        };
        fetchRecent();
    }, []);
    const navigate = useNavigate();
    return (
        <div className="bg-bg-surface border border-bg-border rounded-card overflow-hidden mt-8 shadow-sm">
            <div className="flex justify-between items-center p-6 border-b border-bg-border">
                <h2 className="text-xl font-bold text-text-main">Recent Renovations</h2>
                <Link to="/" className="text-primary hover:text-primary-hover text-sm font-medium transition-colors">
                    View all &rarr;
                </Link>
            </div>

            <div className="overflow-x-auto">
                <table className="w-full text-left border-collapse">
                    <thead>
                        <tr className="border-b border-bg-border text-xs text-text-muted uppercase tracking-wider bg-bg-base/30">
                            <th className="p-4 font-medium">Project Name</th>
                            <th className="p-4 font-medium">Client</th>
                            <th className="p-4 font-medium">Status</th>
                            <th className="p-4 font-medium">Budget</th>
                            <th className="p-4 font-medium"></th>
                        </tr>
                    </thead>
                    <tbody className="text-sm">
                        {isLoading && (
                            <tr>
                                <td colSpan={5} className="p-4 text-center text-text-muted animate-pulse">
                                    Ładowanie danych...
                                </td>
                            </tr>
                        )}
                        
                        {!isLoading && error && (
                            <tr>
                                <td colSpan={5} className="p-4 text-center text-red-400">
                                    Nie udało się załadować danych.
                                </td>
                            </tr>
                        )}

                        {!isLoading && !error && renovations.length === 0 && (
                            <tr>
                                <td colSpan={5} className="p-4 text-center text-text-muted">
                                    Brak aktywnych projektów.
                                </td>
                            </tr>
                        )}

                        {!isLoading && !error && renovations.map((project) => (
                            <tr 
                                key={project.id} 
                                onClick={() => navigate({ 
                                    to: '/admin/renovations/$id', 
                                    params: { id: project.id } 
                                })}
                                className="border-b border-bg-border last:border-0 hover:bg-bg-base/50 transition-colors cursor-pointer group"
                            >
                                <td className="p-4 font-medium text-text-main">{project.name}</td>
                                <td className="p-4 text-text-muted">{project.client_name || '-'}</td>
                                <td className="p-4">
                                    <StatusBadge status={project.status} />
                                </td>
                                <td className="p-4 font-medium text-text-main">
                                    {project.total_budget ? `${project.total_budget.toLocaleString('pl-PL')} zł` : '0 zł'}
                                </td>
                                <td className="p-4 text-right text-text-muted group-hover:text-primary transition-colors">
                                    &gt;
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};