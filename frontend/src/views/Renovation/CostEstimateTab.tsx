import { useState, useEffect } from 'react';
import { apiClient } from '../../services/api';

interface LaborTask {
    id: string;
    label: string;
    status: 'pending' | 'in_progress' | 'completed';
    unit_price: number;
    unit: string;
    quantity: number;
    amount: number;
}

export const CostEstimateTab = ({ renovationId }: { renovationId: string }) => {
    const [tasks, setTasks] = useState<LaborTask[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(false);

    useEffect(() => {
        const fetchTasks = async () => {
            try {
                const response = await apiClient.get(`/api/renovations/${renovationId}`);
                
                const rawTasks = response.data?.data?.LaborTasks || [];
                const formattedTasks = rawTasks.map((t: any) => ({
                    id: t.ID,
                    label: t.Label,
                    status: t.Status,
                    unit_price: t.UnitPrice,
                    unit: t.Unit,
                    quantity: t.Quantity,
                    amount: t.Amount
                }));
                
                setTasks(formattedTasks);
            } catch (err) {
                console.error("Błąd pobierania kosztorysu", err);
                setError(true);
            } finally {
                setIsLoading(false);
            }
        };

        fetchTasks();
    }, [renovationId]);

    if (isLoading) return <div className="p-6 text-text-muted animate-pulse">Analiza kosztorysu...</div>;
    if (error) return <div className="p-6 text-red-500">Błąd odczytu danych kosztorysu.</div>;

    const formatPLN = (amount: number) => `${amount.toLocaleString('pl-PL', { minimumFractionDigits: 2 })} zł`;

    return (
        <div className="bg-bg-surface border border-bg-border rounded-card overflow-hidden">
            <div className="p-6 border-b border-bg-border flex justify-between items-center">
                <h2 className="text-xl font-bold text-text-main">Zestawienie Prac (Labor Tasks)</h2>
                {/* Przycisk do dodawania nowego zadania (przygotowany pod przyszłą funkcję) */}
                <button className="bg-primary hover:bg-primary-hover text-white px-4 py-2 rounded text-sm font-bold transition-colors">
                    + Dodaj zadanie
                </button>
            </div>
            
            <div className="overflow-x-auto">
                <table className="w-full text-left border-collapse">
                    <thead>
                        <tr className="border-b border-bg-border text-xs text-text-muted uppercase tracking-wider bg-bg-base/30">
                            <th className="p-4 font-medium">Nazwa usługi (Label)</th>
                            <th className="p-4 font-medium text-center">Status</th>
                            <th className="p-4 font-medium text-right">Ilość</th>
                            <th className="p-4 font-medium text-right">Cena jedn.</th>
                            <th className="p-4 font-medium text-right">Razem</th>
                        </tr>
                    </thead>
                    <tbody className="text-sm">
                        {tasks.length === 0 ? (
                            <tr>
                                <td colSpan={5} className="p-6 text-center text-text-muted">Kosztorys jest pusty.</td>
                            </tr>
                        ) : (
                            tasks.map((task) => (
                                <tr key={task.id} className="border-b border-bg-border last:border-0 hover:bg-bg-base/50 transition-colors">
                                    <td className="p-4 font-medium text-text-main">{task.label}</td>
                                    <td className="p-4 text-center">
                                        <span className={`px-2 py-1 rounded text-xs font-medium uppercase ${
                                            task.status === 'completed' ? 'bg-emerald-900/50 text-emerald-400' :
                                            task.status === 'in_progress' ? 'bg-blue-900/50 text-blue-400' :
                                            'bg-slate-700 text-slate-300'
                                        }`}>
                                            {task.status.replace('_', ' ')}
                                        </span>
                                    </td>
                                    <td className="p-4 text-right text-text-muted">{task.quantity} {task.unit}</td>
                                    <td className="p-4 text-right text-text-muted">{formatPLN(task.unit_price)}</td>
                                    <td className="p-4 text-right font-bold text-text-main">{formatPLN(task.amount)}</td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
};